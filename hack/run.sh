#!/bin/bash

kind get kubeconfig >/dev/null 2>&1 || kind create cluster

kubectl apply -f crds/
kubectl apply -f testdata/cardtemplate-dev.yaml

go run main.go -debug -kubeconfig $HOME/.kube/config
