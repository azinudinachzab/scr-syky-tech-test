package service

import (
	"context"
	"time"

	"github.com/azinudinachzab/scr-syky-tech-test/model"
	"github.com/azinudinachzab/scr-syky-tech-test/repository"
	"github.com/go-playground/validator/v10"
)

type Dependency struct {
	Validator *validator.Validate
	Repo      repository.Repository
	Conf      model.Configuration
}

type AppService struct {
	validator *validator.Validate
	repo      repository.Repository
	conf      model.Configuration
}

func NewService(dep Dependency) *AppService {
	return &AppService{
		validator: dep.Validator,
		repo:      dep.Repo,
		conf:      dep.Conf,
	}
}

type Service interface {
	Registration(ctx context.Context, req model.RegistrationRequest) error
	SendDailyBirthdayPromo(ctx context.Context, execDate time.Time)
}
