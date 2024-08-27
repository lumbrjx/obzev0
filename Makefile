.PHONY: all

all: create-cluster deploy-controller deploy-daemonset setup-prometheus port-forward-prometheus

deploy-controller:
	@if [ -z "$$IMAGE_TAG" ]; then \
		echo "Usage: make deploy-controller IMAGE_TAG=<controller-image-tag>"; \
		exit 1; \
	fi
	cd controller && ./deployController.sh $$IMAGE_TAG
	@echo "Deployment completed successfully.###################################################################"

deploy-daemonset:
	@echo "Applying the DaemonSet..."
	@if ! kubectl apply -f daemonset.yaml; then \
		echo "Failed to apply the DaemonSet."; \
		exit 1; \
	fi
	@echo "Applying the Service..."
	@if ! kubectl apply -f service.yaml; then \
		echo "Failed to apply the Service."; \
		exit 1; \
	fi
	@echo "Labeling nodes..."
	@for node in kind-worker kind-worker2 kind-worker3; do \
		if ! kubectl label node $$node node-role.kubernetes.io/worker=; then \
			echo "Failed to label node $$node."; \
			exit 1; \
		fi \
	done
	@echo "Deployment completed successfully.###################################################################"

cleanup-controller:
	@echo "Undeploying the Controller..."
	@if ! cd controller && ./cleanUpController.sh; then \
		echo "Failed to undeploy the Controller."; \
		exit 1; \
	fi
	@echo "Clean up completed successfully.###################################################################"

cleanup-cluster:
	@echo "Deleting cluster..."
	@if ! kind delete cluster; then \
		echo "Failed to delete the cluster."; \
		exit 1; \
	fi


KIND_CONFIG=kind.yaml

create-cluster:
	kind create cluster --config $(KIND_CONFIG)

setup-prometheus:
	@echo "Creating namespace..."
	@if ! kubectl create namespace prometheus; then \
		echo "Failed to create namespace."; \
		exit 1; \
	fi
	@echo "Applying Prometheus config..."
	@if ! kubectl apply -f prometheus-config.yaml -n prometheus; then \
		echo "Failed to apply Prometheus config."; \
		exit 1; \
	fi
	@echo "Applying Prometheus Deployment..."
	@if ! kubectl apply -f prometheus-deployment.yaml -n prometheus; then \
		echo "Failed to apply Prometheus Deployment."; \
		exit 1; \
	fi
	@echo "Applying Prometheus Service..."
	@if ! kubectl apply -f prometheus-service.yaml -n prometheus; then \
		echo "Failed to apply Prometheus Service."; \
		exit 1; \
	fi
	@echo "###################################################################"
port-forward-prometheus:
	@echo "Access Prometheus at: http://localhost:9090"
	kubectl wait --for=condition=ready pod -l app=prometheus-server -n prometheus --timeout=120s
	@if ! kubectl port-forward service/prometheus-service 9090:80 -n prometheus; then \
		echo "Failed to get Prometheus Service."; \
		exit 1; \
	fi

