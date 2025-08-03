package notification_test

import (
	"os"
	"testing"

	"github.com/ahobsonsayers/twitchets/notification"
	"github.com/ahobsonsayers/twitchets/test"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestSqsSendTicketMessage(t *testing.T) {
	t.Skip("Can only be run manually locally with environment variables set. Comment to run.")

	_ = godotenv.Load(test.ProjectDirectoryJoin(t, ".env"))

	queueUrl := os.Getenv("SQS_QUEUE_URL")
	require.NotEmpty(t, queueUrl, "SQS_QUEUE_URL is not set")

	region := os.Getenv("AWS_REGION")
	require.NotEmpty(t, region, "AWS_REGION is not set")

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	require.NotEmpty(t, accessKeyID, "AWS_ACCESS_KEY_ID is not set")

	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	require.NotEmpty(t, secretAccessKey, "AWS_SECRET_ACCESS_KEY is not set")

	client, err := notification.NewSqsClient(notification.SqsConfig{
		QueueUrl:        queueUrl,
		Region:          region,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	})
	require.NoError(t, err)

	ticket := testNotificationTicket()
	err = client.SendTicketNotification(ticket)
	require.NoError(t, err)
}
