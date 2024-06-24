package backup

import (
	"context"

	"github.com/Polo44444/harpo/alerting"
	"github.com/Polo44444/harpo/config"
	"github.com/Polo44444/harpo/storing"
)

type processor interface {
	process(ctx context.Context, folder config.Folder, storages map[string]storing.Provider, notifiers map[string]alerting.Provider)

	// setNext sets the next processor in the chain and returns it
	setNext() processor
}
