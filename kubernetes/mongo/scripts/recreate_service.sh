#!/bin/sh
##
# Script to just deploy the MongoDB Service & StatefulSet back onto the exising Kubernetes cluster.
##

# Show persistent volume claims are still reserved even though mongod stateful-set not deployed
kubectl get --namespace=youmnibus  persistentvolumes

# Deploy just the mongodb service with mongod stateful-set only
kubectl apply --namespace=youmnibus  -f ../resources/mongodb-service.yaml
sleep 5

# Print current deployment state (unlikely to be finished yet)
kubectl get --namespace=youmnibus  all 
kubectl get --namespace=youmnibus  persistentvolumes
echo
echo "Keep running the following command until all 'mongod-n' pods are shown as running:  kubectl get --namespace=youmnibus  all"
echo

