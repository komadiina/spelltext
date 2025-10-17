#!/bin/bash
kubectl delete ns spelltext
kubectl delete pv --all
kubectl delete pvc --all