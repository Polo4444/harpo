package archiving

import (
	"context"
	"io"

	"github.com/Polo44444/harpo/models"
)

const (
	ZipProvider models.ProviderEntity = "ZIP"
	TarProvider models.ProviderEntity = "TAR"
)

// Provider interface
type Provider interface {

	/*Archive creates an archive from the src and writes it to the dst.
	`src` can be a file or a directory.
	`ignoreErrors` is a flag that indicates if the provider should ignore errors when creating the archive.
	*/
	Archive(ctx context.Context, src string, dst io.Writer, ignoreErrors bool) error

	// Extract extracts the archive from the src and writes it to the dst.
	// The dst must be a directory.
	Extract(ctx context.Context, src io.Reader, dst string, ignoreErrors bool) error

	// Ext returns the extension of the archive with the dot.
	Ext() string
}

// GetProvider returns a provider based on the entity and the config
func GetProvider(entity models.ProviderEntity, config models.ProviderConfig) (Provider, error) {

	var err error = nil
	var prvd Provider = nil

	switch entity {
	case ZipProvider:
		prvd, err = newZipProvider(config)
	case TarProvider:
		prvd, err = newTarProvider(config)
	default:
		err = models.ErrProviderNotSupported
	}

	return prvd, err
}
