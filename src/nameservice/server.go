package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	pb "github.com/samcfinan/microservices-demo/src/nameservice/genproto"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
)

var (
	cat          pb.ListProductsResponse
	log          *logrus.Logger
	extraLatency time.Duration

	port = flag.Int("port", 9556, "port to listen at")

	db    *gorm.DB
	dbErr error
)

func init() {
	log = logrus.New()
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
}

func main() {
	// go initTracing()
	// go initProfiling("nameservice", "1.0.0")
	flag.Parse()

	db, dbErr = gorm.Open("postgres", fmt.Sprintf("host=postgres port=%s user=%s dbname=%s password=%s sslmode=disable", os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASS")))
	defer db.Close()
	if dbErr != nil {
		log.Panic("Cannot connect to Postgres:", dbErr)
	}

	log.Infof("starting grpc server at :%d", *port)
	run(*port)

	// Blocking
	select {}
}

func run(port int) string {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	svc := &nameServer{}
	pb.RegisterNameServiceServer(srv, svc)
	// healthpb.RegisterHealthServer(srv, svc)
	go srv.Serve(l)
	return l.Addr().String()
}

func initJaegerTracing() {
	svcAddr := os.Getenv("JAEGER_SERVICE_ADDR")
	if svcAddr == "" {
		log.Info("jaeger initialization disabled.")
		return
	}
	// Register the Jaeger exporter to be able to retrieve
	// the collected spans.
	exporter, err := jaeger.NewExporter(jaeger.Options{
		Endpoint: fmt.Sprintf("http://%s", svcAddr),
		Process: jaeger.Process{
			ServiceName: "nameservice",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	log.Info("jaeger initialization completed.")
}

func initStats(exporter *stackdriver.Exporter) {
	view.SetReportingPeriod(60 * time.Second)
	view.RegisterExporter(exporter)
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Info("Error registering default server views")
	} else {
		log.Info("Registered default server views")
	}
}

func initStackDriverTracing() {
	for i := 1; i <= 3; i++ {
		exporter, err := stackdriver.NewExporter(stackdriver.Options{})
		if err != nil {
			log.Warnf("failed to initialize stackdriver exporter: %+v", err)
		} else {
			trace.RegisterExporter(exporter)
			trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
			log.Info("registered stackdriver tracing")

			// Register the views to collect server stats.
			initStats(exporter)
			return
		}
		d := time.Second * 10 * time.Duration(i)
		log.Infof("sleeping %v to retry initializing stackdriver exporter", d)
		time.Sleep(d)
	}
	log.Warn("could not initialize stackdriver exporter after retrying, giving up")
}

func initTracing() {
	initJaegerTracing()
	initStackDriverTracing()
}

func initProfiling(service, version string) {
	for i := 1; i <= 3; i++ {
		if err := profiler.Start(profiler.Config{
			Service:        service,
			ServiceVersion: version,
			// ProjectID must be set if not running on GCP.
			// ProjectID: "my-project",
		}); err != nil {
			log.Warnf("failed to start profiler: %+v", err)
		} else {
			log.Info("started stackdriver profiler")
			return
		}
		d := time.Second * 10 * time.Duration(i)
		log.Infof("sleeping %v to retry initializing stackdriver profiler", d)
		time.Sleep(d)
	}
	log.Warn("could not initialize stackdriver profiler after retrying, giving up")
}

type nameServer struct{}

type Name struct {
	Name   string
	Length int32
}

func (n *nameServer) CheckName(ctx context.Context, nr *pb.NameRequest) (*pb.NameResponse, error) {
	name := nr.GetName()
	nameLength := int32(len(name))
	na := Name{Name: name, Length: nameLength}
	db.Create(&na)
	return &pb.NameResponse{Name: name, NameLength: nameLength}, nil
}

func (n *nameServer) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (n *nameServer) Watch(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}
