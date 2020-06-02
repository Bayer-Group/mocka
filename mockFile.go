package mocka

import (
	"fmt"
	"io"
	"os"
)

// mockFile represents a mock memory only file
type mockFile struct {
	name   string
	buf    []byte
	offset int64
	base   int64
	limit  int64
}

// Close is an implemented mock function for io.Closer
func (*mockFile) Close() error {
	return nil
}

// Read is an implemented mock function for io.Reader
func (f *mockFile) Read(p []byte) (n int, err error) {
	n, err = f.ReadAt(p, f.offset)
	f.offset += int64(n)

	return n, err
}

// Seek is an implemented mock function for io.Seeker
func (f *mockFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	default:
		return 0, fmt.Errorf("mocka: invalid whence received - %v", whence)
	case io.SeekStart:
		offset += f.base
	case io.SeekCurrent:
		offset += f.offset
	case io.SeekEnd:
		offset += f.limit
	}

	if offset < f.base {
		return 0, fmt.Errorf("mocka: desired offset %v is below the offset base of %v", offset, f.base)
	} else if offset > f.limit {
		return 0, fmt.Errorf("mocka: desired offset %v is above the offset limit of %v", offset, f.limit)
	}

	f.offset = offset
	return f.offset, nil
}

// Stat is an implemented mock function for os.File
func (f *mockFile) Stat() (os.FileInfo, error) {
	return mockFileInfo{name: f.name, size: f.limit}, nil
}

// ReadAt is an implemented mock function for io.ReaderAt
func (f *mockFile) ReadAt(p []byte, startingPos int64) (n int, err error) {
	bytesToRead := len(p) + int(startingPos)
	remainingBytes := len(f.buf) - int(startingPos)

	if bytesToRead > remainingBytes {
		bytesToRead = remainingBytes + int(startingPos)
	}

	if n = copy(p, f.buf[startingPos:bytesToRead]); n == 0 {
		return n, io.EOF
	}

	return n, nil
}

// ReadAt is an implemented mock function for io.Writer
func (f *mockFile) Write(b []byte) (int, error) {
	bytesToWrite := int64(len(b))
	n := 0

	for off := f.offset; off < bytesToWrite && off < f.limit; off++ {
		f.buf[off] = b[n]
		n++
	}

	if n != len(b) {
		return n, io.ErrShortWrite
	}

	return n, nil
}

// Truncate is an implemented mock function for os.File
func (f *mockFile) Truncate(size int64) error {
	if size >= f.limit || size < 0 {
		return nil
	}

	f.buf = f.buf[:size]
	f.limit = int64(len(f.buf))

	return nil
}
