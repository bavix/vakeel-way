package build

import (
	"context"
	"net"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/bavix/vakeel-way/internal/app"
	way "github.com/bavix/vakeel-way/pkg/api/vakeel_way"
	"github.com/bavix/vakeel-way/pkg/zerolog/interceptor"
)

// RunGRPCServer starts a gRPC server on the TCP port specified by the `GRPCAddr`
// field of the `config` field of the `Builder` receiver. It listens on the TCP
// port specified by the `GRPCAddr` field of the `config` field of the `Builder`
// receiver. If the port is already in use, this function returns an error. It
// registers the gRPC service implementation with the gRPC server. Then it
// starts serving requests in a separate goroutine. The function blocks until the
// server is stopped or an error occurs.
//
// ctx - The context.Context used to stop the server.
// Returns an error if there is a problem with listening on the TCP port.
func (b *Builder) RunGRPCServer(ctx context.Context) error {
	// Listen on the TCP port specified by the `GRPCAddr` field of the `config`
	// field of the `Builder` receiver. If the port is already in use, an error
	// is returned.
	listen, err := net.Listen(b.config.GRPC.Network, b.config.GRPC.Addr())
	if err != nil {
		return err
	}

	// Get the logger from the context.
	logger := zerolog.Ctx(ctx)

	// Create a new gRPC server.
	server := grpc.NewServer(
		// Set the stream interceptor to add a logger to the context.
		grpc.StreamInterceptor(
			interceptor.StreamInterceptor(logger), // Add a logger to the context.
		),
		// Set the unary interceptor to add a logger to the context.
		grpc.UnaryInterceptor(
			interceptor.UnaryInterceptor(logger), // Add a logger to the context.
		),
	)

	// Start a goroutine that listens for the context to be closed. When the
	// context is closed, it closes the listener. This ensures that the server
	// is stopped when the context is closed.
	//
	// This goroutine is needed to ensure that the server is stopped when the
	// context is closed. The server is stopped by calling the Stop method on
	// the gRPC server. This method blocks until all active RPCs are finished.
	//
	// The goroutine is started after the gRPC server is started. This ensures
	// that the server is stopped after all active RPCs are finished.
	go func() {
		// Wait for the context to be closed.
		<-ctx.Done()

		// Stop the server after all active RPCs are finished. The server is
		// stopped by calling the Stop method on the gRPC server. This method
		// blocks until all active RPCs are finished.
		server.Stop()
	}()

	// Register the gRPC service implementation with the gRPC server.
	way.RegisterStateServiceServer(server, app.NewGRPCServer(b.checkerUsecase(ctx)))

	// Register reflection service on gRPC server. This allows clients to
	// discover the services and methods offered by the server.
	reflection.Register(server)

	// Start serving requests in a separate goroutine. This method blocks until
	// the server is stopped or an error occurs.

	// Log the address of the server.
	logger.Info().Str("addr", b.config.GRPC.Addr()).Msg("Starting gRPC server")

	// Start serving requests.
	return server.Serve(listen)
}
