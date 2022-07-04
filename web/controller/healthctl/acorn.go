package healthctl

import (
	"github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-backend-service-common/acorns/controller"
)

// --- implementing Acorn ---

func New() auacornapi.Acorn {
	return &HealthCtlImpl{}
}

func (a *HealthCtlImpl) IsHealthController() bool {
	return true
}

func (a *HealthCtlImpl) AcornName() string {
	return controller.HealthControllerAcornName
}

func (a *HealthCtlImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	return nil
}

func (a *HealthCtlImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	return nil
}

func (a *HealthCtlImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	return nil
}
