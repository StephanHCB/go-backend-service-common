package consumer

import (
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
	auconfigenv "github.com/StephanHCB/go-autumn-config-env"
	"github.com/StephanHCB/go-backend-service-common/repository/config"
	"github.com/StephanHCB/go-backend-service-common/repository/logging"
	"github.com/StephanHCB/go-backend-service-common/repository/vault"
)

const contractTestConfigurationPath = "../../resources/contract-test-config.yaml"

const keyMyConfigKey = "MY_CONFIG_KEY"

var (
	configImpl  *config.ConfigImpl
	loggingImpl *logging.LoggingImpl
	vaultImpl   *vault.Impl

	customConfigItems = []auconfigapi.ConfigItem{
		{
			Key:         keyMyConfigKey,
			EnvName:     keyMyConfigKey,
			Default:     "",
			Description: "demo config item that is to be filled from a secret",
			Validate:    auconfigenv.ObtainNotEmptyValidator(),
		},
	}
	customConfigImpl *customConfigMock
)

// simplest possible custom config implementation with the one key we need

type customConfigMock struct {
	VMyConfigKey string
}

func (c *customConfigMock) Obtain(getter func(key string) string) {
	c.VMyConfigKey = getter(keyMyConfigKey)
}

// global setup for contract tests

func tstSetup() error {
	// setup custom config to write the secret value to
	customConfigImpl = &customConfigMock{}

	// setup test configuration
	allExtraConfigItems := append(vault.ConfigItems, customConfigItems...)
	configImpl = config.New(customConfigImpl, allExtraConfigItems).(*config.ConfigImpl)
	auconfigenv.LocalConfigFileName = contractTestConfigurationPath
	err := configImpl.Read()
	if err != nil {
		return err
	}
	// intentionally not validating the configuration
	configImpl.ObtainPredefinedValues()
	// not obtaining custom config here, see during actual test

	// setup logging
	loggingImpl = logging.New().(*logging.LoggingImpl)
	loggingImpl.Configuration = configImpl
	loggingImpl.Setup()
	configImpl.Logging = loggingImpl

	// setup vault
	vaultImpl = vault.New().(*vault.Impl)
	vaultImpl.Configuration = configImpl
	vaultImpl.Logging = loggingImpl

	return nil
}
