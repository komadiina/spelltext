@echo off
echo "[INFO] starting minikube tunnel..."
start cmd /c "minikube tunnel"

echo "[INFO] installing helm charts..."
helm install spelltext chart/ -n spelltext --create-namespace

echo "[INFO] setting current kubectl context..."
kubectl config set-context --current --namespace=spelltext

echo "[INFO] port-forwarding NATS (client port, 4222@TCP) for clients..."
start cmd /c "kubectl port-forward nats-spelltext-0 4222:4222"

echo "[INFO] spelltext server started!"