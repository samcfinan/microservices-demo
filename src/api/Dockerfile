FROM golang:1.12-alpine
RUN apk add --no-cache ca-certificates git
ENV PROJECT github.com/samcfinan/microservices-demo/src/api
WORKDIR /go/src/$PROJECT

# restore dependencies
ENV GO111MODULE=on
COPY . .
COPY go.mod ./
COPY go.sum ./
# RUN go mod download
COPY ./vendor ./vendor
EXPOSE 8080
RUN go build -mod vendor .
ENTRYPOINT ["./api"]
