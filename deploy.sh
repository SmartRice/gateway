kubectl apply -f k8s
kubectl set image deployments/gateway-deployment server=duongvu089x/gateway:$SHA

