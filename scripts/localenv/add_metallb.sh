#!/bin/bash

set -e

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

# Example calculation: Use the .10.1 for start and .10.50 for the end of the range
# This is a simplistic calculation that might not suit all subnets or requirements
range_start="${ip1}.${ip2}.10.1"
range_end="${ip1}.${ip2}.10.50"

# Combine them into the range string expected by MetalLB
RANGE="${range_start}-${range_end}"
export RANGE

# Prepare your YAML file with the ${RANGE} placeholder and use envsubst, apply the YAML file
envsubst < ./metallb-config-template.yaml | kubectl apply -f -
