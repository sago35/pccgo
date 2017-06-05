package exec

import (
	"fmt"
	"io"
	"os/exec"
	"syscall"
	"time"
)

type Cmd struct {
	exec.Cmd
	Stdout io.Writer `json:"-"`
	Stderr io.Writer `json:"-"`
	Join   bool
	Target string
}

func (c Cmd) Run() int {
	e := exec.Command(c.Path, c.Args...)
	e.Env = c.Env
	e.Dir = c.Dir
	e.Stdout = c.Stdout
	e.Stderr = c.Stderr

	var exitStatus int = 0
	if err := e.Run(); err != nil {
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				exitStatus = s.ExitStatus()
			} else {
				exitStatus = 2
			}
		}
	} else {
		exitStatus = 0
	}

	return exitStatus
}
