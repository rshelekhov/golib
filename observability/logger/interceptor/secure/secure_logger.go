package secure

import (
	"context"
	"log/slog"
	"regexp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type SecureLogger struct {
	log *slog.Logger
}

func NewSecureLogger(log *slog.Logger) *SecureLogger {
	return &SecureLogger{log: log}
}

// UnaryServerInterceptor returns a new unary server interceptor with secure logging
func (sl *SecureLogger) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Log request
		sl.log.Info("request received",
			slog.String("method", info.FullMethod),
			slog.Any("request", req),
		)

		// Call handler
		resp, err := handler(ctx, req)

		// Log response with masked sensitive data
		if err != nil {
			sl.log.Error("request failed",
				slog.String("method", info.FullMethod),
				slog.String("error", err.Error()),
				slog.String("code", status.Code(err).String()),
			)
		} else {
			maskedResp := sl.maskSensitiveData(resp)
			sl.log.Info("response sent",
				slog.String("method", info.FullMethod),
				slog.Any("response", maskedResp),
			)
		}

		return resp, err
	}
}

// maskSensitiveData masks tokens and other sensitive information in response
func (sl *SecureLogger) maskSensitiveData(resp any) any {
	if resp == nil {
		return nil
	}

	// Convert to proto message if possible
	if protoMsg, ok := resp.(proto.Message); ok {
		// Marshal to JSON
		jsonBytes, err := protojson.Marshal(protoMsg)
		if err != nil {
			return resp
		}

		jsonStr := string(jsonBytes)

		// Mask sensitive fields
		jsonStr = sl.maskJSONField(jsonStr, "accessToken")
		jsonStr = sl.maskJSONField(jsonStr, "refreshToken")
		jsonStr = sl.maskJSONField(jsonStr, "access_token")
		jsonStr = sl.maskJSONField(jsonStr, "refresh_token")
		jsonStr = sl.maskJSONField(jsonStr, "password")
		jsonStr = sl.maskJSONField(jsonStr, "token")

		return jsonStr
	}

	return resp
}

// maskJSONField masks a specific field in JSON string using regex
func (sl *SecureLogger) maskJSONField(jsonStr, fieldName string) string {
	// Regex pattern to match: "fieldName":"any_value"
	pattern := `("` + regexp.QuoteMeta(fieldName) + `"):"([^"]*?)"`
	re := regexp.MustCompile(pattern)

	// Replace with masked value
	return re.ReplaceAllString(jsonStr, `$1:"***MASKED***"`)
}
