
#!/bin/bash

while true; do
	date -u
	echo "attempting to connect to vmi $1 at namespace $2 streaming console log to stdout"
	/usr/bin/kubevirt-console-debugger $1 $2
	sleep 2
done

