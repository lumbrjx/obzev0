#!/bin/bash

echo "Undeploying the Controller..."
if ! make undeploy; then
    echo "Failed to undepoly the Controller."
    exit 1
fi
echo "Deleting the role binding..."
if ! kubectl delete clusterrolebinding permissive-binding; then
    echo "Failed to delete the role binding."
    exit 1
fi

echo "Clean up completed successfully."

