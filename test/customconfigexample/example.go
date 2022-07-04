package customconfigexample

import (
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
	"github.com/StephanHCB/go-backend-service-common/repository/config"
)

type CustomConfigurationWithOneField interface {
	Obtain(func(key string) string)

	MyCustomField() string
}

const (
	KeyMyCustomField = "MY_CUSTOM_FIELD"
)

var CustomConfigItems = []auconfigapi.ConfigItem{
	{
		Key:         KeyMyCustomField,
		EnvName:     KeyMyCustomField,
		Default:     "",
		Description: "an example custom config field",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
}

func New() auacornapi.Acorn {
	customConfigInstance := &CustomConfigurationWithOneFieldImpl{}
	return config.New(customConfigInstance, CustomConfigItems)
}

// implementing the interface here

type CustomConfigurationWithOneFieldImpl struct {
	VMyCustomField string
}

func (c *CustomConfigurationWithOneFieldImpl) Obtain(getter func(key string) string) {
	c.VMyCustomField = getter(KeyMyCustomField)
}

func (c *CustomConfigurationWithOneFieldImpl) MyCustomField() string {
	return c.VMyCustomField
}
