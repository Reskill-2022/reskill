package email

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/constants"
	"github.com/thealamu/linkedinsignin/errors"
	"github.com/thealamu/linkedinsignin/model"
)

type (
	Emailer interface {
		// Welcome sends a welcome email to the user
		Welcome(ctx context.Context, user *model.User) error
	}

	sesEmailer struct {
		client *ses.Client
		logger zerolog.Logger
	}

	mailchimp struct {
		apiKey string
	}
)

func NewMailChimp(apiKey string, logger zerolog.Logger) (*mailchimp, error) {
	if apiKey == "" {
		return nil, errors.New("apiKey is required", 400)
	}

	return &mailchimp{apiKey: apiKey}, nil
}

func (m *mailchimp) Welcome(ctx context.Context, user *model.User) error {
	return nil
}

func NewSES(ctx context.Context, logger zerolog.Logger) (*sesEmailer, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	return &sesEmailer{
		client: ses.NewFromConfig(cfg),
		logger: logger,
	}, nil
}

func (s *sesEmailer) Welcome(ctx context.Context, user *model.User) error {
	s.logger.Info().Msgf("Sending welcome email to '%s'", user.Email)

	payload := fmt.Sprintf(`{"name": "%s"}`, user.Name)

	dst := types.Destination{
		ToAddresses: []string{user.Email},
	}

	_, err := s.client.SendTemplatedEmail(ctx, &ses.SendTemplatedEmailInput{
		Destination:  &dst,
		Source:       aws.String(constants.DefaultSourceEmail),
		Template:     aws.String("basic-welcome"),
		TemplateData: &payload,
	})
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
