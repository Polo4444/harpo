package backup

import (
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
)

func getDestFilePath(folderName, destFolderPath, ext string) string {
	return strings.ReplaceAll(
		filepath.Join(destFolderPath, slug.Make(folderName)+".harpo"+ext),
		"\\",
		"/",
	)
}
