package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/thealamu/linkedinsignin/linkedin"
	"github.com/thealamu/linkedinsignin/model"
	"github.com/thealamu/linkedinsignin/repository"
	"github.com/thealamu/linkedinsignin/requests"
	"net/http"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (u *UserController) CreateUser(userCreator repository.UserCreator, service linkedin.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var requestBody requests.CreateUserRequest

		err := c.Bind(&requestBody)
		if err != nil {
			return HandleError(c, err, http.StatusBadRequest)
		}

		profile, err := service.GetProfile(requestBody.Email)
		if err != nil {
			return HandleError(c, err, http.StatusBadRequest)
		}

		// validate profile has required fields
		if profile.Name == "" || profile.Location == "" || profile.Phone == "" {
			return HandleError(c, fmt.Errorf("some fields are missing on your LinkedIn"), http.StatusBadRequest)
		}

		u := model.User{
			Email:       requestBody.Email,
			Name:        profile.Name,
			LinkedInURL: profile.ProfileURL,
			Location:    profile.Location,
			Phone:       profile.Phone,
		}

		user, err := userCreator.CreateUser(ctx, u)
		if err != nil {
			return HandleError(c, err, http.StatusInternalServerError)
		}

		return HandleSuccess(c, user, http.StatusCreated)
	}
}

func (u *UserController) UpdateUser(userUpdater repository.UserUpdater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var requestBody requests.UpdateUserRequest

		err := c.Bind(&requestBody)
		if err != nil {
			return HandleError(c, err, http.StatusBadRequest)
		}

		update := model.User{
			Email: c.Param("email"),
		}

		{
			if requestBody.Representation != "" {
				update.Representation = requestBody.Representation
			}

			if requestBody.Gender != "" {
				update.Gender = requestBody.Gender
			}

			if requestBody.AgeGroup != "" {
				update.AgeGroup = requestBody.AgeGroup
			}

			if requestBody.EmploymentStatus != "" {
				update.EmploymentStatus = requestBody.EmploymentStatus
			}

			if requestBody.HighestSchool != "" {
				update.HighestSchool = requestBody.HighestSchool
			}

			if requestBody.CanWorkInUSA != "" {
				update.CanWorkInUSA = requestBody.CanWorkInUSA
			}

			if requestBody.LearningTrack != "" {
				update.LearningTrack = requestBody.LearningTrack
			}

			if requestBody.TechExperience != "" {
				update.TechExperience = requestBody.TechExperience
			}

			if requestBody.HoursPerWeek != "" {
				update.HoursPerWeek = requestBody.HoursPerWeek
			}

			if requestBody.Referral != "" {
				update.Referral = requestBody.Referral
			}
		}

		user, err := userUpdater.UpdateUser(ctx, update)
		if err != nil {
			return HandleError(c, err, http.StatusInternalServerError)
		}

		return HandleSuccess(c, user, http.StatusOK)
	}
}

func (u *UserController) GetUser(userGetter repository.UserGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		user, err := userGetter.GetUser(ctx, c.Param("email"))
		if err != nil {
			return HandleError(c, err, http.StatusInternalServerError)
		}

		return HandleSuccess(c, user, http.StatusOK)
	}
}
