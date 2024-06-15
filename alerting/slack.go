package alerting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type slackProvider struct {
	webhookURL string
}

type SlackMessage struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type   string      `json:"type"`
	Text   *BlockText  `json:"text,omitempty"`
	Fields []BlockText `json:"fields,omitempty"`
}

type BlockText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func BuildSlackConfig(webhookURL string) ProviderConfig {
	return ProviderConfig{"webhook_url": webhookURL}
}

func newSlackProvider(config ProviderConfig) (*slackProvider, error) {

	prvd := &slackProvider{
		webhookURL: config["webhook_url"].(string),
	}

	return prvd, nil
}

func (s *slackProvider) Send(ctx context.Context, m *Message) error {

	// Validate the message
	err := m.Validate()
	if err != nil {
		return err
	}

	details := ""
	if m.Details != nil {
		details = m.Details.Error()
	}

	// Send the message
	msg := SlackMessage{
		Blocks: []Block{
			{
				Type: "section",
				Text: &BlockText{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*%s*", m.Subject),
				},
			},
			{
				Type: "section",
				Fields: []BlockText{
					{Type: "mrkdwn", Text: fmt.Sprintf("*User:*\n%s", m.Entity)},
					{Type: "mrkdwn", Text: fmt.Sprintf("*Location:*\n%s", m.LocationToString())},
					{Type: "mrkdwn", Text: fmt.Sprintf("*Level:*\n%s", string(m.Level))},
					{Type: "mrkdwn", Text: fmt.Sprintf("*Details:*\n%s", details)},
				},
			},
		},
	}

	reqBodyBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.webhookURL, bytes.NewReader(reqBodyBytes))
	if err != nil {
		return err
	}

	// Headers
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	// Errors
	if err != nil && err == context.Canceled {
		return fmt.Errorf("request canceled. %w", err)
	}

	if err != nil && err == context.DeadlineExceeded {
		return fmt.Errorf("request deadline exceeded. %w", err)
	}

	if err != nil {
		return fmt.Errorf("failed to send message. %w", err)
	}

	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return fmt.Errorf("failed to send message. status Code: %d", resp.StatusCode)
	}

	return nil
}

func (s *slackProvider) Close(_ context.Context) error {
	return nil
}
