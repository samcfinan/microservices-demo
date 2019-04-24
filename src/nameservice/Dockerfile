FROM golang:1.12-alpine
RUN apk add --no-cache ca-certificates git

ENV PROJECT github.com/samcfinan/microservices-demo/src/nameservice
WORKDIR /go/src/$PROJECT

# restore dependencies
ENV GO111MODULE=on
COPY . .
# COPY go.mod ./
# COPY go.sum ./
COPY ./vendor .
# RUN go mod download
# gRPC port
# RUN apk add --no-cache ca-certificates
# RUN GRPC_HEALTH_PROBE_VERSION=v0.2.0 && \
#     wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
#     chmod +x /bin/grpc_health_probe

EXPOSE 9556
RUN go build -mod vendor .
ENTRYPOINT ["./nameservice"]
