kubectl apply -f gateway-deployment.yml
kubectl set image gateway-deployment server=duongvu089x/gateway:$SHA --all

