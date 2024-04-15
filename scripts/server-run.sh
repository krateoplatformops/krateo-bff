#!/bin/bash

# Startup Kind
kind get kubeconfig >/dev/null 2>&1 || kind create cluster

export KRATEO_BFF_DEBUG=true
export KRATEO_BFF_DUMP_ENV=true
export KRATEO_BFF_PORT=8090
export AUTHN_STORE_NAMESPACE=demo-system

# All CRDs
kubectl apply -f crds/

# Create the 'demo' namespace
kubectl apply -f testdata/ns.yaml
# CardTemplate sample
#kubectl apply -f testdata/cardtemplate-demo.yaml
# Column sample
#kubectl apply -f testdata/column-demo.yaml
# Row sample
#kubectl apply -f testdata/row-demo.yaml

# FormTemplate sample
#kubectl apply -f testdata/formtemplate.sample.yaml
# Dummy 'FireworksApp' CRD (just for test/demo scopes)
#kubectl apply -f testdata/fireworksapp.crd.yaml
#kubectl apply -f testdata/fireworksapp.sample.yaml

# Install roles
kubectl apply -f testdata/clusterrole-widgets-viewer.yaml
kubectl apply -f testdata/clusterrole-layout-viewer.yaml
kubectl apply -f testdata/clusterrole-formtemplates-viewer.yaml
kubectl apply -f testdata/clusterrole-apps-viewer.yaml

# Issue 20240415
kubectl apply -f testdata/issue-20240415.yaml

go run main.go -kubeconfig $HOME/.kube/config
