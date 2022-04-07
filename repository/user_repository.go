package repository

import (
	"fmt"
	"github.com/thealamu/linkedinsignin/model"
)

type UserRepository struct{}

var _ UserRepositoryInterface = (*UserRepository)(nil)

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u *UserRepository) CreateUser(user model.User) (*model.User, error) {
	return nil, fmt.Errorf("not implemented")
}
