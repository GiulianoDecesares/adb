package adb

import (
	"context"
	"github.com/pterm/pterm"
	"testing"
	"time"
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

func TestSwitchToFastboot(testingContext *testing.T) {
	adbCli := NewAdb("C:/adb/adb.exe")
	defer adbCli.Stop()

	devices, err := adbCli.Devices()

	if err != nil {
		testingContext.Fatal(err)
	}

	if len(devices) == 0 {
		testingContext.Fatal("no devices")
	}

	device := devices[0]
	fastbootCli := NewFastboot("C:/adb/fastboot.exe")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	fastbootDevice, err := device.SwitchToFastboot(fastbootCli, ctx)

	if err != nil {
		testingContext.Fatal(err)
	}

	err = fastbootDevice.WaitUntilReady(ctx)

	if err != nil {
		testingContext.Fatal(err)
	}
}
