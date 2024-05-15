package adb

import (
	"context"
)

type IAdb interface {
	Check() error
	Start() error
	Stop() error
	GetVersion() (string, error)
	Devices() ([]IDevice, error)
	ReleaseDevice(device *Device)
	ExecuteCommand(command ...string) error
	ExecuteCommandWithContext(ctx context.Context, command ...string) *BufferedOutput
	ExecuteCommandWithReturn(command ...string) (string, error)
}
