//go:build !windows

package cli

import (
	"context"
	"testing"
)

func TestBuildCommandLeavesProcessAttributesUnset(t *testing.T) {
	tool := NewCLI("tools/adb")
	cmd := tool.buildCommand("devices", "-l")

	if cmd.SysProcAttr != nil {
		t.Fatalf("SysProcAttr = %#v, want nil", cmd.SysProcAttr)
	}
	if got, want := cmd.Args, []string{"tools/adb", "devices", "-l"}; !sameArgs(got, want) {
		t.Fatalf("Args = %q, want %q", got, want)
	}
}

func TestBuildCommandWithContextLeavesProcessAttributesUnset(t *testing.T) {
	tool := NewCLI("tools/fastboot")
	cmd := tool.buildCommandWithContext(context.Background(), "devices")

	if cmd.SysProcAttr != nil {
		t.Fatalf("SysProcAttr = %#v, want nil", cmd.SysProcAttr)
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
