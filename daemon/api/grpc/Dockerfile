# FROM golang:1.22.5-alpine AS builder
#
# WORKDIR /app
#
# COPY daemon/go.mod daemon/go.sum ./
# COPY common ./common
# COPY daemon/ .
#
# RUN sed -i 's|replace obzev0/common => ../common|replace obzev0/common => ./common|' go.mod
# RUN go mod download
# RUN CGO_ENABLED=0 GOOS=linux go build -o grpc-server ./api/grpc/server.go
#
# FROM alpine:latest AS grpc_health_probe_downloader
# RUN wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.11/grpc_health_probe-linux-amd64 && \
#     chmod +x /bin/grpc_health_probe
#
# FROM alpine:latest  
#
# WORKDIR /root/
#
# COPY --from=builder /app/grpc-server .
# COPY --from=builder /app/common ./common
# COPY --from=grpc_health_probe_downloader /bin/grpc_health_probe /bin/grpc_health_probe
#
# EXPOSE 50051
# EXPOSE 2112
#
# CMD ["./grpc-server"]
#
# Stage 1: Build eBPF program and Go application
FROM golang:1.23.0-bookworm AS builder

# Install dependencies for building eBPF programs
RUN apt-get update && \
    apt-get install -y \
    clang \
    llvm \
    libelf-dev \
    linux-headers-6.1.0-22-amd64 \
    build-essential \
    iproute2 \
    libbpf-dev \
    && apt-get clean && \
    rm -rf /var/lib/apt/lists/*
	
# Set working directory
WORKDIR /app

# Copy Go modules files and download dependencies
COPY daemon/go.mod daemon/go.sum ./
COPY common ./common
COPY daemon/ .

# Modify go.mod and build Go application
RUN sed -i 's|replace obzev0/common => ../common|replace obzev0/common => ./common|' go.mod
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o grpc-server ./api/grpc/server.go

# Copy and compile eBPF program
COPY daemon/api/grpc/tcAnalyser/bpf/ /app/bpf/
WORKDIR /app/bpf
RUN clang -I/usr/include -I/usr/include/x86_64-linux-gnu -I/usr/include/x86_64-linux-gnu/bits -I/usr/include/x86_64-linux-gnu/sys -I/usr/include/bpf -O2 -g -target bpf -c tc.c -o tc.o

# Stage 2: Download grpc-health-probe
FROM debian:latest AS grpc_health_probe_downloader
RUN apt-get update && \
    apt-get install -y wget && \
    wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.11/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Stage 3: Create final image
FROM debian:latest

# Install iproute2 to handle network configuration
RUN apt-get update && \
    apt-get install -y iproute2 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /root/

# Copy compiled Go application and eBPF program
COPY --from=builder /app/grpc-server .
COPY --from=builder /app/bpf/tc.o /root/
COPY --from=grpc_health_probe_downloader /bin/grpc_health_probe /bin/grpc_health_probe

# Expose ports
EXPOSE 50051
EXPOSE 2112

# Default command to run your application
CMD ["./grpc-server"]

