FROM golang:1.12-alpine
RUN apk add --no-cache ca-certificates git && \
      wget -qO/go/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
      chmod +x /go/bin/dep

ENV PROJECT github.com/samcfinan/microservices-demo/src/nameservice
WORKDIR /go/src/$PROJECT

# restore dependencies
ENV GO111MODULE=on
COPY . .
COPY go.mod ./
COPY go.sum ./
RUN go mod download
EXPOSE 8080
EXPOSE 9556
RUN go build .
ENTRYPOINT ["./nameservice"]
