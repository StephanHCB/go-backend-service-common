package repository

import "context"

const ConfigurationAcornName = "configuration"

// CustomConfiguration is just used as a type constraint for the part that YOU need to provide.
//
// You should declare an interface of your own that extends it (i.e. has its method),
// but also implements type safe accessors for all your custom configuration variables.
type CustomConfiguration interface {
	// Obtain gets called when all configuration strings are parsed and validated.
	//
	// It is given an accessor function that retrieves the parsed configuration values by their key.
	Obtain(func(key string) string)

	// ... in your implementation, you should put accessors here
}

// Configuration is the central singleton representing the configuration.
//
// In normal operation, all values come from environment variables, but for localhost convenience we
// also support reading a yaml file.
type Configuration interface {
	IsConfiguration() bool

	// Read the configuration (does not log, so needs no logging yet, and not context aware)
	//
	// In order of decreasing precedence:
	// - environment variable
	// - local flat yaml file (intended for localhost only)
	// - default value
	Read() error

	// Validate the configuration (logs detailed validation errors, so needs logging set up)
	Validate(ctx context.Context) error

	// Custom gives you access to your custom configuration value object.
	//
	// TODO sorry for forcing you to type cast for now, but interfaces with generics aren't quite there yet.
	// It's funny, but if you make the return type of this method a generic type argument, there is actually
	// no way to implement the method (except for one specific type, which is precisely NOT what we need).
	//
	// I have found it convenient to do the type cast in the acorn setup, meaning, if your code needs any
	// custom configuration values, just obtain this reference at setup time rather than casting every
	// time you need a configuration value.
	Custom() CustomConfiguration

	// accessors for common configuration properties which are always provided for you

	ApplicationName() string

	ServerAddress() string
	ServerPort() uint16
	MetricsPort() uint16

	Environment() string
	Platform() string

	PlainLogging() bool

	VaultServer() string
	VaultCertificateFile() string
	VaultSecretPath() string

	LocalVault() bool
	LocalVaultToken() string

	VaultKubernetesRole() string
	VaultKubernetesTokenPath() string
	VaultKubernetesBackend() string

	CorsAllowOrigin() string
}
