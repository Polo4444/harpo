/* //go:build integration_test */

package alerting

import (
	"errors"
	"testing"
)

var testMessage = NewMessage("test@domain.com", ErrorMessage, []string{"func1", "line1"}, "test subject", errors.New("test details"))

func TestInit(t *testing.T) {

	testSlackConf = BuildSlackConfig(
		"webhook_url",
	)

	testSentryConf = BuildSentryConfig(
		"dsn",
	)
}
