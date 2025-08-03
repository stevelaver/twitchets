package notification

import (
	"context"
	"fmt"

	"github.com/ahobsonsayers/twigots"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SqsClient struct {
	queueUrl string
	client   *sqs.Client
}

var _ Client = SqsClient{}

func (c SqsClient) SendTicketNotification(ticket twigots.TicketListing) error {
	message, err := RenderMessage(ticket, WithHeader(), WithFooter())
	if err != nil {
		return err
	}

	_, err = c.client.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    &c.queueUrl,
		MessageBody: aws.String(message),
	})
	if err != nil {
		return fmt.Errorf("failed to send sqs message: %w", err)
	}

	return nil
}

type SqsConfig struct {
	QueueUrl        string `json:"queueUrl"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

func NewSqsClient(cfg SqsConfig) (SqsClient, error) {
	loadOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.Region),
	}
	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		loadOptions = append(loadOptions, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), loadOptions...)
	if err != nil {
		return SqsClient{}, fmt.Errorf("failed to load aws config: %w", err)
	}

	sqsClient := sqs.NewFromConfig(awsCfg)

	return SqsClient{
		queueUrl: cfg.QueueUrl,
		client:   sqsClient,
	}, nil
}
