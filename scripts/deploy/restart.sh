#!/bin/bash
kubectl delete ns spelltext
kubectl delete pv --all
kubectl delete pvc --all
helm install spelltext k8s/ -f k8s/values.yaml -n spelltext --create-namespace
kubectl config set-context --current --namespace=spelltext
kubectl get po