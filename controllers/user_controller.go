package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/thealamu/linkedinsignin/model"
	"github.com/thealamu/linkedinsignin/repository"
	"github.com/thealamu/linkedinsignin/requests"
	"net/http"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (u *UserController) CreateUser(userCreator repository.UserCreator) echo.HandlerFunc {
	return func(c echo.Context) error {
		var requestBody requests.CreateUserRequest

		err := c.Bind(&requestBody)
		if err != nil {
			return HandleError(c, err, http.StatusBadRequest)
		}

		//todo: pull and validate user info from linkedin

		u := model.User{
			Email: requestBody.Email,
		}

		user, err := userCreator.CreateUser(u)
		if err != nil {
			return HandleError(c, err, http.StatusInternalServerError)
		}

		return HandleSuccess(c, user, http.StatusCreated)
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
