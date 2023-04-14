FROM quay.io/fedora/fedora:38

RUN cd /usr/bin && \
	curl -L https://github.com/kubevirt/kubevirt/releases/download/v0.59.0/virtctl-v0.59.0-linux-amd64 --output virtctl && \
	chmod 755 virtctl && \
	curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
	chmod 755 kubectl && \
	mkdir /manifests

COPY console-launcher.yaml /manifests/console-launcher.yaml
COPY entrypoint.sh /usr/bin/entrypoint.sh
COPY vm-console-dbug-launcher.sh /usr/bin/vm-console-dbug-launcher.sh

ENTRYPOINT ["/usr/bin/entrypoint.sh"]

