
REGISTRY ?= quay.io/dvossel
TAG ?= latest
RHCOS_TAG ?= rhcos
VM_CONSOLE_DBUG_IMAGE_NAME ?= kubevirt-console-debugger
VM_CONSOLE_DBUG_IMG ?= $(REGISTRY)/$(VM_CONSOLE_DBUG_IMAGE_NAME)

.PHONY: build
build:
	go build .

.PHONY: docker-build
docker-build:
	docker build . -t $(VM_CONSOLE_DBUG_IMG):$(TAG) --file Dockerfile

.PHONY: docker-push
docker-push:
	docker push $(VM_CONSOLE_DBUG_IMG):$(TAG)

.PHONY: docker-build-rhcos
docker-build-rhcos:
	docker build . -t $(VM_CONSOLE_DBUG_IMG):$(RHCOS_TAG) --file Dockerfile.rhcos-debug

.PHONY: docker-push-rhcos
docker-push-rhcos:
	docker push $(VM_CONSOLE_DBUG_IMG):$(RHCOS_TAG)
