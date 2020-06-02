package mocka

import (
	"os"
	"time"
)

// variables used for unit testing
var _timeNow = time.Now

// mockFileInfo describes meta data about a mock file.
type mockFileInfo struct {
	name string
	size int64
}

// Name is an implemented mock function for os.FileInfo
func (fileInfo mockFileInfo) Name() string {
	return fileInfo.name
}

// Size is an implemented mock function for os.FileInfo
func (fileInfo mockFileInfo) Size() int64 {
	return fileInfo.size
}

// Mode is an implemented mock function for os.FileInfo
func (mockFileInfo) Mode() os.FileMode {
	return os.FileMode(0)
}

// ModTime is an implemented mock function for os.FileInfo
func (mockFileInfo) ModTime() time.Time {
	return _timeNow()
}

// IsDir is an implemented mock function for os.FileInfo
func (mockFileInfo) IsDir() bool {
	return false
}

// Sys is an implemented mock function for os.FileInfo
func (mockFileInfo) Sys() interface{} {
	return nil
}
