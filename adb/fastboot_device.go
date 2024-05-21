package adb

import (
	"context"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type FastbootDevice struct {
	serialNo    string
	fastbootCli *Fastboot
}

func NewFastbootDevice(id string, fastboot *Fastboot) *FastbootDevice {
	return &FastbootDevice{
		serialNo:    id,
		fastbootCli: fastboot,
	}
}

func (device *FastbootDevice) GetSerial() string { return device.serialNo }

func (device *FastbootDevice) GetProduct() (string, error) {
	output, err := device.Run("getvar", "product")

	if err != nil {
		return "", err
	}

	outputData := strings.Split(output, "\n")
	productName := ""

	for _, line := range outputData {
		if strings.Contains(line, "product") {
			lineData := strings.Split(line, ":")

			if len(lineData) != 2 {
				return "", errors.New("failed to parse product output")
			}

			productName = strings.TrimSpace(lineData[1])
		}
	}

	if len(productName) == 0 {
		return "", errors.New("failed to parse product output")
	}

	return productName, nil
}

func (device *FastbootDevice) Reboot(ctx context.Context) error {
	_, err := device.Run("reboot")

	if err != nil {
		return err
	}

	return device.WaitUntilReady(ctx)
}

func (device *FastbootDevice) Run(rawCommand ...string) (string, error) {
	var arguments []string

	arguments = append(arguments, "-s", device.serialNo)
	arguments = append(arguments, rawCommand...)

	return device.fastbootCli.Run(arguments...)
}

func (device *FastbootDevice) SwitchToAdb(adbCli *Adb, ctx context.Context) (*AdbDevice, error) {
	if adbCli == nil || adbCli.Check() != nil {
		return nil, errors.New("null or unavailable fastboot CLI")
	}

	if err := device.WaitUntilReady(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to wait for fastboot device ready")
	}

	_, err := device.Run("continue")

	if err != nil {
		return nil, errors.Wrap(err, "unable to reboot to adb mode")
	}

	adbDevice := NewDevice(device.GetSerial(), adbCli)

	if err := adbDevice.WaitUntilReady(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to wait for adb device ready")
	}

	return adbDevice, nil
}

func (device *FastbootDevice) WaitUntilReady(ctx context.Context) error {
	deviceResponsive := false

	for !deviceResponsive {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			time.Sleep(time.Second)

			if _, err := device.GetProduct(); err == nil {
				deviceResponsive = true
			}
		}
	}

	return nil
}
