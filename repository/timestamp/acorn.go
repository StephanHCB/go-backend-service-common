package timestamp

import (
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-backend-service-common/acorns/repository"
	"time"
)

func New() auacornapi.Acorn {
	return &TimestampImpl{}
}

// NewNoAcorn replaces everything in the full Acorn lifecycle for this component, no need to call anything else.
func NewNoAcorn(nowFunc func() time.Time) repository.Timestamp {
	return &TimestampImpl{
		Timestamp: nowFunc,
	}
}

func (r *TimestampImpl) IsTimestamp() bool {
	return true
}

func (r *TimestampImpl) AcornName() string {
	return repository.TimestampAcornName
}

func (r *TimestampImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	r.Timestamp = time.Now

	return nil
}

func (r *TimestampImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	return nil
}

func (r *TimestampImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	return nil
}
