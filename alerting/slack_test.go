/* //go:build integration_test */

package alerting

import (
	"context"
	"testing"
	"time"
)

var testSlackConf = BuildSlackConfig(
	"webhook_url",
)

func TestSendWithSlack(t *testing.T) {

	// Init
	TestInit(t)

	// timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// We create init slack provider
	p, err := GetProvider(SlackProvider, testSlackConf)
	if err != nil {
		t.Fatalf("Error creating Slack provider: %s", err.Error())
	}
	defer p.Close(ctx)

	// We send the message
	err = p.Send(ctx, testMessage)
	if err != nil {
		t.Fatalf("Error sending message: %s", err.Error())
	}

	t.Log("Message sent successfully\n")
}
