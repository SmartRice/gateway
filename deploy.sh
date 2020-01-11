kubectl apply -f gateway-deployment.yaml
kubectl set image deployment.apps/gateway-deployment gateway=vuda/gateway:$SHA --all

kubectl apply -f ingress-service.yaml


