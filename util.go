package memfs

import "path/filepath"

func clean(path string) string {
	return filepath.Clean(filepath.FromSlash(path))
}
