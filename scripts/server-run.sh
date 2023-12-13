#!/bin/bash

export KRATEO_BFF_DEBUG=true

go run main.go -debug -kubeconfig $HOME/.kube/config
