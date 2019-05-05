module github.com/samcfinan/microservices-demo/src/api

go 1.12

require (
	cloud.google.com/go v0.37.2
	contrib.go.opencensus.io/exporter/stackdriver v0.5.0
	git.apache.org/thrift.git v0.12.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/golang/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.1
	github.com/grpc-ecosystem/grpc-gateway v1.6.2 // indirect
	github.com/kisielk/gotool v1.0.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/samcfinan/microservices-demo/src/frontend v0.0.0-20190424020159-6c03bf4ac1d8
	github.com/samcfinan/microservices-demo/src/nameservice v0.0.0-20190424140046-2c2b6f0822af
	github.com/sirupsen/logrus v1.3.0
	go.opencensus.io v0.19.2
	google.golang.org/grpc v1.19.0
)
