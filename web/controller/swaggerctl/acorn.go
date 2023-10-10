package swaggerctl

import (
	"github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-backend-service-common/acorns/controller"
)

// --- implementing Acorn ---

func New() auacornapi.Acorn {
	return &SwaggerCtlImpl{}
}

// NewNoAcorn performs the full Acorn lifecycle for this component, no further setup necessary
func NewNoAcorn() controller.SwaggerController {
	return &SwaggerCtlImpl{}
}

func (a *SwaggerCtlImpl) IsSwaggerController() bool {
	return true
}

func (a *SwaggerCtlImpl) AcornName() string {
	return controller.SwaggerControllerAcornName
}

func (a *SwaggerCtlImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	return nil
}

func (a *SwaggerCtlImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	return nil
}

func (a *SwaggerCtlImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	return nil
}
