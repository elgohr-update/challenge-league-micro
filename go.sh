set -euxo pipefail
bash docker-build.sh
kubectl delete -f kubernetes/
kubectl apply -f kubernetes/
