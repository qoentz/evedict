package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"strings"
)

type MailService struct {
	client *sesv2.Client
}

func NewMailService(accessKey string, secretAccessKey string, region string) (*MailService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &MailService{
		client: sesv2.NewFromConfig(cfg),
	}, nil
}

func (m *MailService) SendContactEmail(name, email, subject, message string) error {
	emailMessage := fmt.Sprintf("From: %s (%s)\n\n%s", name, email, message)

	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String("contact@evedict.com"),
		Destination: &types.Destination{
			ToAddresses: []string{"contact@evedict.com"},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(subject),
					Charset: aws.String("UTF-8"),
				},
				Body: &types.Body{
					Text: &types.Content{
						Data:    aws.String(emailMessage),
						Charset: aws.String("UTF-8"),
					},
				},
			},
		},
	}

	_, err := m.client.SendEmail(context.TODO(), input)
	return err
}

func (m *MailService) ValidateContactForm(name, email, subject, message string) []string {
	var errors []string

	if strings.TrimSpace(email) == "" {
		errors = append(errors, "Email is required")
	} else if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		errors = append(errors, "Invalid email format")
	}

	if strings.TrimSpace(subject) == "" {
		errors = append(errors, "Subject is required")
	}

	if strings.TrimSpace(message) == "" {
		errors = append(errors, "Message is required")
	}

	return errors
}
