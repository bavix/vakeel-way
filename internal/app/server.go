package app

import (
	"github.com/bavix/apis/pkg/uuidconv"
	"github.com/bavix/vakeel-way/internal/domain/usecases"
	way "github.com/bavix/vakeel-way/pkg/api/vakeel_way"
)

var _ = way.StateServiceServer(&GRPCServer{})

// NewGRPCServer creates a new instance of the GRPCServer struct.
//
// It takes a *usecases.Checker as a parameter and returns a pointer to a GRPCServer struct.
// The GRPCServer struct implements the way.StateServiceServer interface and is used to provide the StateService
// RPC service. The checker parameter is used to send events to the checker.
//
// Parameters:
//   - checker: A *usecases.Checker used to send events to the checker.
//
// Returns:
//   - A pointer to a GRPCServer struct.
func NewGRPCServer(
	checker *usecases.Checker,
) *GRPCServer {
	// Create a new instance of the GRPCServer struct.
	// The GRPCServer struct implements the way.StateServiceServer interface and is used to provide the StateService
	// RPC service. The checker parameter is used to send events to the checker.
	return &GRPCServer{
		// The checker field is used to send events to the checker.
		checker: checker,
	}
}

// GRPCServer is a gRPC server implementation that provides the StateService
// RPC service. It implements the way.StateServiceServer interface.
type GRPCServer struct {
	checker *usecases.Checker

	way.UnimplementedStateServiceServer
}

// Update handles the Update RPC call.
//
// It receives a stream of UpdateRequest messages from the client and responds
// with an UpdateResponse message for each request. It continues to receive
// requests until the client closes the stream.
//
// Each UpdateRequest message contains a list of UUIDs that need to be updated.
// These UUIDs are used to uniquely identify the request and can be used to track
// the request throughout the system.
//
// For each UpdateRequest message, the server sends an empty UpdateResponse
// message to indicate that the update operation was successful.
//
// If there is a problem with receiving or sending messages, an error is returned.
func (s *GRPCServer) Update(stream way.StateService_UpdateServer) error {
	// Process requests from the client stream.
	for {
		// Receive the next request from the client.
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		// Get the list of UUIDs from the request.
		for _, id := range req.GetIds() {
			// Convert the UUID to a string.
			sid := uuidconv.DoubleInt2UUID(id.GetHigh(), id.GetLow())
			// Send the UUID to the checker.
			s.checker.Send(sid)
		}

		// Send an empty UpdateResponse message to the client.
		if err := stream.SendMsg(&way.UpdateResponse{}); err != nil {
			return err
		}
	}
}
