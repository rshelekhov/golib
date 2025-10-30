package requestid

// Constants for request ID handling
const (
	// Header is the HTTP/gRPC metadata header name for request ID
	Header = "X-Request-ID"

	// CtxKey is the context key used to store request ID
	CtxKey = "RequestID"
)
