/* //go:build integration_test */

package alerting

var testMessage = NewMessage("sespolo@gmail.com", ErrorMessage, []string{"info 1", "info 2"}, "test subject", "test details", nil)
