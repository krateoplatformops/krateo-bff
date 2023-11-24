# Set the shell to bash always
SHELL := /bin/bash

# Look for a .env file, and if present, set make variables from it.
ifneq (,$(wildcard ./.env))
	include .env
	export $(shell sed 's/=.*//' .env)
endif

CLUSTER_NAME ?= local-dev
KUBECONFIG ?= $(HOME)/.kube/config

VERSION := $(shell git describe --always --tags | sed 's/-/./2' | sed 's/-/./2')
ifndef VERSION
VERSION := 0.0.0
endif

# Tools
KIND=$(shell which kind)
LINT=$(shell which golangci-lint)
KUBECTL=$(shell which kubectl)
HELM=$(shell which helm)


.DEFAULT_GOAL := help


.PHONY: test
test: ## Run all the Go test
	go test -v ./...

.PHONY: lint
lint: ## Check the Go coding conventions.
	$(LINT) run

.PHONY: tidy
tidy: ## Ensure that all Go imports are satisfied.
	go mod tidy

.PHONY: generate
generate: tidy ## Generate all CRDs.
	go generate ./...

.PHONY: kind-up
kind-up: ## Starts a KinD cluster for local development.
	@$(KIND) get kubeconfig --name $(CLUSTER_NAME) >/dev/null 2>&1 || \
		$(KIND) create cluster --name=$(CLUSTER_NAME)


.PHONY: kind-down
kind-down: ## Shuts down the KinD cluster.
	@$(KIND) delete cluster --name=$(CLUSTER_NAME)

.PHONY: kind-certs
kind-certs: ## Copy CA.crt from kind.
	rm ca.crt || true
	docker cp $(CLUSTER_NAME)-control-plane:/etc/kubernetes/pki/ca.crt ca.crt
	base64 -i ca.crt

.PHONY: demo
demo: ## Starts demo.
	$(KUBECTL) apply -f crds/
	$(KUBECTL) apply -f testdata/cardtemplate-dev.yaml
	cp $(HOME)/.kube/config kubeconfig
	go run main.go -kubeconfig kubeconfig


.PHONY: help
help: ## Print this help.
	@grep -E '^[a-zA-Z\._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'