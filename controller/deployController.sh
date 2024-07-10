#!/bin/bash

if [ $# -ne 1 ]; then
    echo "Usage: $0 <controller-image-tag>"
    exit 1
fi

IMAGE_TAG="$1"

echo "Installing CRDs into the cluster..."
if ! make install; then
    echo "Failed to install CRDs."
    exit 1
fi

echo "Deploying the controller image into the cluster..."
if ! make deploy IMG=lumbrjx/obzev0poc:${IMAGE_TAG}; then
    echo "Failed to deploy the controller image."
    exit 1
fi

echo "Granting access and permissions..."

for config in controller-manager-clusterrole.yaml controller-manager-rolebinding.yaml; do
    if ! kubectl apply -f "$config"; then
        echo "Failed to apply $config."
        exit 1
    fi
done

if ! kubectl create clusterrolebinding permissive-binding \
    --clusterrole=cluster-admin \
    --serviceaccount=controller-system:controller-controller-manager; then
    echo "Failed to create the clusterrolebinding."
    exit 1
fi

echo "Deployment completed successfully."

