package adb

import (
	"github.com/pterm/pterm"
	"testing"
)

func TestDevices(context *testing.T) {
	adb := NewAdbWithLogger("C:/adb/adb.exe", pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace).WithMaxWidth(pterm.GetTerminalWidth()))
	defer adb.Stop()

	devices, err := adb.Devices()

	if err != nil {
		context.Fatal(err)
	}

	for _, device := range devices {
		context.Log(device.GetModel())
	}
}
