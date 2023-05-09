#!/bin/bash

function apply_pod() {
	local namespace=$1
	local vminame=$2
	cat <<EOF | kubectl apply -n $namespace -f -
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: vmi-console-debug
  name: $vminame-console-logger
spec:
  containers:
  - args:
    - $vminame
    - $namespace
    stdin: true
    tty: true 
    image: quay.io/dvossel/kubevirt-console-debugger:latest
    imagePullPolicy: IfNotPresent
    name: dbug
    resources:
      requests:
        cpu: 10m
        memory: 100Mi
    securityContext:
      allowPrivilegeEscalation: false
      capabilities:
        drop:
        - ALL
  restartPolicy: Always
  securityContext:
    seccompProfile:
      type: RuntimeDefault
  serviceAccount: vmi-console-debug
  serviceAccountName: vmi-console-debug
EOF
}

function apply_rbac() {
	local namespace=$1
	cat <<EOF | kubectl apply -n $namespace -f -
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vmi-console-debug
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: vmi-console-debug
rules:
- apiGroups: ["kubevirt.io"]
  resources: ["virtualmachineinstances"]
  verbs: ["list", "get"]
- apiGroups: ["subresources.kubevirt.io"]
  resources: ["virtualmachineinstances/console"]
  verbs: ["get", "update"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: vmi-console-debug
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: vmi-console-debug
subjects:
- kind: ServiceAccount
  name: vmi-console-debug
EOF
}

while true; do
	kubectl get vm -A --no-headers |
	while read -r  line; do
		VM_NAME=$(echo $line | awk '{print $2}')
		VM_NAMESPACE=$(echo $line | awk '{print $1}')
		echo "------------------------BEGIN----------------------------------"
		echo "processing debug pod for vm $VM_NAME at namespace $VM_NAMESPACE"
		apply_rbac $VM_NAMESPACE
		apply_pod $VM_NAMESPACE $VM_NAME
		echo "-------------------------END-----------------------------------"
	done
	sleep 5
done
