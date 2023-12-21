#!/bin/bash

kubectl apply -f crds/
kubectl apply -f testdata/cardtemplate-sample.yaml
