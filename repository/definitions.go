package repository

import "github.com/thealamu/linkedinsignin/model"

type (
	UserCreator interface {
		CreateUser(user model.User) (*model.User, error)
	}

	UserRepositoryInterface interface {
		UserCreator
	}
)