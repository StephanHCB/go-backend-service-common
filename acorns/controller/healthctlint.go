package controller

import (
	"context"
	"github.com/go-chi/chi/v5"
)

const HealthControllerAcornName = "healthctl"

// HealthController provides a basic health endpoint
type HealthController interface {
	IsHealthController() bool

	WireUp(ctx context.Context, router chi.Router)
}
