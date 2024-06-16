#!/bin/bash

set -e

kind delete cluster
kind create cluster --config kind-config.yaml

kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.14.4/config/manifests/metallb-native.yaml
{ grep -q  "controller"; kill $!; } < <(kubectl get pod -w -n metallb-system)
kubectl wait --namespace metallb-system \
                --for=condition=ready pod \
                --selector=app=metallb \
                --timeout=90s
lb_subnet=$(docker network inspect kind --format '{{json .}}' | jq -r '.IPAM.Config[0].Subnet')

# Extract the base IP and calculate the range start and end
# This example assumes the subnet is in CIDR notation e.g., 172.20.0.0/16
IFS='/' read -r base_ip cidr <<< "$lb_subnet"
IFS='.' read -r ip1 ip2 ip3 ip4 <<< "$base_ip"

echo "Using $base_ip/$_cidr subnet for MetalLB"

# Example calculation: Use the .10.1 for start and .10.50 for the end of the range
# This is a simplistic calculation that might not suit all subnets or requirements
range_start="${ip1}.${ip2}.10.1"
range_end="${ip1}.${ip2}.10.50"

# Combine them into the range string expected by MetalLB
RANGE="${range_start}-${range_end}"
export RANGE

# Prepare your YAML file with the ${RANGE} placeholder and use envsubst
envsubst < metallb-config-template.yaml | kubectl apply -f -

# Apply the YAML file
#kubectl apply -f metallb-config.yaml

helm repo add stunner https://l7mp.io/stunner
helm repo update
#disable errors
set +e
helm install stunner-gateway-operator stunner/stunner-gateway-operator --create-namespace --namespace=stunner
set -e



kubectl apply -f stunner-gateway.yaml
kubectl apply -f stunner-gwcc.yaml


kubectl apply -f cloudretro-setup-coordinator.yaml
kubectl apply -f cloudretro-setup-workers.yaml

printf "\nStart cloudretro: \n"

kubectl wait --namespace cloudretro --for=condition=available --timeout=600s deployment/coordinator-deployment deployment/worker-deployment
while ! kubectl get svc coordinator-lb-svc -n cloudretro -o jsonpath='{.status.loadBalancer.ingress[0].ip}' | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$'; do
  echo "Waiting for coordinator LoadBalancer IP..."
  sleep 10
done
echo "Coordinator LoadBalancer IP is ready."

./worker-config.sh

{ grep -q  "udp"; kill $!; } < <(kubectl get service -w -n stunner)

./coordinator-config.sh

# print the external IP
EXTERNAL_IP=$(kubectl get service -n cloudretro coordinator-lb-svc -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

printf "\nCloudRetro is ready to use. Access the coordinator at http://${EXTERNAL_IP}:8000\n"
