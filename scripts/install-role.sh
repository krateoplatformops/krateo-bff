#!/bin/bash

kubectl apply -f testdata/clusterrole-widgets-viewer.yaml
kubectl apply -f testdata/clusterrole-layout-viewer.yaml
kubectl apply -f testdata/clusterrole-formtemplates-viewer.yaml