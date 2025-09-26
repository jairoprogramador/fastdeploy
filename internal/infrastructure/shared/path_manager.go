package shared

import (
	"path/filepath"
)

func GetPath(paths ...string) string {
	return filepath.Join(paths...)
}
