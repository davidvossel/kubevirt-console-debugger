package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh/terminal"

	"kubevirt.io/client-go/kubecli"
)

func attachConsole(stdinReader, stdoutReader *io.PipeReader,
	stdinWriter, stdoutWriter *io.PipeWriter,
	message string, resChan <-chan error,
	autoResponder map[string]string) (err error) {

	// Always start by sending a new line to toggle console
	_, err = stdinWriter.Write([]byte("\n"))
	if err != nil {
		return err
	}

	stopChan := make(chan struct{}, 1)
	writeStop := make(chan error)
	readStop := make(chan error)
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		state, err := terminal.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("Make raw terminal failed: %s", err)
		}
		defer terminal.Restore(int(os.Stdin.Fd()), state)
	}
	fmt.Fprint(os.Stderr, message)

	in := os.Stdin
	out := os.Stdout

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		<-interrupt
		close(stopChan)
	}()

	// read output from console
	go func() {
		buf := make([]byte, 1024, 1024)
		line := ""
		for {
			n, err := stdoutReader.Read(buf)
			if err != nil && err != io.EOF {
				readStop <- err
				return
			}
			if n == 0 && err == io.EOF {
				return
			}

			split := strings.Split(string(buf[0:n]), "\n")
			for i, str := range split {
				line = line + str

				for key, val := range autoResponder {
					if strings.Contains(line, key) {
						fmt.Printf("--- MATCHED: [%s]\nKEY [%s]\nVAL [%s]\n", line, key, val)
						// Writing auto response
						_, err = stdinWriter.Write([]byte(val))
						if err == io.EOF {
							readStop <- err
							return
						}
						// ensure we only attempt the response once, so we don't get into a loop if no console password is set
						delete(autoResponder, key)
						// clear the line if a match occurs
						line = ""
					}
				}

				if i != len(split)-1 {
					// clear the line once a full new line is detected and no match
					line = ""
				}
			}

			// Writing out to the console connection
			_, err = out.Write(buf[0:n])
			if err == io.EOF {
				return
			}
		}

	}()

	// write input to console
	go func() {
		defer close(writeStop)
		buf := make([]byte, 1024, 1024)
		for {
			// reading from stdin
			n, err := in.Read(buf)
			if err != nil && err != io.EOF {
				writeStop <- err
				return
			}
			if n == 0 && err == io.EOF {
				return
			}

			// the escape sequence
			if buf[0] == 29 {
				return
			}
			// Writing out to the console connection
			_, err = stdinWriter.Write(buf[0:n])
			if err == io.EOF {
				return
			}
		}
	}()

	select {
	case <-stopChan:
	case err = <-readStop:
	case err = <-writeStop:
	case err = <-resChan:
	}

	return err
}

func run(args []string) error {

	timeout := 0

	if len(args) < 3 {
		fmt.Printf("invalid args\n")
		os.Exit(1)
	}
	vmi := args[1]
	namespace := args[2]

	virtCli, err := kubecli.GetKubevirtSubresourceClientFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return err
	}

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	// in -> stdinWriter | stdinReader -> console
	// out <- stdoutReader | stdoutWriter <- console
	// Wait until the virtual machine is in running phase, user interrupt or timeout
	resChan := make(chan error)
	runningChan := make(chan error)
	waitInterrupt := make(chan os.Signal, 1)
	signal.Notify(waitInterrupt, os.Interrupt)

	go func() {
		con, err := virtCli.VirtualMachineInstance(namespace).SerialConsole(vmi, &kubecli.SerialConsoleOptions{ConnectionTimeout: time.Duration(timeout) * time.Minute})
		runningChan <- err

		if err != nil {
			fmt.Printf("ERROR during setup: %v\n", err)
			return
		}

		resChan <- con.Stream(kubecli.StreamOptions{
			In:  stdinReader,
			Out: stdoutWriter,
		})
	}()

	select {
	case <-waitInterrupt:
		// Make a new line in the terminal
		fmt.Println()
		return nil
	case err = <-runningChan:
		if err != nil {
			fmt.Printf("ERROR after setup: %v\n", err)
			return err
		}
	}

	autoResponses := map[string]string{}

	autoResponses[" login:"] = "core\n"
	autoResponses["Password:"] = "core\n"
	autoResponses["~]$"] = "journalctl -f --no-tail\n"

	err = attachConsole(stdinReader, stdoutReader, stdinWriter, stdoutWriter,
		fmt.Sprint("Successfully connected to ", vmi, " console. The escape sequence is ^]\n"),
		resChan,
		autoResponses)

	if err != nil {
		if e, ok := err.(*websocket.CloseError); ok && e.Code == websocket.CloseAbnormalClosure {
			fmt.Fprint(os.Stderr, "\nYou were disconnected from the console. This has one of the following reasons:"+
				"\n - another user connected to the console of the target vm"+
				"\n - network issues\n")
		} else {
			fmt.Fprint(os.Stderr, "\nVirt Console Disconnection Reason: %v\n", err)
		}
		return err
	}
	return nil
}

func main() {

	fmt.Printf("ARGS %v\n", os.Args)

	run(os.Args)
}
