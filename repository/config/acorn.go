package config

import (
	"context"
	"github.com/StephanHCB/go-autumn-acorn-registry/api"
	auconfigenv "github.com/StephanHCB/go-autumn-config-env"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	auzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	"github.com/StephanHCB/go-backend-service-common/acorns/repository"
)

// --- implementing Acorn ---

// New() see plumbing.go

func (r *ConfigImpl) IsConfiguration() bool {
	return true
}

func (r *ConfigImpl) AcornName() string {
	return repository.ConfigurationAcornName
}

func (r *ConfigImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	r.Logging = registry.GetAcornByName(repository.LoggingAcornName).(repository.Logging)

	// load the configuration, so it is available for logging setup (but don't validate it yet)
	if err := r.Read(); err != nil {
		// we do not have logging yet, and cannot read configuration, so this is going to be incomplete by necessity
		auzerolog.SetupJsonLogging(ApplicationName)
		aulogging.Logger.NoCtx().Error().WithErr(err).Print("failed to obtain configuration. BAILING OUT")
		return err
	}

	// make unvalidated configuration values available for logging setup
	r.ObtainValuesNeededForLogging()

	return nil
}

func (r *ConfigImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	err := registry.SetupAfter(r.Logging.(auacornapi.Acorn))
	if err != nil {
		return err
	}

	ctx := auzerolog.AddLoggerToCtx(context.Background())

	if err := r.Validate(ctx); err != nil {
		r.Logging.Logger().Ctx(ctx).Error().WithErr(err).Print("failed to validate configuration. BAILING OUT. See log messages above for individual errors")
		return err
	}

	r.ObtainPredefinedValues()
	r.CustomConfiguration.Obtain(auconfigenv.Get)

	r.Logging.Logger().Ctx(ctx).Info().Print("successfully set up configuration and logging")

	return nil
}

func (r *ConfigImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	return nil
}
