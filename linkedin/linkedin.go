package linkedin

import (
	"fmt"
	"github.com/rs/zerolog"
)

type (
	Service interface {
		GetProfile(GetProfileInput) (GetProfileOutput, error)
	}

	GetProfileInput struct {
		Email string
	}

	GetProfileOutput struct {
		FullName       string
		ProfilePhoto   string
		Location       string
		Email          string
		Phone          string
		WorkExperience interface{}
	}

	lkd struct {
		logger zerolog.Logger
		apiKey string
	}
)

func New(logger zerolog.Logger, apiKey string) Service {
	return &lkd{
		logger: logger,
		apiKey: apiKey,
	}
}

func (l *lkd) GetProfile(input GetProfileInput) (GetProfileOutput, error) {
	return GetProfileOutput{}, fmt.Errorf("not implemented")
}
