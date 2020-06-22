package binclude

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Debug if set to true files are read via os.Open() and the bincluded files are
// ignored, use when developing.
var Debug = false

// Include this file/ directory (including subdirectories) relative to the package path (noop)
// The path is walked via filepath.Walk and all files found are included
func Include(name string) {}

// FileSystem implements access to a collection of named files.
type FileSystem map[string]*File

// check that the http.FileSystem interface is implemented
var _ http.FileSystem = new(FileSystem)

// Open returns a File using the http.File interface
func (fs FileSystem) Open(name string) (http.File, error) {
	if Debug {
		name = filepath.FromSlash(name)

		return os.Open(name)
	}

	name = strings.TrimPrefix(name, "./")
	if f, ok := fs[name]; ok {
		f.reader = bytes.NewReader(f.Content)
		f.path = name
		f.fs = &fs
		return f, nil
	}

	return nil, &os.PathError{"open", name, errors.New("File does not exist in binclude map")}
}

// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.
func (fs FileSystem) Stat(name string) (os.FileInfo, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Stat()
}

// ReadFile reads the file named by filename and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.
func (fs FileSystem) ReadFile(filename string) ([]byte, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
func (fs FileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, _ := f.Readdir(-1)
	f.Close()
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

// CopyFile copies a specific file from a binclude FileSystem to the hosts FileSystem.
// Permissions are copied from the included file.
func (fs FileSystem) CopyFile(bincludePath, hostPath string) error {
	src, err := fs.Open(bincludePath)
	if err != nil {
		return err
	}
	defer src.Close()

	info, _ := src.Stat()

	log.Println("Copy Path:", hostPath, "binclude path:", bincludePath)

	dst, err := os.OpenFile(hostPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	info, err = os.Stat(hostPath)
	if err != nil {
	return err
}

	log.Println("Info after Copy:", info.IsDir(), info.Name(), info.Mode(), info.Size())

	return nil
}
// File implements the io.Reader, io.Seeker, io.Closer and http.File interfaces
type File struct {
	Filename string
	Mode     os.FileMode
	ModTime  time.Time
	Content  []byte
	reader   *bytes.Reader
	path     string
	fs       *FileSystem
}

// check that the http.File interface is implemented
var _ http.File = new(File)

// Read implements the io.Reader interface.
func (f *File) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}

// Name returns the name of the file as presented to Open.
func (f *File) Name() string {
	return f.path
}

// Close closes the File, rendering it unusable for I/O.
func (f *File) Close() error {
	f.reader = nil
	return nil
}

// Size returns the original length of the underlying byte slice.
// Size is the number of bytes available for reading via ReadAt.
// The returned value is always the same and is not affected by calls
// to any other method.
func (f *File) Size() int64 {
	return int64(len(f.Content))
}

// Readdir reads the contents of the directory associated with file and
// returns a slice of up to n FileInfo values, as would be returned
// by Lstat, in directory order. Subsequent calls on the same file will yield
// further FileInfos.
func (f *File) Readdir(count int) (infos []os.FileInfo, err error) {
	fileDir := f.Name()
	if !f.Mode.IsDir() {
		fileDir = filepath.Dir(f.path)
	}

	for path, file := range *f.fs {
		if filepath.Dir(path) != fileDir {
			continue
		}

		info, _ := file.Stat()

		infos = append(infos, info)
	}

	return infos, nil
}

// Stat returns the FileInfo structure describing file.
// Error is always nil
func (f *File) Stat() (os.FileInfo, error) {
	return &FileInfo{
		name:    f.Filename,
		mode:    f.Mode,
		size:    f.Size(),
		modtime: f.ModTime,
	}, nil
}

// Seek implements the io.Seeker interface.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.reader.Seek(offset, whence)
}

// FileInfo implements the os.FileInfo interface.
type FileInfo struct {
	name    string
	mode    os.FileMode
	modtime time.Time
	size    int64
}

// check that the os.FileInfo interface is implemented
var _ os.FileInfo = new(FileInfo)

// Name returns the base name of the file
func (info *FileInfo) Name() string {
	return info.name
}

// Size returns the length in bytes
func (info *FileInfo) Size() int64 {
	return info.size
}

// Mode returns the file mode bits
func (info *FileInfo) Mode() os.FileMode {
	return info.mode
}

// ModTime returns the modification time (returns current time)
func (info *FileInfo) ModTime() time.Time {
	return info.modtime
}

// IsDir abbreviation for Mode().IsDir()
func (info *FileInfo) IsDir() bool {
	return info.Mode().IsDir()
}

// Sys underlying data source (returns nil)
func (info *FileInfo) Sys() interface{} {
	return nil
}
