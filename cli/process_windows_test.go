package cli

import (
	"context"
	"os/exec"
	"syscall"
	"testing"

	"golang.org/x/sys/windows"
)

func TestBuildCommandHidesConsoleWindow(t *testing.T) {
	tool := NewCLI(`tools\adb.exe`)
	cmd := tool.buildCommand("devices", "-l")

	assertNoConsole(t, cmd)
	if got, want := cmd.Args, []string{`tools\adb.exe`, "devices", "-l"}; !sameArgs(got, want) {
		t.Fatalf("Args = %q, want %q", got, want)
	}
}

func TestBuildCommandWithContextHidesConsoleWindow(t *testing.T) {
	tool := NewCLI(`tools\fastboot.exe`)
	cmd := tool.buildCommandWithContext(context.Background(), "devices")

	assertNoConsole(t, cmd)
	if got, want := cmd.Args, []string{`tools\fastboot.exe`, "devices"}; !sameArgs(got, want) {
		t.Fatalf("Args = %q, want %q", got, want)
	}
}

func TestConfigureCommandPreservesExistingAttributes(t *testing.T) {
	cmd := &exec.Cmd{
		SysProcAttr: &syscall.SysProcAttr{
			CreationFlags:    windows.CREATE_NEW_PROCESS_GROUP,
			NoInheritHandles: true,
		},
	}

	configureCommand(cmd)

	if !cmd.SysProcAttr.NoInheritHandles {
		t.Fatal("NoInheritHandles was cleared")
	}
	if cmd.SysProcAttr.CreationFlags&windows.CREATE_NEW_PROCESS_GROUP == 0 {
		t.Fatal("existing CreationFlags were cleared")
	}
	assertNoConsole(t, cmd)
}

func assertNoConsole(t *testing.T, cmd *exec.Cmd) {
	t.Helper()

	if cmd.SysProcAttr == nil {
		t.Fatal("SysProcAttr was not configured")
	}
	if !cmd.SysProcAttr.HideWindow {
		t.Fatal("HideWindow was not set")
	}
	if cmd.SysProcAttr.CreationFlags&windows.CREATE_NO_WINDOW == 0 {
		t.Fatal("CREATE_NO_WINDOW was not set")
	}
}

func sameArgs(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
