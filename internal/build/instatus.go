package build

import "github.com/bavix/vakeel-way/internal/infra/instatus"

// inStatusClient returns a new instance of the instatus.Api struct.
//
// The instatus.Api struct is used to interact with the Instatus API.
// It provides methods for sending status updates to the Instatus service.
//
// The function returns a pointer to an instatus.Api struct.
func (b *Builder) inStatusClient() *instatus.Api {
	// Create a new instance of the instatus.Api struct.
	// The struct is created with default settings.
	return instatus.NewApi()
}
