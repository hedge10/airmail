package mail

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type (
	Ses struct {
		AwsAccessKeyId     string
		AwsSecretAccessKey string
		Region             string
	}
)

func (s Ses) Send(e *Email) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.Region),
		Credentials: credentials.NewStaticCredentials(s.AwsAccessKeyId, s.AwsSecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	to := make([]*string, 0, len(e.To))
	for _, t := range e.To {
		to = append(to, &t.Email)
	}

	svc := ses.New(sess)
	_, sending_err := svc.SendEmail(&ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: to,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(e.Message),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(e.Subject),
			},
		},
		Source: aws.String(e.From.Email),
	})

	return sending_err
}
