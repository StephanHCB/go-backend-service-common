package config

// stuff you normally won't need to look at unless you're working on the configuration itself

import (
	"context"
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
	auconfigenv "github.com/StephanHCB/go-autumn-config-env"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	auzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	"github.com/StephanHCB/go-backend-service-common/acorns/repository"
)

type ConfigImpl struct {
	Logging           repository.Logging
	validationContext context.Context

	// place to store the parsed and validated config values for quick no-parse access
	VApplicationName   string
	VServerAddress     string
	VEnvironment       string
	VPlatform          string
	VLogstyle          string
	VVaultServer       string
	VVaultCertFile     string
	VVaultSecretPath   string
	VLocalVaultToken   string
	VVaultK8sRole      string
	VVaultK8sTokenPath string
	VVaultK8sBackend   string

	VServerPortValue  uint16
	VMetricsPortValue uint16

	CustomConfiguration repository.CustomConfiguration
}

// New initially creates the instance with no logging (circular dependency)
func New(customConfig repository.CustomConfiguration, additionalConfigItems []auconfigapi.ConfigItem) auacornapi.Acorn {
	instance := &ConfigImpl{
		CustomConfiguration: customConfig,
	}

	allConfigItems := PredefinedConfigItems
	for _, item := range additionalConfigItems {
		allConfigItems = append(allConfigItems, item)
	}

	warnFunc := func(message string) {
		if instance.Logging != nil && instance.validationContext != nil {
			instance.Logging.Logger().Ctx(instance.validationContext).Error().Print(message)
		}
	}

	err := auconfigenv.Setup(allConfigItems, warnFunc)
	if err != nil {
		// we do not have logging yet, and cannot read configuration, so this is going to be incomplete by necessity
		auzerolog.SetupJsonLogging(ApplicationName)
		aulogging.Logger.NoCtx().Fatal().WithErr(err).Print("failed to read configuration defaults from code - only strings are supported! BAILING OUT")
	}

	return instance
}

func (r *ConfigImpl) Read() error {
	return auconfigenv.Read()
}

func (r *ConfigImpl) Validate(ctx context.Context) error {
	r.validationContext = ctx
	return auconfigenv.Validate()
}

func (r *ConfigImpl) ObtainPredefinedValues() {
	r.VApplicationName = auconfigenv.Get(KeyApplicationName)
	r.VServerAddress = auconfigenv.Get(KeyServerAddress)
	r.VEnvironment = auconfigenv.Get(KeyEnvironment)
	r.VPlatform = auconfigenv.Get(KeyPlatform)
	r.VLogstyle = auconfigenv.Get(KeyLogstyle)
	r.VVaultServer = auconfigenv.Get(KeyVaultServer)
	r.VVaultCertFile = auconfigenv.Get(KeyVaultCertificateFile)
	r.VVaultSecretPath = auconfigenv.Get(KeyVaultSecretPath)
	r.VLocalVaultToken = auconfigenv.Get(KeyLocalVaultToken)
	r.VVaultK8sRole = auconfigenv.Get(KeyVaultKubernetesRole)
	r.VVaultK8sTokenPath = auconfigenv.Get(KeyVaultKubernetesTokenPath)
	r.VVaultK8sBackend = auconfigenv.Get(KeyVaultKubernetesBackend)

	// after validate, these cannot fail any more
	vServerPortValue, _ := auconfigenv.AToUint(auconfigenv.Get(KeyServerPort))
	r.VServerPortValue = uint16(vServerPortValue)

	vMetricsPortValue, _ := auconfigenv.AToUint(auconfigenv.Get(KeyMetricsPort))
	r.VMetricsPortValue = uint16(vMetricsPortValue)
}
