package interceptor

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// serverStreamWrapper is a wrapper around a gRPC server stream.
//
// It adds a context.Context field to the wrapper. This field is used to
// store the context.Context object that is used to log messages related
// to the gRPC stream.
type serverStreamWrapper struct {
	ss  grpc.ServerStream
	ctx context.Context //nolint:containedctx
}

// Context returns the context.Context object stored in the wrapper.
//
// It is used to get the context.Context object that is used to log messages
// related to the gRPC stream.
func (w serverStreamWrapper) Context() context.Context {
	return w.ctx
}

// RecvMsg receives a message from the stream.
//
// It is used to receive a message from the stream and log messages related
// to the message.
func (w serverStreamWrapper) RecvMsg(msg interface{}) error {
	return w.ss.RecvMsg(msg)
}

// SendMsg sends a message to the stream.
//
// It is used to send a message to the stream and log messages related
// to the message.
func (w serverStreamWrapper) SendMsg(msg interface{}) error {
	return w.ss.SendMsg(msg)
}

// SendHeader sends a metadata header to the stream.
//
// It is used to send a metadata header to the stream and log messages
// related to the header.
func (w serverStreamWrapper) SendHeader(md metadata.MD) error {
	return w.ss.SendHeader(md)
}

// SetHeader sets a metadata header on the stream.
//
// It is used to set a metadata header on the stream and log messages
// related to the header.
func (w serverStreamWrapper) SetHeader(md metadata.MD) error {
	return w.ss.SetHeader(md)
}

// SetTrailer sets a metadata trailer on the stream.
//
// It is used to set a metadata trailer on the stream and log messages
// related to the trailer.
func (w serverStreamWrapper) SetTrailer(md metadata.MD) {
	w.ss.SetTrailer(md)
}

// StreamInterceptor is a gRPC interceptor that adds a logger to the context.
// The logger can be used to log messages related to the gRPC stream.
//
// It takes a logger as a parameter and returns a grpc.StreamServerInterceptor.
// The returned interceptor is used to intercept the gRPC stream requests.
//
// The interceptor function is called for each gRPC stream request.
// It takes the server, the stream, the server info, and the handler.
// It returns an error.
func StreamInterceptor(logger *zerolog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{}, // The server object.
		ss grpc.ServerStream, // The stream object.
		_ *grpc.StreamServerInfo, // The server info.
		handler grpc.StreamHandler, // The handler function for the stream.
	) error {
		// Create a serverStreamWrapper object with the stream and context.
		// The context is created with the logger.
		//
		// It is used to log messages related to the gRPC stream.
		return handler(srv, serverStreamWrapper{
			ss:  ss,
			ctx: logger.WithContext(ss.Context()),
		})
	}
}
