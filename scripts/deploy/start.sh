#!/bin/bash

helm install spelltext k8s/ -f k8s/values.yaml -n spelltext --create-namespace
kubectl config set-context --current --namespace=spelltext
kubectl get po