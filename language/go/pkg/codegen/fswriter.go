package codegen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
)

// A FileSystemWriter provides the means to write files into a collection of
// named files.
type FileSystemWriter interface {
	Create(name string) (io.WriteCloser, error)
}

// OSFileSystemWriter allows writing into the OS filesystem.
type OSFileSystemWriter struct {
	Root string
}

// OpenWriter opens an OS file for writing.
func (w *OSFileSystemWriter) Create(name string) (io.WriteCloser, error) {
	p := path.Join(w.Root, name)
	d := path.Dir(p)
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, errors.WithStack(err)
	}
	return os.Create(p)
}

// MemoryFileSystemWriter allows writing into an in-memory map.
type MemoryFileSystemWriter struct {
	files map[string]*closableBuffer
}

// OpenWriter opens an OS file for writing.
func (w *MemoryFileSystemWriter) OpenWriter(name string) (io.WriteCloser, error) {
	var buf closableBuffer
	w.files[name] = &buf
	return &buf, nil
}

// Get returns the bytes of the named file.
func (w *MemoryFileSystemWriter) Get(name string) []byte {
	return w.files[name].Bytes()
}

type closableBuffer struct {
	bytes.Buffer
	closed bool
}

func (buf *closableBuffer) Write(b []byte) (n int, err error) {
	if buf.closed {
		return 0, errors.WithStack(fmt.Errorf("writing to closed file"))
	}
	return buf.Buffer.Write(b)
}

func (*closableBuffer) Close() error {
	return nil
}
