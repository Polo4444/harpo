/* //go:build integration_test */

package alerting

import (
	"errors"
	"testing"
)

var testMessage = NewMessage("sespolo@gmail.com", ErrorMessage, []string{"info 1", "info 2"}, "test subject", errors.New("test details"))

func TestInit(t *testing.T) {

	testSlackConf = BuildSlackConfig(
		"webhook_url",
	)

	testSentryConf = BuildSentryConfig(
		"dsn",
	)
}
