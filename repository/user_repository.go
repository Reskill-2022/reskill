package repository

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/errors"
	"github.com/thealamu/linkedinsignin/model"
	"google.golang.org/api/option"
)

type UserRepository struct {
	logger  zerolog.Logger
	client1 *firestore.Client
	client2 *firestore.Client
}

var _ UserRepositoryInterface = (*UserRepository)(nil)

func NewUserRepository(logger zerolog.Logger) *UserRepository {
	r := &UserRepository{
		logger: logger,
	}

	r.client1 = getClient("./service-account-1.json")
	r.client2 = getClient("./service-account-2.json")

	return r
}

func getClient(saFile string) *firestore.Client {
	ctx := context.Background()

	sa := option.WithCredentialsFile(saFile)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return client
}

func (u *UserRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	u.logger.Debug().Msgf("Firestore: creating user with email: %s", user.Email)

	gotUser, err := u.GetUser(ctx, user.Email)
	if err == nil || gotUser != nil {
		return gotUser, nil
	}

	_, err = u.client.Collection("users").Doc(user.Email).Set(ctx, user)
	if err != nil {
		return nil, errors.From(err, "failed to create user", 500)
	}
	return &user, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, user model.User) (*model.User, error) {
	u.logger.Debug().Msgf("Firestore: updating user with email: %s", user.Email)

	_, err := u.client.Collection("users").Doc(user.Email).Update(ctx, []firestore.Update{
		{Path: "representation", Value: user.Representation},
		{Path: "gender", Value: user.Gender},
		{Path: "age_group", Value: user.AgeGroup},
		{Path: "employment_status", Value: user.EmploymentStatus},
		{Path: "highest_school", Value: user.HighestSchool},
		{Path: "optional_major", Value: user.OptionalMajor},
		{Path: "can_work_in_usa", Value: user.CanWorkInUSA},
		{Path: "learning_track", Value: user.LearningTrack},
		// {Path: "tech_experience", Value: user.TechExperience},
		{Path: "hours_per_week", Value: user.HoursPerWeek},
		{Path: "referral", Value: user.Referral},
		{Path: "referral_other", Value: user.ReferralOther},
		{Path: "enrolled", Value: user.Enrolled},
		{Path: "timezone", Value: user.Timezone},
		{Path: "phone", Value: user.Phone},
		{Path: "photo", Value: user.Photo},
		{Path: "gitaccount", Value: user.GitAccount},
		{Path: "figmaaccount", Value: user.FigmaAccount},
		{Path: "git_yes", Value: user.GitYes},
		{Path: "figma_yes", Value: user.FigmaYes},
		{Path: "city", Value: user.City},
		{Path: "state", Value: user.State},
		{Path: "professional_experience", Value: user.ProfessionalExperience},
		{Path: "industries", Value: user.Industries},
		// {Path: "will_change_job", Value: user.WillChangeJob},
		// {Path: "will_change_job_role", Value: user.WillChangeJobRole},
		// {Path: "open_to_meet", Value: user.OpenToMeet},
		{Path: "racial_demographic", Value: user.RacialDemographic},
		{Path: "prior_knowledge", Value: user.PriorKnowledge},
		{Path: "linkedin_url", Value: user.LinkedInURL},
	})
	if err != nil {
		return nil, errors.From(err, "failed to update user data", 500)
	}
	return &user, nil
}

func (u *UserRepository) GetUser(ctx context.Context, email string) (*model.User, error) {
	u.logger.Debug().Msgf("Firestore: getting user with email: %s", email)

	data, err := u.client.Collection("users").Doc(email).Get(ctx)
	if err != nil {
		return nil, errors.From(err, "User Account Not Found", 404)
	}

	user := model.User{}
	err = data.DataTo(&user)
	if err != nil {
		return nil, errors.From(err, "failed to bind user data", 500)
	}

	return &user, nil
}
