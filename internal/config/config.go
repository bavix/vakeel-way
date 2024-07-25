package config

import (
	"net"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

// Webhooks is a slice of WebhookConfig.
type Webhooks []WebhookConfig

// AsMap converts the slice of WebhookConfig into a map.
//
// The function takes the slice of WebhookConfig as input and returns a map
// with the ID of the WebhookConfig as the key and the target URL as the value.
// The map is created with preallocated capacity to avoid resizing during iteration.
//
// Returns:
// - A map[uuid.UUID]string containing the converted data.
func (w Webhooks) AsMap() map[uuid.UUID]string {
	// Create a map with preallocated capacity for the length of the slice.
	// This is done to avoid resizing the map during the iteration.
	m := make(map[uuid.UUID]string, len(w))

	// Iterate over each WebhookConfig in the slice.
	// The range keyword is used to iterate over the slice and get the index and value.
	for i := range w {
		// Use the ID of the WebhookConfig as the key in the map,
		// and the target of the WebhookConfig as the value.
		m[w[i].ID] = w[i].Target
	}

	// Return the WebhooksMap containing the converted data.
	return m
}

// Config represents the configuration of the application.
//
// It contains the configuration for the logger and the gRPC server.
type Config struct {
	// Log is the configuration for the logger.
	//
	// The logger configuration contains the log level and the log format.
	Log LogConfig `yaml:"log"`

	// GRPC is the configuration of the gRPC server.
	//
	// The gRPC server configuration contains the network, address, and maximum message size.
	GRPC GRPCConfig `yaml:"grpc"`

	// Webhooks is the configuration for the webhooks.
	//
	// The webhook configuration contains the unique identifier and the target URL of the webhook.
	Webhooks Webhooks `yaml:"webhooks"`
}

// WebhookConfig represents the configuration for the webhook.
//
// It contains the unique identifier and the target URL of the webhook.
type WebhookConfig struct {
	// ID is the unique identifier for the webhook.
	//
	// The ID is a string that uniquely identifies the webhook.
	// It is used to distinguish between different webhooks.
	ID uuid.UUID `yaml:"id"`

	// Target is the target URL of the webhook.
	//
	// The target URL is the URL that will be notified when an event is triggered.
	// It should be a valid URL that the webhook can reach.
	//
	// Example: "https://example.com/webhook"
	Target string `yaml:"target"`
}

// LogConfig represents the configuration for the logger.
//
// It contains the log level.
type LogConfig struct {
	// Level is the log level.
	//
	// The log level determines the severity of the log messages that will be logged.
	// The possible values are:
	// - "debug" for low-level debugging information
	// - "info" for informational messages
	// - "warn" for warnings or potential issues
	// - "error" for errors that should be addressed
	// - "fatal" for critical errors that cause the application to exit
	Level string `yaml:"level"`
}

// GRPCConfig represents the configuration of the gRPC server.
//
// It contains the network protocol, host address, and port number to use for the gRPC server.
type GRPCConfig struct {
	// Network is the network protocol to use for the gRPC server.
	// It is the transport protocol that the server will use to communicate with clients.
	// The possible values are:
	// - "tcp" for IPv4 or IPv6
	// - "udp" for connectionless communication
	Network string `yaml:"network"`

	// Host is the host address to use for the gRPC server.
	// It can be an IP address or a hostname.
	Host string `yaml:"host"`

	// Port is the port number to use for the gRPC server.
	// It is the port number where the gRPC server will listen for incoming connections.
	Port string `yaml:"port"`
}

// Addr returns the address of the gRPC server as a string.
//
// The address is formed by joining the Host and Port fields of the GRPCConfig
// using the net.JoinHostPort function. The resulting string has the format
// "host:port".
//
// Parameters:
// - None
//
// Returns:
// - string: The address of the gRPC server in the format "host:port".
func (c GRPCConfig) Addr() string {
	// Join the Host and Port fields of the GRPCConfig using the net.JoinHostPort
	// function. The resulting string has the format "host:port".
	return net.JoinHostPort(c.Host, c.Port)
}

// New reads the configuration from a YAML file and returns an instance of Config.
// It takes the path to the YAML file as a parameter and returns the parsed configuration
// or an error if there was an issue reading or parsing the file.
//
// The path parameter is a string that represents the path to the YAML file.
// It returns a Config instance and an error.
func New(path string) (Config, error) {
	// Create a new Config instance with default values
	// The default values are:
	// - log level: info
	// - network: tcp
	// - host: 0.0.0.0
	// - port: 4643
	cfg := Config{
		Log: LogConfig{
			Level: "info",
		},
		GRPC: GRPCConfig{
			Network: "tcp",
			Host:    "0.0.0.0",
			Port:    "4643",
		},
	}

	// Check if the YAML file exists
	// If the file does not exist, return the default config
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}

	// Read the contents of the YAML file
	// If there is an issue reading the file, return the error
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	// Decode the YAML contents into the Config instance
	// The Unmarshal function decodes the YAML data into the specified value.
	// It takes the YAML data as a byte slice and a pointer to the value to decode into.
	// In this case, we are decoding the YAML data into the Config instance.
	// If there is an issue decoding the YAML, return the error
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	// Return the Config instance and nil (indicating success)
	return cfg, nil
}
