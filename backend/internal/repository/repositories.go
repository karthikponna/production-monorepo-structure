package repository

import "github.com/karthikponna/production-monorepo-structure/internal/server"

type Repositories struct{}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}