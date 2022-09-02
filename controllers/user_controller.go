package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"unicode"

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
		// var body bytes.Buffer
		// _, err := io.Copy(&body, c.Request().Body)
		// if err != nil {
		// 	u.logger.Debug().Msgf("Error reading request body: %s", err)
		// }

		// cp := body.Bytes()
		// u.logger.Debug().Msgf("Request Body: %s", cp)

		err := json.NewDecoder(c.Request().Body).Decode(&requestBody)
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

		profile, err := service.GetProfile(authCode, redirectURI)
		if err != nil {
			return u.HandleError(c, err, errors.CodeFrom(err))
		}

		// do validations
		if profile.Name == "" {
			return u.HandleError(c, errors.New("Invalid Profile. Found No Name", 400), http.StatusBadRequest)
		}

		if profile.Photo == "" {
			return u.HandleError(c, errors.New("Invalid Profile. Please Set Your Profile Picture on LinkedIn", 400), http.StatusBadRequest)
		}

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
		// var body bytes.Buffer
		// _, err := io.Copy(&body, c.Request().Body)
		// if err != nil {
		// 	u.logger.Debug().Msgf("Error reading request body: %s", err)
		// }

		// cp := body.Bytes()
		// u.logger.Debug().Msgf("Request Body: %s", cp)

		err := json.NewDecoder(c.Request().Body).Decode(&requestBody)
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
			// if requestBody.Timezone != "" {
			// 	update.Timezone = requestBody.Timezone
			// }

			if requestBody.LinkedInURL == "" {
				return u.HandleError(c, errors.New("Missing Fields! LinkedIn URL is required", 400), http.StatusBadRequest)
			}
			match, err := isValidLinkedIn(requestBody.LinkedInURL)
			if err != nil {
				return u.HandleError(c, err, errors.CodeFrom(err))
			}
			if !match {
				return u.HandleError(c, errors.New("Invalid LinkedIn URL", 400), http.StatusBadRequest)
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
			if strings.ToUpper(requestBody.CanWorkInUSA) != "YES" {
				return u.HandleError(c, errors.New("It is Required that You can Work in the USA", 400), http.StatusBadRequest)
			}
			update.CanWorkInUSA = requestBody.CanWorkInUSA

			if requestBody.LearningTrack == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose a Learning Track", 400), http.StatusBadRequest)
			}
			update.LearningTrack = requestBody.LearningTrack

			if requestBody.HoursPerWeek == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose Hours available Per Week", 400), http.StatusBadRequest)
			}
			update.HoursPerWeek = requestBody.HoursPerWeek

			if requestBody.Referral == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose your Referral", 400), http.StatusBadRequest)
			}
			update.Referral = requestBody.Referral

			if requestBody.Photo == "" || requestBody.Photo == "null" {
				return u.HandleError(c, errors.New("Missing Fields! Please upload a picture", 400), http.StatusBadRequest)
			}
			update.Photo = requestBody.Photo

			if requestBody.City == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please set a City", 400), http.StatusBadRequest)
			}
			update.City = requestBody.City

			// if requestBody.State == "" {
			// 	return u.HandleError(c, errors.New("Missing Fields! Please set a State", 400), http.StatusBadRequest)
			// }
			// update.State = requestBody.State

			if requestBody.ProfessionalExperience == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose a Professional Experience", 400), http.StatusBadRequest)
			}
			update.ProfessionalExperience = requestBody.ProfessionalExperience

			//industries := strings.Split(requestBody.Industries, ",")
			//if len(industries) < 1 {
			//	return u.HandleError(c, errors.New("Missing Fields! Please add at least one Industry", 400), http.StatusBadRequest)
			//}
			//for _, industry := range industries {
			//	if !isAlpha(industry) {
			//		return u.HandleError(c, errors.New("One of the Industries is invalid", 400), http.StatusBadRequest)
			//	}
			//}
			if err := validateIndustries(requestBody.Industries); err != nil {
				return u.HandleError(c, err, http.StatusBadRequest)
			}
			//industriesStr := strings.Join(industries, ",")
			update.Industries = requestBody.Industries

			//if requestBody.RacialDemographic == "" {
			//	return u.HandleError(c, errors.New("Missing Fields! Please choose a Racial Demographic", 400), http.StatusBadRequest)
			//}
			//update.RacialDemographic = requestBody.RacialDemographic

			if requestBody.PriorKnowledge == "" {
				return u.HandleError(c, errors.New("Missing Fields! Please choose Prior Knowledge level", 400), http.StatusBadRequest)
			}
			update.PriorKnowledge = requestBody.PriorKnowledge

			if requestBody.ReferralOther != "" {
				// referralOther is optional
				update.ReferralOther = requestBody.ReferralOther
			}

			if requestBody.OptionalMajor != "" {
				return u.HandleError(c, errors.New("Missing Fields! Please add a Field of Study", 400), http.StatusBadRequest)
			}
			if !isAlpha(requestBody.OptionalMajor) {
				return u.HandleError(c, errors.New("Invalid Field of Study", 400), http.StatusBadRequest)
			}
			update.OptionalMajor = requestBody.OptionalMajor

		}

		update.Enrolled = true
		user, err := userUpdater.UpdateUser(ctx, *update)
		if err != nil {
			return u.HandleError(c, err, errors.CodeFrom(err))
		}

		if err := emailer.Welcome(ctx, user); err != nil {
			u.logger.Error().Err(err).Msg("failed to send welcome email")
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

func isValidLinkedIn(url string) (bool, error) {
	validRoot1 := "https://www.linkedin.com/in/"
	validRoot2 := "https://linkedin.com/in/"

	var (
		hasRoot   bool
		afterRoot string
	)

	if strings.HasPrefix(url, validRoot1) {
		hasRoot = true
		afterRoot = url[len(validRoot1):]

	} else if strings.HasPrefix(url, validRoot2) {
		hasRoot = true
		afterRoot = url[len(validRoot2):]
	}

	if !hasRoot {
		return false, nil //errors.New("Invalid LinkedIn URL", 400)
	}

	if strings.TrimSpace(afterRoot) == "" {
		return false, nil //errors.New("Invalid LinkedIn URL", 400)
	}

	return true, nil
}

func isAlphaNum(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func isAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func isNum(s string) bool {
	for _, r := range s {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func validateIndustries(s string) error {
	industries := strings.Split(s, ",")
	if len(industries) < 1 {
		return errors.New("Missing Fields! Please add at least one Industry", 400)
	}
	for _, industry := range industries {
		industry = strings.TrimSpace(industry)
		if !isAlpha(industry) {
			return errors.New("One of the Industries is invalid", 400)
		}
	}

	return nil
}
