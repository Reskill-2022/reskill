package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/email"
	"github.com/thealamu/linkedinsignin/errors"
	"github.com/thealamu/linkedinsignin/linkedin"
	"github.com/thealamu/linkedinsignin/model"
	"github.com/thealamu/linkedinsignin/repository"
	"github.com/thealamu/linkedinsignin/requests"
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
			return u.HandleError(c, errors.New("Invalid JSON Request Body", 400), http.StatusBadRequest)
		}

		authCode := requestBody.AuthCode
		if authCode == "" {
			return u.HandleError(c, errors.New("Auth Code is required", 400), http.StatusBadRequest)
		}
		redirectURI := requestBody.RedirectURI
		if redirectURI == "" {
			return u.HandleError(c, errors.New("Redirect URI is required", 400), http.StatusBadRequest)
		}
		// userEmail := strings.ToLower(requestBody.Email)
		// if userEmail == "" {
		// 	return u.HandleError(c, errors.New("Email is required", 400), http.StatusBadRequest)
		// }

		profile, err := service.GetProfile(authCode, redirectURI)
		if err != nil {
			return u.HandleError(c, err, errors.CodeFrom(err))
		}

		// do validations
		if profile.Name == "" {
			return u.HandleError(c, errors.New("Invalid Profile. Found No Name", 400), http.StatusBadRequest)
		}
		// if profile.Location == "" {
		// 	return u.HandleError(c, errors.New("Invalid Profile. Please Set Your City and State of Residence on LinkedIn", 400), http.StatusBadRequest)
		// }
		if profile.Photo == "" {
			return u.HandleError(c, errors.New("Invalid Profile. Please Set Your Profile Picture on LinkedIn", 400), http.StatusBadRequest)
		}
		// if !profile.HasExperience {
		// 	return u.HandleError(c, errors.New("Invalid Profile. Please Add Your Work Experience on LinkedIn", 400), http.StatusBadRequest)
		// }

		//country := profile.Location
		//i := strings.LastIndex(profile.Location, ",")
		//if i > 0 {
		//	country = profile.Location[i+2:]
		//}
		//if country != "United States" {
		//	return u.HandleError(c, errors.New("Invalid Profile. For United States Only", 400), http.StatusBadRequest)
		//}

		firstname, lastname := u.splitNames(profile.Name)

		data := model.User{
			Email:       profile.Email,
			Name:        profile.Name,
			FirstName:   firstname,
			LastName:    lastname,
			LinkedInURL: profile.ProfileURL,
			Location:    profile.Location,
			Phone:       profile.Phone,
			Photo:       profile.Photo,
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
			if requestBody.Timezone != "" {
				update.Timezone = requestBody.Timezone
			}

			if requestBody.LinkedInURL == "" {
				return u.HandleError(c, errors.New("Missing Fields! LinkedIn URL is required", 400), http.StatusBadRequest)
			}
			update.LinkedInURL = requestBody.LinkedInURL

			if requestBody.Phone == "" {
				return u.HandleError(c, errors.New("Missing Fields! Phone Number is required", 400), http.StatusBadRequest)
			}
			update.Phone = requestBody.Phone

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
			if strings.Title(requestBody.CanWorkInUSA) != "Yes" {
				return u.HandleError(c, errors.New("It is Required that You can Work in the USA", 400), http.StatusBadRequest)
			}
			update.CanWorkInUSA = requestBody.CanWorkInUSA

			if requestBody.LearningTrack == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose a Learning Track", 400), http.StatusBadRequest)
			}
			update.LearningTrack = requestBody.LearningTrack

			if requestBody.TechExperience != "" {
				update.TechExperience = requestBody.TechExperience
			}

			if requestBody.HoursPerWeek == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose Hours available Per Week", 400), http.StatusBadRequest)
			}
			update.HoursPerWeek = requestBody.HoursPerWeek

			if requestBody.Referral == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose your Referral", 400), http.StatusBadRequest)
			}
			update.Referral = requestBody.Referral

			if requestBody.Photo == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please upload a picture", 400), http.StatusBadRequest)
			}
			update.Photo = requestBody.Photo

			if requestBody.ReferralOther != "" {
				// referralOther is optional
				update.ReferralOther = requestBody.ReferralOther
			}

			if requestBody.OptionalMajor != "" {
				// major is optional
				update.OptionalMajor = requestBody.OptionalMajor
			}

			if requestBody.City != "" {
				update.City = requestBody.City
			}

			if requestBody.State != "" {
				update.State = requestBody.State
			}
			if requestBody.ProfessionalExperience != "" {
				update.ProfessionalExperience = requestBody.ProfessionalExperience
			}
			if requestBody.Industries != "" {
				update.Industries = requestBody.Industries
			}
			if requestBody.WillChangeJob != "" {
				update.WillChangeJob = requestBody.WillChangeJob
			}
			if requestBody.WillChangeJobRole != "" {
				update.WillChangeJobRole = requestBody.WillChangeJobRole
			}
			if requestBody.OpenToMeet != "" {
				update.OpenToMeet = requestBody.OpenToMeet
			}
			if requestBody.RacialDemographic != "" {
				update.RacialDemographic = requestBody.RacialDemographic
			}
			if requestBody.PriorKnowledge != "" {
				update.PriorKnowledge = requestBody.PriorKnowledge
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

		userEmail := c.Param("email")
		if userEmail == "" {
			return u.HandleError(c, errors.New(" Email is required", 400), http.StatusBadRequest)
		}

		user, err := userGetter.GetUser(ctx, userEmail)
		if err != nil {
			return u.HandleError(c, err, errors.CodeFrom(err))
		}

		return HandleSuccess(c, user, http.StatusOK)
	}
}

func (u *UserController) splitNames(name string) (string, string) {
	names := strings.Split(name, " ")
	if len(names) == 1 {
		return names[0], ""
	}
	return names[0], names[len(names)-1]
}
