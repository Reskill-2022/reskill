package email

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/model"
)

type (
	Emailer interface {
		// Welcome sends a welcome email to the user
		Welcome(ctx context.Context, user *model.User) error
	}

	ses struct {
		logger zerolog.Logger
	}
)

func New(logger zerolog.Logger) *ses {
	return &ses{
		logger: logger,
	}
}
