package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

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
		tmpl   *template.Template
	}
)

func NewMailChimp(apiKey string, logger zerolog.Logger) (*mailchimp, error) {
	if apiKey == "" {
		return nil, errors.New("apiKey is required", 400)
	}

	tmpl, err := template.New("welcome").Parse(welcomeHTML)
	if err != nil {
		return nil, err
	}

	return &mailchimp{
		apiKey: apiKey,
		tmpl:   tmpl,
	}, nil
}

func (m *mailchimp) Welcome(ctx context.Context, user *model.User) error {
	endpoint := "https://mandrillapp.com/api/1.0/messages/send"

	var buf bytes.Buffer
	if err := m.tmpl.Execute(&buf, user); err != nil {
		return err
	}

	payload := map[string]interface{}{
		"key": m.apiKey,
		"message": map[string]interface{}{
			"html":       buf.String(),
			"subject":    "Welcome to ReskillAmericans",
			"from_email": "info@reskillamericans.org",
			"to": []map[string]interface{}{
				{
					"email": user.Email,
					"name":  user.Name,
				},
			},
			"bcc_address": "info@reskillamericans.org",
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != 200 {
		// dump the response body for debugging
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fmt.Println(buf.String())
		return fmt.Errorf("failed to send email: %d", resp.StatusCode)
	}

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

/*
curl -X POST \
  https://mandrillapp.com/api/1.0/messages/send-template \
  -d '{"key":"","template_name":"","template_content":[],"message":{"html":"","text":"","subject":"","from_email":"","from_name":"","to":[],"headers":{},"important":false,"track_opens":false,"track_clicks":false,"auto_text":false,"auto_html":false,"inline_css":false,"url_strip_qs":false,"preserve_recipients":false,"view_content_link":false,"bcc_address":"","tracking_domain":"","signing_domain":"","return_path_domain":"","merge":false,"merge_language":"mailchimp","global_merge_vars":[],"merge_vars":[],"tags":[],"subaccount":"","google_analytics_domains":[],"google_analytics_campaign":"","metadata":{"website":""},"recipient_metadata":[],"attachments":[],"images":[]},"async":false,"ip_pool":"","send_at":""}'
*/
