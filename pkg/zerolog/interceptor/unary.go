package interceptor

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// UnaryInterceptor is a gRPC interceptor that adds a logger to the context.
// The logger can be used to log messages related to the gRPC request.
//
// It takes a logger as a parameter and returns a grpc.UnaryServerInterceptor.
// The returned interceptor is used to intercept the gRPC unary requests.
//
// The interceptor function is called for each gRPC unary request.
// It takes the inner context, the request, the server info, and the handler.
// It returns the response and an error.
//
// The logger can be accessed using grpc.GetLogger(ctx).
func UnaryInterceptor(logger *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		innerCtx context.Context, // The context of the gRPC request.
		req interface{}, // The request object.
		_ *grpc.UnaryServerInfo, // The server info.
		handler grpc.UnaryHandler, // The handler function for the request.
	) (interface{}, error) {
		// Add the logger to the context.
		// Call the handler.
		return handler(logger.WithContext(innerCtx), req)
	}
}
