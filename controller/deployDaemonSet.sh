#!/bin/bash

echo "Applying the DaemonSet..."
if ! kubectl apply -f daemonset.yaml; then
    echo "Failed to apply the DaemonSet."
    exit 1
fi
echo "Applying the Service..."
if ! kubectl apply -f service.yaml; then
    echo "Failed to apply the Service."
    exit 1
fi

echo "Labeling nodes..."
if ! kubectl label node kind-worker node-role.kubernetes.io/worker=; then
    echo "Failed to label node."
    exit 1
fi
if ! kubectl label node kind-worker2 node-role.kubernetes.io/worker=; then
    echo "Failed to label node."
    exit 1
fi
if ! kubectl label node kind-worker3 node-role.kubernetes.io/worker=; then
    echo "Failed to label node."
    exit 1
fi


echo "Deployment completed successfully."

