package service

import (
	"github.com/karthikponna/production-monorepo-structure/internal/lib/job"
	"github.com/karthikponna/production-monorepo-structure/internal/repository"
	"github.com/karthikponna/production-monorepo-structure/internal/server"
)

type Services struct {
	Auth *AuthService
	Job  *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	return &Services{
		Job:  s.Job,
		Auth: authService,
	}, nil
}