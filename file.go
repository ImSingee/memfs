package memfs

import (
	"errors"
	"io"
	"os"

	"github.com/go-git/go-billy/v5"
)

type File struct {
	name     string
	content  *content
	position int64
	flag     int
	mode     os.FileMode

	isClosed bool
}

var _ billy.File = (*File)(nil)

func (f *File) Name() string {
	return f.name
}

func (f *File) Read(b []byte) (int, error) {
	n, err := f.ReadAt(b, f.position)
	f.position += int64(n)

	if err == io.EOF && n != 0 {
		err = nil
	}

	return n, err
}

func (f *File) ReadAt(b []byte, off int64) (int, error) {
	if f.isClosed {
		return 0, os.ErrClosed
	}

	if !isReadAndWrite(f.flag) && !isReadOnly(f.flag) {
		return 0, errors.New("read not supported")
	}

	n, err := f.content.ReadAt(b, off)

	return n, err
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.isClosed {
		return 0, os.ErrClosed
	}

	switch whence {
	case io.SeekCurrent:
		f.position += offset
	case io.SeekStart:
		f.position = offset
	case io.SeekEnd:
		f.position = int64(f.content.Len()) + offset
	}

	return f.position, nil
}

func (f *File) Write(p []byte) (int, error) {
	if f.isClosed {
		return 0, os.ErrClosed
	}

	if !isReadAndWrite(f.flag) && !isWriteOnly(f.flag) {
		return 0, errors.New("write not supported")
	}

	n, err := f.content.WriteAt(p, f.position)
	f.position += int64(n)

	return n, err
}

func (f *File) Close() error {
	if f.isClosed {
		return os.ErrClosed
	}

	f.isClosed = true
	return nil
}

func (f *File) Truncate(size int64) error {
	if size < int64(len(f.content.bytes)) {
		f.content.bytes = f.content.bytes[:size]
	} else if more := int(size) - len(f.content.bytes); more > 0 {
		f.content.bytes = append(f.content.bytes, make([]byte, more)...)
	}

	return nil
}

func (f *File) Duplicate(filename string, mode os.FileMode, flag int) billy.File {
	newFile := &File{
		name:    filename,
		content: f.content,
		mode:    mode,
		flag:    flag,
	}

	if isAppend(flag) {
		newFile.position = int64(newFile.content.Len())
	}

	if isTruncate(flag) {
		newFile.content.Truncate()
	}

	return newFile
}

func (f *File) Stat() (os.FileInfo, error) {
	return &FileInfo{
		name: f.Name(),
		mode: f.mode,
		size: f.content.Len(),
	}, nil
}

// Lock is a no-op in memfs.
func (f *File) Lock() error {
	return nil
}

// Unlock is a no-op in memfs.
func (f *File) Unlock() error {
	return nil
}
