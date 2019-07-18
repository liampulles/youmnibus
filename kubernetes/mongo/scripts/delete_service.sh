#!/bin/sh
##
# Script to just undeploy the MongoDB Service & StatefulSet but nothing else.
##

# Just delete mongod stateful set + mongodb service onlys (keep rest of k8s environment in place)
kubectl delete --namespace=youmnibus  statefulsets mongod
kubectl delete --namespace=youmnibus  services mongodb-service

# Show persistent volume claims are still reserved even though mongod stateful-set has been undeployed
kubectl get --namespace=youmnibus persistentvolumes

