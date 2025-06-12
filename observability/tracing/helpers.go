package tracing

import "go.opentelemetry.io/otel/attribute"

// Helpers for creating attributes
func String(key, value string) Attribute {
	return attribute.String(key, value)
}

func Int(key string, value int) Attribute {
	return attribute.Int(key, value)
}

func Bool(key string, value bool) Attribute {
	return attribute.Bool(key, value)
}
