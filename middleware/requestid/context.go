package requestid

import "context"

// FromContext extracts the request ID from the context
func FromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(CtxKey).(string)
	return requestID, ok
}

// WithContext adds the request ID to the context
func WithContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, CtxKey, requestID)
}
