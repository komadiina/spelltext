@echo off
echo "[INFO] stopping spelltext..."
echo "[INFO] deleting helm charts..."
helm uninstall spelltext

echo "[INFO] deleting kubectl namespace 'spelltext'..."
kubectl delete ns spelltext

echo "[INFO] stopped!"