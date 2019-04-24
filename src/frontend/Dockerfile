FROM golang:1.12-alpine
RUN apk add --no-cache ca-certificates git
ENV PROJECT github.com/samcfinan/microservices-demo/src/frontend
WORKDIR /go/src/$PROJECT

# restore dependencies
ENV GO111MODULE=on
COPY . .
COPY go.mod ./
COPY go.sum ./
# RUN go mod download
COPY ./templates ./templates
COPY ./static ./static
COPY ./vendor ./vendor
EXPOSE 8080
RUN go build -mod vendor .
ENTRYPOINT ["./frontend"]
