package repository

import (
	"context"
	"github.com/thealamu/linkedinsignin/model"
)

type (
	UserCreator interface {
		CreateUser(ctx context.Context, user model.User) (*model.User, error)
	}

	UserUpdater interface {
		UpdateUser(ctx context.Context, user model.User) (*model.User, error)
	}

	UserGetter interface {
		GetUser(ctx context.Context, email string) (*model.User, error)
	}

	UserRepositoryInterface interface {
		UserCreator
		UserUpdater
		UserGetter
	}
)
