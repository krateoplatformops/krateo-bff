#!/bin/bash

export KRATEO_BFF_DEBUG=true
export KRATEO_BFF_DUMP_ENV=true
export KRATEO_BFF_PORT=8090
export AUTHN_STORE_NAMESPACE=demo-system

# All CRDs
kubectl apply -f crds/
# Create the 'demo' namespace
kubectl apply -f testdata/ns.yaml
# CardTemplate sample
kubectl apply -f testdata/cardtemplate-demo.yaml
# Column sample
kubectl apply -f testdata/column-demo.yaml
# Row sample
kubectl apply -f testdata/row-demo.yaml
# FormDefinition sample
kubectl apply -f testdata/formdefinition.sample.yaml
# FormTemplate sample
kubectl apply -f testdata/formtemplate.sample.yaml
# Dummy 'FireworksApp' CRD (just for test/demo scopes)
kubectl apply -f testdata/fireworksapp.crd.yaml
kubectl apply -f testdata/fireworksapp.sample.yaml

go run main.go -kubeconfig $HOME/.kube/config
