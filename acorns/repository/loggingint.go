package repository

import auloggingapi "github.com/StephanHCB/go-autumn-logging/api"

const LoggingAcornName = "logging"

// Logging is the central singleton representing the logging subsystem.
type Logging interface {
	IsLogging() bool

	// Setup uses the (at this point partially unvalidated) configuration to configure either Json/ECS or Plain logging
	Setup()

	// Logger gives you access to the logging implementation
	Logger() auloggingapi.LoggingImplementation
}
