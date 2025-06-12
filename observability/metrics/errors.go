package metrics

import (
	"sync"
	"log"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	businessErrorsCounter metric.Int64Counter
	initBusinessErrorsOnce sync.Once
)

func initBusinessErrors() {
	initBusinessErrorsOnce.Do(func() {
		meter := OtelMeter()
		var err error
		
		businessErrorsCounter, err = meter.Int64Counter(
			"business_errors_total",
			metric.WithDescription("Total number of business logic errors."),
		)
		if err != nil {
			log.Fatalf("failed to create business_errors_total counter: %v", err)
		}
	})
}

// IncBusinessError increases business error counter
func IncBusinessError(errType, code string) {
	initBusinessErrors()
	businessErrorsCounter.Add(context.Background(), 1, metric.WithAttributes(
		attribute.String("type", errType),
		attribute.String("code", code),
	))
}
