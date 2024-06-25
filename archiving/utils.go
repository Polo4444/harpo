package archiving

import (
	"path/filepath"
	"strings"
)

// sanitizePath sanitizes the path by replacing all the slashes with the correct separator
func sanitizePath(p string) string {
	separator := string(filepath.Separator)
	p = strings.ReplaceAll(p, "/", separator)
	p = strings.ReplaceAll(p, "\\", separator)
	return p
}
