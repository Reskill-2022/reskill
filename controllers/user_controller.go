package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/errors"
	"github.com/thealamu/linkedinsignin/linkedin"
	"github.com/thealamu/linkedinsignin/model"
	"github.com/thealamu/linkedinsignin/repository"
	"github.com/thealamu/linkedinsignin/requests"
	"net/http"
)

type UserController struct {
	logger zerolog.Logger
}

func NewUserController(logger zerolog.Logger) *UserController {
	return &UserController{logger}
}

func (u *UserController) CreateUser(userCreator repository.UserCreator, service linkedin.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var requestBody requests.CreateUserRequest

		err := c.Bind(&requestBody)
		if err != nil {
			return u.HandleError(c, errors.From(err, "Invalid JSON Request Body", 400), http.StatusBadRequest)
		}

		profile, err := service.GetProfile(requestBody.Email)
		if err != nil {
			return u.HandleError(c, errors.From(err, "Failed to Get user profile from LinkedIn", errors.CodeFrom(err)), http.StatusBadRequest)
		}

		data := model.User{
			Email:       requestBody.Email,
			Name:        profile.Name,
			LinkedInURL: profile.ProfileURL,
			Location:    profile.Location,
			Phone:       profile.Phone,
		}

		user, err := userCreator.CreateUser(ctx, data)
		if err != nil {
			return u.HandleError(c, err, errors.CodeFrom(err))
		}

		return HandleSuccess(c, user, http.StatusCreated)
	}
}

func (u *UserController) UpdateUser(userGetter repository.UserGetter, userUpdater repository.UserUpdater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var requestBody requests.UpdateUserRequest

		err := c.Bind(&requestBody)
		if err != nil {
			return u.HandleError(c, err, http.StatusBadRequest)
		}

		update, err := userGetter.GetUser(ctx, c.Param("email"))
		if err != nil {
			return u.HandleError(c, err, errors.CodeFrom(err))
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

		user, err := userUpdater.UpdateUser(ctx, *update)
		if err != nil {
			return u.HandleError(c, err, errors.CodeFrom(err))
		}

		return HandleSuccess(c, user, http.StatusOK)
	}
}

func (u *UserController) GetUser(userGetter repository.UserGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		user, err := userGetter.GetUser(ctx, c.Param("email"))
		if err != nil {
			return u.HandleError(c, err, errors.CodeFrom(err))
		}

		return HandleSuccess(c, user, http.StatusOK)
	}
}
