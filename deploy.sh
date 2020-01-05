kubectl apply -f gateway-deployment.yaml
kubectl set image deployment.apps/gateway-deployment gateway=duongvu089x/gateway:$SHA --all

