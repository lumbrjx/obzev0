# Build the manager binary
FROM golang:1.22 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

# Copy the Go Modules manifests
COPY controller/go.mod go.mod
COPY controller/go.sum go.sum

# Copy the `common` directory to the Docker image
COPY common ./common

# Ensure the `replace` directive in go.mod points to the correct path
RUN sed -i 's|replace obzev0/common => ../common|replace obzev0/common => ./common|' go.mod

# Download Go modules
RUN go mod download

# Copy the source code into the Docker image
COPY controller/cmd/main.go cmd/main.go
COPY controller/api/ api/
COPY controller/internal/controller/ internal/controller/

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go

# Use a minimal base image for the final image
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]

