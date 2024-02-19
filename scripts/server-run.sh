#!/bin/bash

export KRATEO_BFF_DEBUG=true
export KRATEO_BFF_DUMP_ENV=true
export KRATEO_BFF_PORT=8090
export AUTHN_STORE_NAMESPACE=demo-system


kubectl apply -f crds/
kubectl apply -f testdata/ns.yaml
kubectl apply -f testdata/cardtemplate-demo.yaml
kubectl apply -f testdata/column-demo.yaml
kubectl apply -f testdata/row-demo.yaml

go run main.go -kubeconfig $HOME/.kube/config
