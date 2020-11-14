// Package bincludehttp wraps a binclude.FileSystem to be compatible with
// the http.FileSystem interface.
// Wrapping the FS may increase binary size significantly (net/http import).
package bincludehttp

import (
	"net/http"
	"os"

	"github.com/lu4p/binclude"
)

type FileSystem struct {
	realFS *binclude.FileSystem
}

// check that the http.FileSystem interface is implemented
var _ http.FileSystem = new(FileSystem)

func Wrap(fs *binclude.FileSystem) FileSystem {
	return FileSystem{realFS: fs}
}

// Open returns a File using the File interface
func (fs *FileSystem) Open(name string) (http.File, error) {
	return fs.realFS.Open(name)
}

// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.
func (fs *FileSystem) Stat(name string) (os.FileInfo, error) {
	return fs.realFS.Stat(name)
}

// ReadFile reads the file named by filename and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.
func (fs *FileSystem) ReadFile(filename string) ([]byte, error) {
	return fs.realFS.ReadFile(filename)
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
func (fs *FileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return fs.realFS.ReadDir(dirname)
}

// CopyFile copies a specific file from a binclude FileSystem to the hosts FileSystem.
// Permissions are copied from the included file.
func (fs *FileSystem) CopyFile(bincludePath, hostPath string) error {
	return fs.realFS.CopyFile(bincludePath, hostPath)
}

// Decompress turns a FileSystem with compressed files into a filesystem without compressed files
func (fs *FileSystem) Decompress() (err error) {
	return fs.realFS.Decompress()
}

// Compress turns a FileSystem without compressed files into a filesystem with compressed files
func (fs *FileSystem) Compress(algo binclude.Compression) error {
	return fs.realFS.Compress(algo)
}
