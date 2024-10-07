package config

import (
	"github.com/ahobsonsayers/twitchets/cmd/twitchets/notification"
)

type NotificationConfig struct {
	Ntfy     *notification.NtfyConfig     `json:"ntfy"`
	Telegram *notification.TelegramConfig `json:"telegram"`
}

func (c NotificationConfig) Clients() ([]notification.Client, error) {
	clients := []notification.Client{}

	if c.Ntfy != nil {
		ntfyClient, err := notification.NewNtfyClient(*c.Ntfy)
		if err != nil {
			return nil, err
		}
		clients = append(clients, ntfyClient)
	}

	if c.Telegram != nil {
		telegramClient, err := notification.NewTelegramClient(*c.Telegram)
		if err != nil {
			return nil, err
		}
		clients = append(clients, telegramClient)
	}

	return clients, nil
}
