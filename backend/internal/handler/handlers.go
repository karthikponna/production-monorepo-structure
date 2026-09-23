package handler

import (
	"github.com/karthikponna/production-monorepo-structure/internal/server"
	"github.com/karthikponna/production-monorepo-structure/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
	}
}