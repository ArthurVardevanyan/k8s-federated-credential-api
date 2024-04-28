
export KO_DOCKER_REPO=registry.arthurvardevanyan.com/homelab/k8s-federated-credential-api
# https://catalog.redhat.com/software/containers/ubi9/ubi-minimal/615bd9b4075b022acc111bf5?architecture=amd64&image=65e0932f034203e025b55a92
export KO_DEFAULTBASEIMAGE=registry.access.redhat.com/ubi9-minimal:9.3-1612
TAG ?= $(shell date --utc '+%Y%m%d-%H%M')
EXPIRE ?= 1d

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec


.PHONY: run
run:
	go build -C cmd/kfca -o /tmp/kfca && /tmp/kfca

.PHONY: ko-build
ko-build:
	ko build ./cmd/kfca --platform=linux/amd64 --bare --sbom none --image-label quay.expires-after="${EXPIRE}" --tags "${TAG}"
