VALIDATE_PATH = $(GOPATH)/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.1.0
SOURCE_DIR=common/proto/$$PROTO_PATH/$$PROTO_PATH
TARGET_DIR=common/proto/$$PROTO_PATH
DEV_BRANCH = dev
STAGING_BRANCH = staging
RELEASE_BRANCH = release
MAIN_BRANCH = main

.PHONY: all generate-proto build-daemon

all: create-cluster deploy-controller deploy-daemonset setup-prometheus port-forward-prometheus 


build-cli:
	cd cmd/cli && go build -o obzevMini *.go

package-chart: 
	cd chart && helm package . 
	helm repo index . --merge index.yaml

dev-to-staging:
	@git checkout $(STAGING_BRANCH)
	@git merge $(DEV_BRANCH)
	@git push origin $(STAGING_BRANCH)
	@echo "Merged $(DEV_BRANCH) to $(STAGING_BRANCH) and pushed to origin."

staging-to-release: dev-to-staging 
	@git checkout $(RELEASE_BRANCH)
	@git merge $(STAGING_BRANCH)
	@git push origin $(RELEASE_BRANCH)
	@echo "Merged $(STAGING_BRANCH) to $(RELEASE_BRANCH) and pushed to origin."

release-to-main: staging-to-release
	@git checkout $(MAIN_BRANCH)
	@git merge $(RELEASE_BRANCH)
	@git push origin $(MAIN_BRANCH)
	@echo "Merged $(RELEASE_BRANCH) to $(MAIN_BRANCH) and pushed to origin."

# Full pipeline (runs all steps)
pipeline: release-to-main
	@echo "Pipeline completed successfully."

build-daemon:
	@if [ -z "$$TAG" ]; then \
		echo "Usage: make build-daemon TAG=<tag>"; \
		exit 1; \
	fi
	docker build -f daemon/api/grpc/Dockerfile -t lumbrjx/obzev0-grpc-daemon:$$TAG .

build-daemon-stage:
	docker build -f daemon/api/grpc/Dockerfile -t lumbrjx/obzev0-grpc-daemon:staging .
push-daemon-stage:
	docker push lumbrjx/obzev0-grpc-daemon:staging 


build-controller:
	@if [ -z "$$TAG" ]; then \
		echo "Usage: make build-controller TAG=<tag>"; \
		exit 1; \
	fi
	docker build -f controller/Dockerfile -t lumbrjx/obzev0-k8s-controller:$$TAG .

push-daemon:
	@if [ -z "$$TAG" ]; then \
		echo "Usage: make push-daemon TAG=<tag>"; \
		exit 1; \
	fi
	docker push lumbrjx/obzev0-grpc-daemon:$$TAG 

push-controller:
	@if [ -z "$$TAG" ]; then \
		echo "Usage: make push-controller TAG=<tag>"; \
		exit 1; \
	fi
	docker push lumbrjx/obzev0-k8s-controller:$$TAG 

generate-proto:
	@if [ -z "$$PROTO_PATH" ]; then \
		echo "Usage: make generate-proto PROTO_PATH=<protoDir>"; \
		exit 1; \
	fi
	protoc --proto_path=. --proto_path=$(VALIDATE_PATH) \
	       --go_out=. --go-grpc_out=. \
	       --validate_out="lang=go:common/proto/$$PROTO_PATH" \
	       --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
	       common/proto/$$PROTO_PATH/$$PROTO_PATH.proto && mv ${SOURCE_DIR}/* ${TARGET_DIR}/ && rmdir ${SOURCE_DIR}	
	@echo "Code generated successfully."

	 
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

