package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/email"
	"github.com/thealamu/linkedinsignin/errors"
	"github.com/thealamu/linkedinsignin/linkedin"
	"github.com/thealamu/linkedinsignin/model"
	"github.com/thealamu/linkedinsignin/repository"
	"github.com/thealamu/linkedinsignin/requests"
	"io"
	"net/http"
	"strings"
	"time"
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
			return u.HandleError(c, errors.New("Invalid JSON Request Body", 400), http.StatusBadRequest)
		}

		userEmail := strings.ToLower(requestBody.Email)

		profile, err := service.GetProfile(userEmail)
		if err != nil {
			return u.HandleError(c, err, errors.CodeFrom(err))
		}

		// do validations
		if profile.Name == "" {
			return u.HandleError(c, errors.New("Invalid Profile. Found No Name", 400), http.StatusBadRequest)
		}
		if profile.Location == "" {
			return u.HandleError(c, errors.New("Invalid Profile. Please Set Your City and State of Residence on LinkedIn", 400), http.StatusBadRequest)
		}
		if profile.Photo == "" {
			return u.HandleError(c, errors.New("Invalid Profile. Please Set Your Profile Picture on LinkedIn", 400), http.StatusBadRequest)
		}
		if !profile.HasExperience {
			return u.HandleError(c, errors.New("Invalid Profile. Please Add Your Work Experience on LinkedIn", 400), http.StatusBadRequest)
		}

		//country := profile.Location
		//i := strings.LastIndex(profile.Location, ",")
		//if i > 0 {
		//	country = profile.Location[i+2:]
		//}
		//if country != "United States" {
		//	return u.HandleError(c, errors.New("Invalid Profile. For United States Only", 400), http.StatusBadRequest)
		//}

		data := model.User{
			Email:       userEmail,
			Name:        profile.Name,
			LinkedInURL: profile.ProfileURL,
			Location:    profile.Location,
			Phone:       profile.Phone,
			CreatedAt:   time.Now().UTC().String(),
		}

		user, err := userCreator.CreateUser(ctx, data)
		if err != nil {
			return u.HandleError(c, err, errors.CodeFrom(err))
		}

		return HandleSuccess(c, user, http.StatusCreated)
	}
}

func (u *UserController) UpdateUser(userGetter repository.UserGetter, userUpdater repository.UserUpdater, emailer email.Emailer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// dump request headers
		u.logger.Debug().Msgf("Request Headers: %+v", c.Request().Header)

		var requestBody requests.UpdateUserRequest

		// dump request body
		var body bytes.Buffer
		_, err := io.Copy(&body, c.Request().Body)
		if err != nil {
			u.logger.Debug().Msgf("Error reading request body: %s", err)
		}

		cp := body.Bytes()
		u.logger.Debug().Msgf("Request Body: %s", cp)

		err = json.NewDecoder(bytes.NewReader(cp)).Decode(&requestBody)
		if err != nil {
			return u.HandleError(c, err, http.StatusBadRequest)
		}

		update, err := userGetter.GetUser(ctx, c.Param("email"))
		if err != nil {
			return u.HandleError(c, err, errors.CodeFrom(err))
		}
		if update.Enrolled {
			return u.HandleError(c, errors.New("User Already Enrolled", 400), http.StatusBadRequest)
		}

		{
			if requestBody.Representation == "" {
				return u.HandleError(c, errors.New("Missing Fields! Representation is required", 400), http.StatusBadRequest)
			}
			update.Representation = requestBody.Representation

			if requestBody.Gender == "" {
				return u.HandleError(c, errors.New("Missing Fields! Gender is required", 400), http.StatusBadRequest)
			}
			update.Gender = requestBody.Gender

			if requestBody.AgeGroup == "" {
				return u.HandleError(c, errors.New("Missing Fields! Age Group is required", 400), http.StatusBadRequest)
			}
			update.AgeGroup = requestBody.AgeGroup

			if requestBody.EmploymentStatus == "" {
				return u.HandleError(c, errors.New("Missing Fields! Employment Status is required", 400), http.StatusBadRequest)
			}
			update.EmploymentStatus = requestBody.EmploymentStatus

			if requestBody.HighestSchool == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose Highest Education", 400), http.StatusBadRequest)
			}
			update.HighestSchool = requestBody.HighestSchool

			if requestBody.CanWorkInUSA == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose if you can work in USA", 400), http.StatusBadRequest)
			}
			update.CanWorkInUSA = requestBody.CanWorkInUSA

			if requestBody.LearningTrack == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose a Learning Track", 400), http.StatusBadRequest)
			}
			update.LearningTrack = requestBody.LearningTrack

			if requestBody.TechExperience == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please specify Tech Experience", 400), http.StatusBadRequest)
			}
			update.TechExperience = requestBody.TechExperience

			if requestBody.HoursPerWeek == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose Hours available Per Week", 400), http.StatusBadRequest)
			}
			update.HoursPerWeek = requestBody.HoursPerWeek

			if requestBody.Referral == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose your Referral", 400), http.StatusBadRequest)
			}
			update.Referral = requestBody.Referral

			if requestBody.OptionalMajor != "" {
				// major is optional
				update.OptionalMajor = requestBody.OptionalMajor
			}
		}

		update.Enrolled = true
		user, err := userUpdater.UpdateUser(ctx, *update)
		if err != nil {
			return u.HandleError(c, err, errors.CodeFrom(err))
		}

		// welcome the user
		err = emailer.Welcome(ctx, user)
		if err != nil {
			u.logger.Err(err).Msgf("Failed to send welcome email to '%s'", user.Email)
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
