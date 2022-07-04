package logging

import (
	"github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-backend-service-common/acorns/repository"
)

// --- implementing Acorn ---

func New() auacornapi.Acorn {
	return &LoggingImpl{}
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
	r.Logger().NoCtx().Info().Print("logging is now available")

	return nil
}

func (r *LoggingImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	return nil
}
