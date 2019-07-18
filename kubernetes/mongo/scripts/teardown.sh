#!/bin/sh
##
# Script to remove/undepoy all project resources from the local Minikube environment.
##

# Delete mongod stateful set + mongodb service + secrets + host vm configuer daemonset
kubectl delete --namespace=youmnibus statefulsets mongod
kubectl delete --namespace=youmnibus services mongodb-service
kubectl delete --namespace=youmnibus secret shared-bootstrap-data
sleep 3

# Delete persistent volume claims
kubectl delete --namespace=youmnibus persistentvolumeclaims -l role=mongo

