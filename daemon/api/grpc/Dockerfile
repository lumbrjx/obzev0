FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY daemon/go.mod daemon/go.sum ./
COPY common ./common
COPY daemon/ .

RUN sed -i 's|replace obzev0/common => ../common|replace obzev0/common => ./common|' go.mod
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o grpc-server ./api/grpc/server.go

FROM alpine:latest AS grpc_health_probe_downloader
RUN wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.11/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

FROM alpine:latest  

WORKDIR /root/

COPY --from=builder /app/grpc-server .
COPY --from=builder /app/common ./common
COPY --from=grpc_health_probe_downloader /bin/grpc_health_probe /bin/grpc_health_probe

EXPOSE 50051
EXPOSE 2112

CMD ["./grpc-server"]

