#!/bin/env sh

KUBERNETES_VERSION=1.24.7

minikube start --addons olm,registry --insecure-registry=0.0.0.0/0 --kubernetes-version $KUBERNETES_VERSION "$@"
