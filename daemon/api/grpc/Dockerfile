FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY daemon/go.mod daemon/go.sum ./

COPY common ./common

COPY daemon/ .

RUN sed -i 's|replace obzev0/common => ../common|replace obzev0/common => ./common|' go.mod

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o grpc-server ./api/grpc/server.go

FROM alpine:latest  

WORKDIR /root/

COPY --from=builder /app/grpc-server .

COPY --from=builder /app/common ./common

EXPOSE 50051

CMD ["./grpc-server"]

