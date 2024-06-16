/* //go:build integration_test */

package alerting

import (
	"context"
	"testing"
	"time"
)

var testSentryConf = BuildSentryConfig(
	"dsn",
)

func TestSendWithSentry(t *testing.T) {

	// timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// We create init sentry provider
	p, err := GetProvider(SentryProvider, testSentryConf)
	if err != nil {
		t.Fatalf("Error creating Sentry provider: %s", err.Error())
	}
	defer p.Close(ctx)

	// We send the message
	err = p.Send(ctx, testMessage)
	if err != nil {
		t.Fatalf("Error sending message: %s", err.Error())
	}

	t.Log("Message sent successfully\n")
}
