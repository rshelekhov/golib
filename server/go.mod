module github.com/rshelekhov/golib/server

go 1.24.2

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1
	github.com/rshelekhov/golib/middleware/cors v0.0.0
	github.com/rshelekhov/golib/middleware/logging v0.0.0
	github.com/rshelekhov/golib/middleware/recovery v0.0.0
	golang.org/x/sync v0.16.0
	google.golang.org/grpc v1.74.2
)

replace (
	github.com/rshelekhov/golib/middleware/cors => ../middleware/cors
	github.com/rshelekhov/golib/middleware/logging => ../middleware/logging
	github.com/rshelekhov/golib/middleware/recovery => ../middleware/recovery
	github.com/rshelekhov/golib/middleware/validation => ../middleware/validation
)

require (
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk v1.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.37.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250721164621-a45f3dfb1074 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250721164621-a45f3dfb1074 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
