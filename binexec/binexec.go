// Package binexec implements a wrapper for os/exec
package binexec

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lu4p/binclude"
)

// Cmd same as Cmd in the os/exec package
type Cmd struct {
	OsCmd *exec.Cmd
	// Cache if set to true the binary won't be deleted after execution.
	// If the ModTime of the cached binclude file changes the cache gets invalidated automtically.
	Cache bool
}

// Command similar to Command in the os/exec package,
// but copies the executeable to run from bincludePath
// to the host os.
func Command(fs *binclude.FileSystem, bincludePath string, arg ...string) (*Cmd, error) {
	execPath, err := copyCommand(fs, bincludePath)
	if err != nil {
		return nil, err
	}

	cmd := Cmd{
		OsCmd: exec.Command(execPath, arg...),
	}

	return &cmd, nil
}

// copyCommand copy a file from binclude.FileSystem to os.UserCacheDir()
func copyCommand(fs *binclude.FileSystem, bincludePath string) (string, error) {
	dir, _ := os.UserCacheDir()

	info, err := fs.Stat(bincludePath)
	if err != nil {
		return "", err
	}

	nanoSec := strconv.Itoa(info.ModTime().Nanosecond())

	namePart := "_" + filepath.Base(bincludePath)

	execPath := filepath.Join(dir, nanoSec+namePart)

	// exit early if file is already cached
	_, err = os.Stat(execPath)
	if err == nil {
		return execPath, nil
	}

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// remove invalidated cache files
	for _, info := range infos {
		if strings.HasSuffix(info.Name(), namePart) {
			os.Remove(filepath.Join(dir, info.Name())) // don't check for error because we don't really care if the file is removed
			// no break because there could be multiple cached versions
		}
	}

	return execPath, fs.CopyFile(bincludePath, execPath)
}

// CommandContext similar to CommandContext in the os/exec
// package but copies the executeable to run from bincludePath
// to the host os.
func CommandContext(ctx context.Context, fs *binclude.FileSystem, bincludePath string, arg ...string) (*Cmd, error) {
	execPath, err := copyCommand(fs, bincludePath)
	if err != nil {
		return nil, err
	}

	cmd := Cmd{
		OsCmd: exec.CommandContext(ctx, execPath, arg...),
	}

	return &cmd, nil
}

// Run is similar to (*Cmd).Run() in the os/exec package,
// but deletes the executable at Cmd.Path if c.Cache is false
func (c *Cmd) Run() error {
	if !c.Cache {
		defer os.Remove(c.OsCmd.Path)
	}

	return c.OsCmd.Run()
}

// Start same as (*Cmd).Start() in the os/exec package
func (c *Cmd) Start() error {
	return c.OsCmd.Start()
}

// StderrPipe same as (*Cmd).StderrPipe() in the os/exec package
func (c *Cmd) StderrPipe() (io.ReadCloser, error) {
	return c.OsCmd.StderrPipe()
}

// StdinPipe same as (*Cmd).StdinPipe() in the os/exec package
func (c *Cmd) StdinPipe() (io.WriteCloser, error) {
	return c.OsCmd.StdinPipe()
}

// StdoutPipe same as (*Cmd).StdoutPipe() in the os/exec package
func (c *Cmd) StdoutPipe() (io.ReadCloser, error) {
	return c.OsCmd.StdoutPipe()
}

// String same as (*Cmd).String() in the os/exec package
func (c *Cmd) String() string {
	return c.OsCmd.String()
}

// Wait is similar to (*Cmd).Wait() in the os/exec package,
// but deletes the executable at Cmd.Path if c.Cache is false
func (c *Cmd) Wait() error {
	if !c.Cache {
		defer os.Remove(c.OsCmd.Path)
	}

	return c.OsCmd.Wait()
}
