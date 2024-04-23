
REGISTRY ?= quay.io/dvossel
TAG ?= latest
VM_CONSOLE_DBUG_IMAGE_NAME ?= kubevirt-console-debugger
VM_CONSOLE_DBUG_IMG ?= $(REGISTRY)/$(VM_CONSOLE_DBUG_IMAGE_NAME)

.PHONY: build
build:
	go build main.go

.PHONY: docker-build
docker-build:
	docker build . -t $(VM_CONSOLE_DBUG_IMG):$(TAG) --file Dockerfile

.PHONY: docker-push
docker-push:
	docker push $(VM_CONSOLE_DBUG_IMG):$(TAG)
