package repository

import (
	"context"
	"fmt"
	"github.com/thealamu/linkedinsignin/errors"
	"github.com/thealamu/linkedinsignin/model"
)

type UserRepository struct {
	data map[string]*model.User
}

var _ UserRepositoryInterface = (*UserRepository)(nil)

func NewUserRepository() *UserRepository {
	return &UserRepository{
		data: make(map[string]*model.User),
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	if _, ok := u.data[user.Email]; ok {
		return nil, errors.New(fmt.Sprintf("User with email %s already exists", user.Email), 400)
	}
	u.data[user.Email] = &user
	return &user, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, user model.User) (*model.User, error) {
	if _, ok := u.data[user.Email]; !ok {
		return nil, errors.New(fmt.Sprintf("User with email %s does not exist", user.Email), 400)
	}
	u.data[user.Email] = &user
	return &user, nil
}

func (u *UserRepository) GetUser(ctx context.Context, email string) (*model.User, error) {
	if _, ok := u.data[email]; !ok {
		return nil, errors.New(fmt.Sprintf("User with email %s does not exist", email), 400)
	}
	return u.data[email], nil
}
