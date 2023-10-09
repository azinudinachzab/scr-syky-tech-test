package service

import (
	"context"

	"github.com/azinudinachzab/scr-syky-tech-test/model"
	"github.com/azinudinachzab/scr-syky-tech-test/pkg/clock"
	"github.com/azinudinachzab/scr-syky-tech-test/pkg/password"
	"github.com/azinudinachzab/scr-syky-tech-test/repository"
	"github.com/go-playground/validator/v10"
)

type Dependency struct {
	Validator *validator.Validate
	Repo      repository.Repository
	Clock     clock.Time
	Hash      password.Hasher
	Conf      model.Configuration
}

type AppService struct {
	validator *validator.Validate
	repo      repository.Repository
	clock     clock.Time
	hash      password.Hasher
	conf      model.Configuration
}

func NewService(dep Dependency) *AppService {
	return &AppService{
		validator: dep.Validator,
		repo:      dep.Repo,
		clock:     dep.Clock,
		hash:      dep.Hash,
		conf:      dep.Conf,
	}
}

type Service interface {
	Registration(ctx context.Context, req model.RegistrationRequest) error
}
