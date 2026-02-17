TAG ?= $(shell date --utc '+"%Y.%m.%d.%H%M%S"'-local)
EXPIRE ?= 26w
WORKSPACE_RESULTS_PATH ?= /tmp/image
# Image URL to use all building/pushing image targets
IMG ?= registry.arthurvardevanyan.com/homelab/k8s-federated-credential-api:$(TAG)
export KO_DOCKER_REPO=$(shell echo $(IMG) | cut -d: -f1)
# https://catalog.redhat.com/software/containers/ubi9/ubi-micro/615bdf943f6014fa45ae1b58?architecture=amd64&image=662a8edd22c80ead7411ec6c&container-tabs=overview
export KO_DEFAULTBASEIMAGE=cgr.dev/chainguard/static

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

.PHONY: build
build:
	go build -C cmd/kfca -o /tmp/kfca

.PHONY: run
run: build
	/tmp/kfca

.PHONY: ko-build
ko-build:
	ko build ./cmd/kfca --platform=linux/amd64,linux/arm64 --bare --sbom none --image-label quay.expires-after="${EXPIRE}" --tags "${TAG}"

.PHONY: ko-build-pipeline
ko-build-pipeline:
	ko build ./cmd/kfca --platform=linux/amd64,linux/arm64 --bare --sbom none --image-label quay.expires-after="${EXPIRE}" --tags "${TAG}"
	echo "${KO_DOCKER_REPO}:${TAG}" > "${WORKSPACE_RESULTS_PATH}"
