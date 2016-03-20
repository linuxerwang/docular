package dar

import (
	"io"
	"os"
)

type Reader struct {
	r    io.ReaderAt
	File []*File
}

type ReadCloser struct {
	f *os.File
	Reader
}

type File struct {
	FileHeader
	zipr         io.ReaderAt
	zipsize      int64
	headerOffset int64
}

func OpenReader(path string) (*ReadCloser, error) {
	return nil, nil
}

// Close closes the dar file, rendering it unusable for I/O.
func (rc *ReadCloser) Close() error {
	return rc.f.Close()
}

// Open returns a ReadCloser that provides access to the File's contents.
// Multiple files may be read concurrently.
func (f *File) Open() (rc io.ReadCloser, err error) {
	return nil, nil
}
