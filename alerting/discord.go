package alerting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Polo44444/harpo/models"
)

type discordProvider struct {
	webhookURL string
}

type discordEmbedFields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type discordEmbed struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Color       int                  `json:"color"`
	Fields      []discordEmbedFields `json:"fields"`
}

type discordMessage struct {
	Content string         `json:"content"`
	Embeds  []discordEmbed `json:"embeds"`
}

func BuildDiscordConfig(webhookURL string) models.ProviderConfig {
	return models.ProviderConfig{"webhook_url": webhookURL}
}

func newDiscordProvider(config models.ProviderConfig) (*discordProvider, error) {

	prvd := &discordProvider{
		webhookURL: config["webhook_url"].(string),
	}

	return prvd, nil
}

func (d *discordProvider) levelToColor(level MessageLevel) int {
	switch level {
	case DebugMessage:
		return 0xeb459e
	case InfoMessage:
		return 0x3498db
	case WarningMessage:
		return 0xf1c40f
	case ErrorMessage:
		return 0xed4245
	case FatalMessage:
		return 0xed4245
	default:
		return 0xed4245
	}
}

func (d *discordProvider) Send(ctx context.Context, m *Message) error {

	// Validate the message
	err := m.Validate()
	if err != nil {
		return err
	}

	mErr := ""
	if m.Err != nil {
		mErr = m.Err.Error()
	}

	// Send the message
	msg := discordMessage{
		Content: m.Subject,
		Embeds: []discordEmbed{
			{
				Title:       m.Entity,
				Description: m.Details,
				Color:       d.levelToColor(m.Level),
				Fields: []discordEmbedFields{
					{
						Name:   "Level",
						Value:  string(m.Level),
						Inline: true,
					},
					{
						Name:   "Extras",
						Value:  m.ExtrasToString(),
						Inline: false,
					},
				},
			},
		},
	}

	if mErr != "" {
		msg.Embeds[0].Fields = append(msg.Embeds[0].Fields, discordEmbedFields{
			Name:   "Error",
			Value:  mErr,
			Inline: false,
		})
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", d.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return httpRequestError(err)
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return fmt.Errorf("failed to send discord message. Status Code: %s", resp.Status)
	}

	return nil
}

func (d *discordProvider) Close(_ context.Context) error {
	return nil
}
