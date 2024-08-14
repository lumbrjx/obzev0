#!/bin/bash


echo "Creating namespace..."
if ! kubectl create namespace prometheus; then
    echo "Failed to create namespace."
    exit 1
fi

echo "Applying Prometheus config..."
if ! kubectl apply -f prometheus-config.yaml -n prometheus; then
    echo "Failed to apply Prometheus config."
    exit 1
fi
echo "Applying Prometheus Deployment..."
if ! kubectl apply -f prometheus-deployment.yaml -n prometheus; then
    echo "Failed to apply Prometheus Deployment."
    exit 1
fi
echo "Applying Prometheus Service..."
if ! kubectl apply -f prometheus-service.yaml -n prometheus; then
    echo "Failed to apply Prometheus Service."
    exit 1
fi
echo "Access Prometheus at: http://localhost:9090"
if ! kubectl port-forward service/prometheus-service 9090:80 -n prometheus; then
    echo "Failed to get Prometheus Service."
    exit 1
fi



