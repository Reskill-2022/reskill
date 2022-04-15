package email

import (
	"context"
	"github.com/thealamu/linkedinsignin/model"
)

type (
	Emailer interface {
		// Welcome sends a welcome email to the user
		Welcome(ctx context.Context, user *model.User) error
	}
)
