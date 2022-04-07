package repository

import (
	"context"
	"fmt"
	"github.com/thealamu/linkedinsignin/model"
)

type UserRepository struct{}

var _ UserRepositoryInterface = (*UserRepository)(nil)

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u *UserRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (u *UserRepository) UpdateUser(ctx context.Context, user model.User) (*model.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (u *UserRepository) GetUser(ctx context.Context, email string) (*model.User, error) {
	return nil, fmt.Errorf("not implemented")
}
