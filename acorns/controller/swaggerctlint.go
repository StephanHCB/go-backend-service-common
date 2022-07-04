package controller

import (
	"context"
	"github.com/go-chi/chi/v5"
)

const SwaggerControllerAcornName = "swaggerctl"

// SwaggerController provides the swagger ui and serves our openapi v3 spec from docs/openapi-v3-spec.json
type SwaggerController interface {
	IsSwaggerController() bool

	WireUp(ctx context.Context, router chi.Router)
}
