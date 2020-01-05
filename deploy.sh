kubectl apply -f ./gateway-deployment.yaml
kubectl set image gateway-deployment server=duongvu089x/gateway:$SHA --all

