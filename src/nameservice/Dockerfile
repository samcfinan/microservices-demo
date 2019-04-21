FROM golang:1.12-alpine
RUN apk add --no-cache ca-certificates git && \
      wget -qO/go/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
      chmod +x /go/bin/dep

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
EXPOSE 9556
RUN go build -mod vendor .
ENTRYPOINT ["./nameservice"]
