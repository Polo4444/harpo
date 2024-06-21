package alerting

import (
	"context"
	"testing"
	"time"

	"github.com/Polo44444/harpo/models"
)

var testDiscordConf = models.ProviderConfig{
	"webhook_url": "webhook_url",
}

func TestSendWithDiscord(t *testing.T) {

	// timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// We create init discord provider
	p, err := GetProvider(DiscordProvider, testDiscordConf)
	if err != nil {
		t.Fatalf("Error creating discord provider: %s", err.Error())
	}
	defer p.Close(ctx)

	// We send the message
	err = p.Send(ctx, testMessage)
	if err != nil {
		t.Fatalf("Error sending message: %s", err.Error())
	}

	t.Log("Message sent successfully\n")
}
