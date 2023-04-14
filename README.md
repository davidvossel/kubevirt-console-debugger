This is a debug tool and not meant for production.

The kubevirt-console-logger creates a pod that streams the console of a vm to stdout. The vm's console logs are then accessible by viewing the pod's logs. 

# Install

```
kubectl apply -f kubevirt-console-logger.yaml
```

Then verify the console logger is online. This deployment watches for VMs to come online and creates pods to stream the console logs.
```
kubectl get pods -n vm-logger
NAME                                         READY   STATUS    RESTARTS   AGE
vm-console-debug-launcher-5879d67756-gfz47   1/1     Running   0          2m42s
```

# Example

create a vm called my-vm in the default namespace

```
kubectl create my-vm.yaml
```

Watch the VM and the console logger come online
```
kubectl get pods
NAME                            READY   STATUS    RESTARTS   AGE
virt-launcher-my-vm-7fjpz   2/2     Running   0          56s
my-vm-console-logger        1/1     Running   0          2m16s
```

View the VM's console history at any point by streaming the logger's container logs
```
kubectl logs my-vm-console-logger
[    0.000000] Initializing cgroup subsys cpuset
[    0.000000] Initializing cgroup subsys cpu
[    0.000000] Initializing cgroup subsys cpuacct
[    0.000000] Linux version 4.4.0-28-generic (buildd@lcy01-13) (gcc version 5.3.1 20160413 (Ubuntu 5.3.1-14ubuntu2.1) ) #47-Ubuntu SMP Fri Jun 24 10:09:13 UTC 2016 (Ubuntu 4.4.0-28.47-generic 4.4.13)
[    0.000000] Command line: LABEL=cirros-rootfs ro console=tty1 console=ttyS0
[    0.000000] KERNEL supported cpus:
[    0.000000]   Intel GenuineIntel
[    0.000000]   AMD AuthenticAMD
[    0.000000]   Centaur CentaurHauls
[    0.000000] x86/fpu: xstate_offset[2]:  576, xstate_sizes[2]:  256
[    0.000000] x86/fpu: xstate_offset[5]: 1088, xstate_sizes[5]:   64
[    0.000000] x86/fpu: xstate_offset[6]: 1152, xstate_sizes[6]:  512
[    0.000000] x86/fpu: xstate_offset[7]: 1664, xstate_sizes[7]: 1024
[    0.000000] x86/fpu: Supporting XSAVE feature 0x01: 'x87 floating point registers'
[    0.000000] x86/fpu: Supporting XSAVE feature 0x02: 'SSE registers'
[    0.000000] x86/fpu: Supporting XSAVE feature 0x04: 'AVX registers'
[    0.000000] x86/fpu: Supporting XSAVE feature 0x20: 'AVX-512 opmask'
[    0.000000] x86/fpu: Supporting XSAVE feature 0x40: 'AVX-512 Hi256'
[    0.000000] x86/fpu: Supporting XSAVE feature 0x80: 'AVX-512 ZMM_Hi256'
[    0.000000] x86/fpu: Enabled xstate features 0xe7, context size is 2688 bytes, using 'standard' format.
[    0.000000] x86/fpu: Using 'eager' FPU context switches.
[    0.000000] e820: BIOS-provided physical RAM map:
[    0.000000] BIOS-e820: [mem 0x0000000000000000-0x000000000009fbff] usable
[    0.000000] BIOS-e820: [mem 0x000000000009fc00-0x000000000009ffff] reserved
[    0.000000] BIOS-e820: [mem 0x00000000000f0000-0x00000000000fffff] reserved
[    0.000000] BIOS-e820: [mem 0x0000000000100000-0x0000000007fdcfff] usable
[    0.000000] BIOS-e820: [mem 0x0000000007fdd000-0x0000000007ffffff] reserved
[    0.000000] BIOS-e820: [mem 0x00000000b0000000-0x00000000bfffffff] reserved
[    0.000000] BIOS-e820: [mem 0x00000000fed1c000-0x00000000fed1ffff] reserved
.
.
.
```

