package archiving

import (
	"context"
	"fmt"
	"io"

	"github.com/Polo44444/harpo/models"
	"github.com/mholt/archiver/v4"
)

type CompressionType string

const (
	GzCompressionType CompressionType = "Gz"
)

type tarProvider struct {
	compression string // Compression method: Gz
	method      int    // Method/level of compression
}

func BuildTarConfig(method int, compression CompressionType) models.ProviderConfig {
	return models.ProviderConfig{
		"compression": string(compression),
		"method":      method,
	}
}

func newTarProvider(config models.ProviderConfig) (*tarProvider, error) {

	prvd := &tarProvider{
		compression: config["compression"].(string),
		method:      config["method"].(int),
	}

	if prvd.compression != string(GzCompressionType) {
		return nil, fmt.Errorf("invalid compression type: %s", prvd.compression)
	}

	// check if the method is valid
	if prvd.method < 0 || prvd.method > 9 { // Gzip compression methods
		return nil, fmt.Errorf("invalid compression method: %d", prvd.method)
	}

	return prvd, nil
}

// Archive creates a tar archive from the src and writes it to the dst.
func (t *tarProvider) Archive(ctx context.Context, src string, dst io.Writer, ignoreErrors bool) error {

	files, err := archiver.FilesFromDisk(nil, map[string]string{
		src: "",
	})
	if err != nil {
		return err
	}

	var compression archiver.Compression

	switch t.compression {
	case string(GzCompressionType):
		compression = archiver.Gz{
			CompressionLevel: t.method,
			Multithreaded:    true,
		}
	default:
		compression = archiver.Gz{
			CompressionLevel: t.method,
			Multithreaded:    true,
		}
	}

	format := archiver.CompressedArchive{
		Compression: compression,
		Archival: archiver.Tar{
			ContinueOnError: ignoreErrors,
		},
	}

	return format.Archive(ctx, dst, files)
}

// Extract extracts the tar archive from the src and writes it to the dst. The dst must be a directory.
func (t *tarProvider) Extract(ctx context.Context, src io.Reader, dst string, ignoreErrors bool) error {

	format := archiver.CompressedArchive{
		Archival: archiver.Tar{
			ContinueOnError: ignoreErrors,
		},
	}

	return format.Extract(ctx, src, nil, NewExtract(dst).handler)
}
