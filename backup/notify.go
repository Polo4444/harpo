package backup

import (
	"context"

	"github.com/Polo44444/harpo/alerting"
)

func notify(ctx context.Context, m *alerting.Message, notifiers map[string]alerting.Provider) {

	for _, notifier := range notifiers {

		mClone := *m
		go notifier.Send(ctx, &mClone)
	}
}

func NotifyInfo(ctx context.Context, folderName, text, details string, notifiers map[string]alerting.Provider) {

	m := &alerting.Message{
		Entity:  "Harpo Backup",
		Subject: text,
		Level:   alerting.InfoMessage,
		Extras:  []string{folderName},
		Details: details,
	}

	notify(ctx, m, notifiers)
}

// NotifyError sends an error message to all notifiers
func NotifyError(ctx context.Context, folderName, text, details string, err error, notifiers map[string]alerting.Provider) {

	m := &alerting.Message{
		Entity:  "Harpo Backup",
		Subject: text,
		Level:   alerting.ErrorMessage,
		Extras:  []string{folderName},
		Details: details,
		Err:     err,
	}

	notify(ctx, m, notifiers)
}
