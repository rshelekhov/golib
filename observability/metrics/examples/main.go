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
	ctx := context.Background()

	// Standard pattern using new Config API
	cfg := metrics.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "production",
		ExporterType:   metrics.ExporterPrometheus,
	}
	meterProvider, handler, err := metrics.Init(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
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
	ctx := context.Background()

	// OTLP exporter for push model
	cfg := metrics.Config{
		ServiceName:    "my-service",
		ServiceVersion: "1.0.0",
		Env:            "production",
		ExporterType:   metrics.ExporterOTLP,
		OTLPEndpoint:   "localhost:4317",
	}
	meterProvider, _, err := metrics.Init(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down: %v", err)
		}
	}()

	meter := metrics.OtelMeter()
	counter, _ := meter.Int64Counter("my_otlp_counter")
	counter.Add(context.Background(), 1)

	fmt.Println("Metrics pushed to OTLP collector")
}

// Note: Stdout exporter for metrics has been removed as it's not practical
// Use Prometheus exporter for local development instead

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
	ExampleGRPCServer()
}
