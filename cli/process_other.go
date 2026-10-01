//go:build !windows

package cli

import "os/exec"

// configureCommand is a no-op outside Windows.
func configureCommand(cmd *exec.Cmd) {}
