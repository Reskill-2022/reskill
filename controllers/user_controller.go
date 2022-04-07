package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (u *UserController) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		return fmt.Errorf("not implemented")
	}
}

func (u *UserController) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		return fmt.Errorf("not implemented")
	}
}

func (u *UserController) GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		return fmt.Errorf("not implemented")
	}
}
