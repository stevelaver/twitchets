package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ahobsonsayers/twitchets/notification"
	"github.com/orsinium-labs/enum"
)

type NotificationType enum.Member[string]

var (
	notificationType = enum.NewBuilder[string, NotificationType]()

	NotificationTypeNtfy     = notificationType.Add(NotificationType{"ntfy"})
	NotificationTypeGotify   = notificationType.Add(NotificationType{"gotify"})
	NotificationTypeTelegram = notificationType.Add(NotificationType{"telegram"})
	NotificationTypeSqs      = notificationType.Add(NotificationType{"sqs"})

	NotificationTypes = notificationType.Enum()
)

func (c *NotificationType) UnmarshalJSON(data []byte) error {
	var notificationTypeString string
	err := json.Unmarshal(data, &notificationTypeString)
	if err != nil {
		return err
	}

	notificationType := NotificationTypes.Parse(notificationTypeString)
	if notificationType == nil {
		return fmt.Errorf("notificationType '%s' is not valid", notificationTypeString)
	}

	*c = *notificationType
	return nil
}

func (c *NotificationType) UnmarshalText(data []byte) error {
	notificationTypeString := string(data)
	notificationType := NotificationTypes.Parse(notificationTypeString)
	if notificationType == nil {
		return fmt.Errorf("notificationType '%s' is not valid", notificationTypeString)
	}

	*c = *notificationType
	return nil
}

type NotificationConfig struct {
	Ntfy     *notification.NtfyConfig     `json:"ntfy"`
	Gotify   *notification.GotifyConfig   `json:"gotify"`
	Telegram *notification.TelegramConfig `json:"telegram"`
	Sqs      *notification.SqsConfig      `json:"sqs"`
}

func (c NotificationConfig) Validate() error {
	if c.Ntfy != nil {
		if !beginsWithHttp(c.Ntfy.Url) {
			return errors.New("ntfy url must begin with 'http://' or 'https://'")
		}
		if c.Ntfy.Topic == "" {
			return errors.New("ntfy topic must be set")
		}
	}

	if c.Gotify != nil {
		if !beginsWithHttp(c.Gotify.Url) {
			return errors.New("gotify url must begin with 'http://' or 'https://'")
		}
		if c.Gotify.Token == "" {
			return errors.New("gotify token cannot be empty")
		}
	}

	if c.Telegram != nil {
		if c.Telegram.ChatId == 0 {
			return errors.New("telegram chat id cannot be empty")
		}
		if c.Telegram.Token == "" {
			return errors.New("telegram token cannot be empty")
		}
	}

	if c.Sqs != nil {
		if c.Sqs.QueueUrl == "" {
			return errors.New("sqs queue url cannot be empty")
		}
	}

	return nil
}

func (c NotificationConfig) Clients() (map[NotificationType]notification.Client, error) {
	clients := map[NotificationType]notification.Client{}

	if c.Ntfy != nil {
		ntfyClient, err := notification.NewNtfyClient(*c.Ntfy)
		if err != nil {
			return nil, fmt.Errorf("failed to setup ntfy client: %w", err)
		}

		clients[NotificationTypeNtfy] = ntfyClient
	}

	if c.Gotify != nil {
		gotifyClient, err := notification.NewGotifyClient(*c.Gotify)
		if err != nil {
			return nil, fmt.Errorf("failed to setup gotify client: %w", err)
		}

		clients[NotificationTypeGotify] = gotifyClient
	}

	if c.Telegram != nil {
		telegramClient, err := notification.NewTelegramClient(*c.Telegram)
		if err != nil {
			return nil, fmt.Errorf("failed to setup telegram client: %w", err)
		}

		clients[NotificationTypeTelegram] = telegramClient
	}

	if c.Sqs != nil {
		sqsClient, err := notification.NewSqsClient(*c.Sqs)
		if err != nil {
			return nil, fmt.Errorf("failed to setup sqs client: %w", err)
		}

		clients[NotificationTypeSqs] = sqsClient
	}

	return clients, nil
}

func beginsWithHttp(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
