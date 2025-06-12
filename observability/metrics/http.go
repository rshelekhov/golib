package metrics

import (
	"net/http"
	"strconv"
	"time"
	"sync"
	"log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/attribute"
)

var (
	httpRequestsCounter metric.Int64Counter
	httpLatencyHistogram metric.Float64Histogram
	httpPanicsCounter metric.Int64Counter
	initHTTPMetricsOnce sync.Once
)

func initHTTPMetrics() {
	initHTTPMetricsOnce.Do(func() {
		meter := OtelMeter()
		var err error
		
		httpRequestsCounter, err = meter.Int64Counter(
			"http_requests_total",
			metric.WithDescription("Total number of HTTP requests."),
		)
		if err != nil {
			log.Fatalf("failed to create http_requests_total counter: %v", err)
		}
		
		httpLatencyHistogram, err = meter.Float64Histogram(
			"http_request_duration_seconds",
			metric.WithDescription("HTTP request latency in seconds."),
		)
		if err != nil {
			log.Fatalf("failed to create http_request_duration_seconds histogram: %v", err)
		}
		
		httpPanicsCounter, err = meter.Int64Counter(
			"http_panics_total",
			metric.WithDescription("Total number of panics in HTTP handlers."),
		)
		if err != nil {
			log.Fatalf("failed to create http_panics_total counter: %v", err)
		}
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Middleware returns http.Handler with otel-metrics
func Middleware(next http.Handler) http.Handler {
	initHTTPMetrics()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		ctx := r.Context()
		defer func() {
			if rec := recover(); rec != nil {
				httpPanicsCounter.Add(ctx, 1, metric.WithAttributes(
					attribute.String("method", r.Method),
					attribute.String("path", r.URL.Path),
				))
				panic(rec)
			}
		}()
		next.ServeHTTP(rec, r)
		status := strconv.Itoa(rec.status)
		httpRequestsCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("method", r.Method),
			attribute.String("path", r.URL.Path),
			attribute.String("status", status),
		))
		httpLatencyHistogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
			attribute.String("method", r.Method),
			attribute.String("path", r.URL.Path),
		))
	})
}
