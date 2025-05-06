
WORKSPACE_RESULTS_PATH ?= /tmp/image
export KO_DOCKER_REPO=registry.arthurvardevanyan.com/homelab/k8s-federated-credential-api
# https://catalog.redhat.com/software/containers/ubi9/ubi-minimal/615bd9b4075b022acc111bf5?architecture=amd64&image=66cddd84df3259c57ceb8f65
export KO_DEFAULTBASEIMAGE=cgr.dev/chainguard/static:latest
TAG ?= $(shell date --utc '+%Y%m%d-%H%M')
EXPIRE ?= 180d

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

.PHONY: goa-gen
goa-gen:
	~/go/bin/goa gen k8s-federated-credential-api/design && \
	jq -c '.components.schemas.ExchangeTokenRequestBody += {"additionalProperties":false}' gen/http/openapi3.json > gen/http/tmp.json && mv gen/http/tmp.json  gen/http/openapi3.json
	yq -i '.components.schemas.ExchangeTokenRequestBody.additionalProperties += false' gen/http/openapi3.yaml
	jq -c '.components.schemas.StatusResult += {"additionalProperties":false}' gen/http/openapi3.json > gen/http/tmp.json && mv gen/http/tmp.json  gen/http/openapi3.json
	yq -i '.components.schemas.StatusResult.additionalProperties += false' gen/http/openapi3.yaml
	jq -c '.components.schemas.Status += {"additionalProperties":false}' gen/http/openapi3.json > gen/http/tmp.json && mv gen/http/tmp.json  gen/http/openapi3.json
	yq -i '.components.schemas.Status.additionalProperties += false' gen/http/openapi3.yaml
	jq -c '.components.schemas.Error += {"additionalProperties":false}' gen/http/openapi3.json > gen/http/tmp.json && mv gen/http/tmp.json  gen/http/openapi3.json
	yq -i '.components.schemas.Error.additionalProperties += false' gen/http/openapi3.yaml

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
