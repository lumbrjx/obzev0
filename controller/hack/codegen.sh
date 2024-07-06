#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
PROJECT_ROOT=$(realpath "${SCRIPT_ROOT}")

# Find the code-generator dir in the go mod cache
CODEGEN_PKG=$(go env GOMODCACHE)/k8s.io/code-generator@v0.26.0

# Ensure the package exists
if [ ! -d "${CODEGEN_PKG}" ]; then
    echo "k8s.io/code-generator package not found. Running go mod download..."
    go mod download k8s.io/code-generator
fi

# Verify again
if [ ! -d "${CODEGEN_PKG}" ]; then
    echo "Error: k8s.io/code-generator package still not found after attempting to download."
    exit 1
fi

echo "Generating code..."
echo "SCRIPT_ROOT: ${SCRIPT_ROOT}"
echo "PROJECT_ROOT: ${PROJECT_ROOT}"
echo "CODEGEN_PKG: ${CODEGEN_PKG}"

bash "${CODEGEN_PKG}"/generate-groups.sh all \
  github.com/lumbrjx/obzev0/pkg/client github.com/lumbrjx/obzev0/api \
  batch.github.com:v1 \
  --output-base "${PROJECT_ROOT}" \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt"

# Move generated files to the correct location
mkdir -p "${PROJECT_ROOT}/pkg"
mv "${PROJECT_ROOT}/github.com/lumbrjx/obzev0/pkg/client" "${PROJECT_ROOT}/pkg/"

echo "Code generation complete. Generated files are in: ${PROJECT_ROOT}/pkg/client"
