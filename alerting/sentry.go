package alerting

import (
	"context"
	"errors"
	"time"

	"github.com/Polo44444/harpo/models"
	"github.com/getsentry/sentry-go"
)

type sentryProvider struct {
	dsn string
}

func BuildSentryConfig(dsn string) models.ProviderConfig {
	return models.ProviderConfig{
		"dsn": dsn,
	}
}

func newSentryProvider(config models.ProviderConfig) (*sentryProvider, error) {

	prvd := &sentryProvider{
		dsn: config["dsn"].(string),
	}

	// Init sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn: prvd.dsn,
	})

	if err != nil {
		return nil, err
	}

	return prvd, nil
}

func (s *sentryProvider) ToSentryLevel(level MessageLevel) sentry.Level {
	switch level {
	case DebugMessage:
		return sentry.LevelDebug
	case InfoMessage:
		return sentry.LevelInfo
	case SuccessMessage:
		return sentry.LevelInfo
	case WarningMessage:
		return sentry.LevelWarning
	case ErrorMessage:
		return sentry.LevelError
	case FatalMessage:
		return sentry.LevelFatal
	default:
		return sentry.LevelError
	}
}

func (s *sentryProvider) Send(_ context.Context, m *Message) error {

	// Validate the message
	err := m.Validate()
	if err != nil {
		return err
	}

	// Build the event
	ev := sentry.NewEvent()
	ev.Message = m.Subject

	// create eexception with the error
	ev.Exception = []sentry.Exception{{
		Type:       m.Subject,
		Value:      m.Details,
		Stacktrace: sentry.ExtractStacktrace(m.Err),
	}}

	ev.Level = s.ToSentryLevel(m.Level)
	ev.User = sentry.User{ID: m.Entity}
	ev.Extra["entity"] = m.Entity
	ev.Extra["extras"] = m.ExtrasToString()
	if m.Err != nil {
		ev.Extra["error"] = m.Err.Error()
	}
	ev.Extra["level"] = m.Level

	// Send the event
	sentry.CaptureEvent(ev)
	return nil
}

func (s *sentryProvider) Close(ctx context.Context) error {

	// Get conetxt deadline
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(2 * time.Second)
	}

	if !sentry.Flush(time.Until(deadline)) {
		return errors.New("failed to flush sentry")
	}

	return nil
}
