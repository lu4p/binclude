// Package binexec implements a wrapper for os/exec
package binexec

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lu4p/binclude"
)

// Cmd same as Cmd in the os/exec package
type Cmd struct {
	OsCmd *exec.Cmd
}

// Command similar to Command in the os/exec package,
// but copies the executeable to run from bincludePath
// to the host os.
func Command(fs binclude.FileSystem, bincludePath string, arg ...string) (*Cmd, error) {
	bincludePath = filepath.FromSlash(bincludePath)
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
func copyCommand(fs binclude.FileSystem, bincludePath string) (string, error) {
	dir, _ := os.UserCacheDir()

	execPath := filepath.Join(dir, filepath.Base(bincludePath))
	return execPath, fs.CopyFile(bincludePath, execPath)
}

// CommandContext similar to CommandContext in the os/exec
// package but copies the executeable to run from bincludePath
// to the host os.
func CommandContext(ctx context.Context, fs binclude.FileSystem, bincludePath string, arg ...string) (*Cmd, error) {
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
// but deletes the executable at Cmd.Path
func (c *Cmd) Run() error {
	defer os.Remove(c.OsCmd.Path)
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
// but deletes the executable at Cmd.Path
func (c *Cmd) Wait() error {
	defer os.Remove(c.OsCmd.Path)
	return c.OsCmd.Wait()
}
