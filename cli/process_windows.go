package cli

import (
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

// configureCommand keeps non-interactive child processes from allocating a
// visible console window. Existing process attributes are preserved.
func configureCommand(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	cmd.SysProcAttr.HideWindow = true
	cmd.SysProcAttr.CreationFlags |= windows.CREATE_NO_WINDOW
}
