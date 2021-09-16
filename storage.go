package memfs

import (
	"os"

	"github.com/go-git/go-billy/v5"
)

type Storage interface {
	Has(path string) bool
	Get(path string) (billy.File, bool)
	New(path string, mode os.FileMode, flag int) (billy.File, error)
	Children(path string) []billy.File
	Rename(from, to string) error
	Remove(path string) error
}
