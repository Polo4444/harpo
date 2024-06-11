package archiving

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v4"
)

type extract struct {
	dst string
}

func NewExtract(dst string) *extract {
	return &extract{
		dst: dst,
	}
}

func (e *extract) handler(ctx context.Context, f archiver.File) error {

	relativePath := filepath.Join(e.dst, f.NameInArchive)

	if f.FileInfo.IsDir() {
		return os.MkdirAll(relativePath, os.ModePerm)
	}

	// Create the directory for the file if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(relativePath), os.ModePerm); err != nil {
		return err
	}

	// Open the file in the archive
	srcFile, err := f.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(relativePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the contents of the file
	_, err = io.Copy(dstFile, srcFile)
	return err
}
