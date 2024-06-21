package alerting

import (
	"context"
	"fmt"

	"github.com/Polo44444/harpo/models"
)

const (
	SentryProvider  models.ProviderEntity = "SENTRY"
	SlackProvider   models.ProviderEntity = "SLACK"
	DiscordProvider models.ProviderEntity = "DISCORD"
)

// Errors
var (
	ErrProviderNotSupported = fmt.Errorf("provider not supported")
	ErrNoSubjectProvided    = fmt.Errorf("no message subject provided")
	ErrNoLocationProvided   = fmt.Errorf("no message location provided")
)

// Define the Provider interface
type Provider interface {
	Send(ctx context.Context, m *Message) error
	Close(ctx context.Context) error
}

// GetProvider returns a provider based on the entity and the config
func GetProvider(entity models.ProviderEntity, config models.ProviderConfig) (Provider, error) {

	var err error = nil
	var prvd Provider = nil

	switch entity {
	case SentryProvider:
		prvd, err = newSentryProvider(config)
	case SlackProvider:
		prvd, err = newSlackProvider(config)
	case DiscordProvider:
		prvd, err = newDiscordProvider(config)
	default:
		err = ErrProviderNotSupported
	}

	return prvd, err
}
