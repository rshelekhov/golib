package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rshelekhov/golib/observability/metrics"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
)

func ExampleInitMeter() {
	// Standard pattern from observability courses
	meterProvider, handler, err := metrics.InitMeter("my-service", "1.0.0", "production")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down: %v", err)
		}
	}()

	http.Handle("/metrics", handler)
	
	meter := metrics.OtelMeter()
	counter, _ := meter.Int64Counter("my_otel_counter", metric.WithDescription("Example otel counter."))
	counter.Add(context.Background(), 1)

	fmt.Println("Prometheus metrics available at /metrics")
}

func ExampleInitMeterOTLP() {
	// OTLP exporter for push model
	meterProvider, err := metrics.InitMeterOTLP(
		context.Background(),
		"my-service",
		"1.0.0", 
		"production",
		"localhost:4317",
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down: %v", err)
		}
	}()

	meter := metrics.OtelMeter()
	counter, _ := meter.Int64Counter("my_otlp_counter")
	counter.Add(context.Background(), 1)

	fmt.Println("Metrics pushed to OTLP collector")
}

func ExampleInitMeterStdout() {
	// Stdout exporter for development
	meterProvider, err := metrics.InitMeterStdout("my-service", "dev", "development")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down: %v", err)
		}
	}()

	meter := metrics.OtelMeter()
	counter, _ := meter.Int64Counter("my_stdout_counter")
	counter.Add(context.Background(), 1)

	fmt.Println("Metrics printed to stdout")
}

func ExampleGRPCServer() {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(metrics.UnaryServerInterceptor()),
		grpc.StreamInterceptor(metrics.StreamServerInterceptor()),
	)
	_ = server // use server
}

func main() {
	ExampleInitMeter()
	ExampleInitMeterOTLP()
	ExampleInitMeterStdout()
	ExampleGRPCServer()
}
