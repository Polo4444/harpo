package archiving

import (
	"archive/zip"
	"context"
	"fmt"
	"io"

	"github.com/Polo44444/harpo/models"
	"github.com/mholt/archiver/v4"
)

type zipProvider struct {
	method uint16 // Method/level of compression
}

func BuildZipConfig(method uint16) models.ProviderConfig {
	return models.ProviderConfig{
		"method": method,
	}
}

func newZipProvider(config models.ProviderConfig) (*zipProvider, error) {

	prvd := &zipProvider{
		method: config["method"].(uint16),
	}

	// check if the method is valid
	if prvd.method != zip.Deflate && prvd.method != zip.Store {
		return nil, fmt.Errorf("invalid compression method: %d", prvd.method)
	}

	return prvd, nil
}

// Archive creates a zip archive from the src and writes it to the dst.
func (z *zipProvider) Archive(ctx context.Context, src string, dst io.Writer, ignoreErrors bool) error {

	files, err := archiver.FilesFromDisk(nil, map[string]string{
		src: "",
	})
	if err != nil {
		return err
	}

	format := archiver.Zip{
		Compression:     z.method,
		ContinueOnError: ignoreErrors,
	}

	return format.Archive(ctx, dst, files)
}

// Extract extracts the zip archive from the src and writes it to the dst. The dst must be a directory.
func (z *zipProvider) Extract(ctx context.Context, src io.Reader, dst string, ignoreErrors bool) error {

	format := archiver.Zip{
		ContinueOnError: ignoreErrors,
	}

	return format.Extract(ctx, src, nil, NewExtract(dst).handler)
}

// Ext returns the extension of the archive.
func (z *zipProvider) Ext() string {
	return ".zip"
}
