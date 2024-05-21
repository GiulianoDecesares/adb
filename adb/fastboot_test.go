package adb

import (
	"context"
	"testing"
	"time"
)

func TestNewDevice(context *testing.T) {
	fastbootCli := NewFastboot("C:/adb/fastboot.exe")
	devices, err := fastbootCli.Devices()

	if err != nil {
		context.Fatal(err)
	}

	for _, device := range devices {
		productName, err := device.GetProduct()

		if err != nil {
			context.Fatal(err)
		}

		context.Log(productName)
	}
}

func TestSwitchToAdb(testingContext *testing.T) {
	fastbootCli := NewFastboot("C:/adb/fastboot.exe")
	devices, err := fastbootCli.Devices()

	if err != nil {
		testingContext.Fatal(err)
	}

	if len(devices) == 0 {
		testingContext.Fatal("no devices")
	}

	device := devices[0]
	adbCli := NewAdb("C:/adb/adb.exe")
	defer adbCli.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	adbDevice, err := device.SwitchToAdb(adbCli, ctx)

	if err != nil {
		testingContext.Fatal(err)
	}

	err = adbDevice.WaitUntilReady(ctx)

	if err != nil {
		testingContext.Fatal(err)
	}
}
