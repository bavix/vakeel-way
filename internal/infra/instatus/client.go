package instatus

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/bavix/vakeel-way/internal/domain/entities"
)

// Api is a client for the Instatus API.
//
// The Instatus API is used to send status updates to the Instatus service.
// The service provides a simple REST API that accepts a POST request with a
// JSON payload to update the status of a service.
//
// Fields:
//
//	client: The HTTP client used to make requests to the Instatus API.
type Api struct {
	// client: The HTTP client used to make requests to the Instatus API.
	client *http.Client
}

// NewApi creates a new Instatus API client.
//
// It takes a variable number of Option functions as parameters and returns a
// pointer to an Api struct. The Option functions are used to customize the
// Api struct before it is returned. If no options are provided, a new Api struct
// is created with default settings.
//
// Option: An option function that takes a pointer to an Api struct and modifies
// it. This allows for customization of the Api struct before it is returned.
// Multiple options can be passed to NewApi to configure the Api struct.
//
// Example:
//
//	api := instatus.NewApi(
//		instatus.WithClient(&http.Client{Timeout: 5 * time.Second}),
//	)
//
//	// Create a new Api struct with default settings for the http.Client.
//	api := &instatus.Api{
//		client: &http.Client{},
//	}
//
//	// Apply all provided options to the Api struct.
//	for _, op := range ops {
//		op(api)
//	}
//
//	return api
//
// Parameters:
//
//	ops: A variadic number of Option functions.
//
// Returns:
//
//	A pointer to an Api struct.
func NewApi(ops ...Option) *Api {
	// Create a new Api struct with default settings for the http.Client.
	api := &Api{
		client: &http.Client{},
	}

	// Apply all provided options to the Api struct.
	for _, op := range ops {
		op(api)
	}

	return api
}

// Option is a function that can be used to configure an Api instance.
//
// It takes a pointer to an Api instance and returns nothing.
type Option func(*Api)

// WithClient returns an Option function that sets the http.Client used to send HTTP requests
// to the Instatus API.
//
// The provided http.Client is used to send HTTP requests to the Instatus API. If the
// http.Client is nil, a new http.Client with default settings is created.
//
// Example:
//
//	api := instatus.NewApi(
//		instatus.WithClient(&http.Client{Timeout: 5 * time.Second}),
//	)
//
// Returns an Option function that sets the http.Client used to send HTTP requests
// to the Instatus API.
func WithClient(c http.Client) Option {
	// The Option function returned by this function sets the http.Client used to send HTTP
	// requests to the Instatus API. If the http.Client is nil, a new http.Client with default
	// settings is created.
	//
	// Parameters:
	// - api: A pointer to the Api struct that will be customized.
	//
	// Returns:
	// None.
	return func(api *Api) {
		// Set the http.Client in the Api struct to the provided *http.Client,
		// or create a new http.Client with default settings if the *http.Client
		// is nil.
		api.client = &c
	}
}

// Send sends a POST request to the given URL with the specified status.
//
// The request is sent with the provided context and the status is used to
// determine the value of the "trigger" field in the request payload.
// The request payload is a JSON object with a single key "trigger" and a
// value that corresponds to the status. The context is used to cancel the
// request if it takes too long to complete.
//
// Returns an error if the request cannot be created, sent, or if the response
// cannot be read.
//
// Parameters:
// - ctx: The context.Context to use for the request.
// - url: The URL to send the request to.
// - status: The entities.Status to use in the request payload.
func (s *Api) Send(ctx context.Context, url string, status entities.Status) error {
	// Create the request payload as a JSON object with a single key "trigger"
	// and a value that corresponds to the status.
	// The payload is created as a string with the JSON object in it.
	// The string is created using fmt.Sprintf() with the status as the
	// parameter.
	payload := fmt.Sprintf(`{"trigger": "%s"}`, status)

	// Create a new HTTP request with the provided context and the specified URL.
	// The request is a POST request with the payload as the request body.
	// The request is created using http.NewRequestWithContext().
	req, err := http.NewRequestWithContext(ctx, "POST", url,
		bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}

	// Set the "Content-Type" header of the request to "application/json" to
	// indicate that the request body is in JSON format.
	// The header is set using the Set() method of the Header map.
	req.Header.Set("Content-Type", "application/json")

	// Send the request and get the response.
	// The request is sent using the Do() method of the client.
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// No error occurred, return nil.
	// The response is not needed, so it is closed immediately.
	return nil
}
