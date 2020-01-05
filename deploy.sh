kubectl apply -f gateway-deployment.yaml
kubectl set image deployment.apps/gateway-deployment server=duongvu089x/gateway:$SHA --all

