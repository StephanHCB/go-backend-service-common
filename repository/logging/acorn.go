package logging

import (
	"github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-backend-service-common/acorns/repository"
)

// --- implementing Acorn ---

func New() auacornapi.Acorn {
	return &LoggingImpl{}
}

// NewNoAcorn constructs and assembles the component, but it cannot set it up yet.
//
// You still need to call Setup() after the configuration was loaded (but not validated)
func NewNoAcorn(configuration repository.Configuration) repository.Logging {
	return &LoggingImpl{
		Configuration: configuration,
	}
}

func (r *LoggingImpl) IsLogging() bool {
	return true
}

func (r *LoggingImpl) AcornName() string {
	return repository.LoggingAcornName
}

func (r *LoggingImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	r.Configuration = registry.GetAcornByName(repository.ConfigurationAcornName).(repository.Configuration)

	return nil
}

func (r *LoggingImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	// configuration was loaded (but not validated) during Assemble

	r.Setup()

	return nil
}

func (r *LoggingImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	return nil
}
