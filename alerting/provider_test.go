/* //go:build integration_test */

package alerting

import (
	"errors"
)

var testMessage = NewMessage("sespolo@gmail.com", ErrorMessage, []string{"info 1", "info 2"}, "test subject", errors.New("test details"))
