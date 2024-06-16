package alerting

import (
	"context"
	"fmt"
)

func httpRequestError(err error) error {

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

	return err
}
