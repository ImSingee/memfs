package memfs

import (
	"os"
)

type Storage interface {
	Has(path string) bool
	Get(path string) (*File, bool)
	New(path string, mode os.FileMode, flag int) (*File, error)
	Children(path string) []*File
	Rename(from, to string) error
	Remove(path string) error
}
