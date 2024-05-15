package adb

import (
	"github.com/GiulianoDecesares/adb/cli"
	"github.com/pterm/pterm"
	"strings"
)

const (
	attachedDevicesString = "List of devices attached"
	adbVersionString      = "Android Debug Bridge version "
)

type Adb struct {
	*cli.CommandlineTool
}

func NewAdb(binPath string) *Adb {
	adb := &Adb{
		CommandlineTool: cli.NewCLI(binPath),
	}

	_ = adb.Start()
	return adb
}

func NewAdbWithLogger(binPath string, logger *pterm.Logger) *Adb {
	adb := &Adb{
		CommandlineTool: cli.NewCLIWithLogger(binPath, logger),
	}

	_ = adb.Start()
	return adb
}

func (adb *Adb) Start() error {
	_, err := adb.Run("start-server")
	return err
}

func (adb *Adb) Stop() error {
	_, err := adb.Run("kill-server")
	return err
}

func (adb *Adb) GetVersion() (string, error) {
	versionString, err := adb.Run("version")
	var split = strings.Split(versionString, "\n")

	// Ugly
	if len(split) > 0 {
		versionString = strings.ReplaceAll(split[0], adbVersionString, "")
	}

	return versionString, err
}

func (adb *Adb) Devices() ([]IDevice, error) {
	rawDevices, err := adb.Run("devices", "-l")

	if err != nil {
		return nil, err
	}

	return adb.parseDevicesString(rawDevices), nil
}

func (adb *Adb) parseDevicesString(rawDevices string) []IDevice {
	var devices = make([]IDevice, 0)
	var devicesStrings = strings.Split(rawDevices, "\n")

	if len(devicesStrings) == 0 {
		return devices
	}

	// Trim whitespace
	for index, deviceString := range devicesStrings {
		devicesStrings[index] = strings.TrimSpace(deviceString)
	}

	// Search for "List of devices attached" and skip everything previous
	// since adb could be down and spam some messages while starting
	// and that could mess up all the parsing
	for index, deviceString := range devicesStrings {
		if attachedDevicesString == deviceString {
			devicesStrings = devicesStrings[index+1:]
			break
		}
	}

	for _, deviceString := range devicesStrings {
		if len(deviceString) == 0 {
			continue
		}

		var deviceData []string

		// Parse device data row
		for _, deviceDataItem := range strings.Split(deviceString, " ") {
			if strings.TrimSpace(deviceDataItem) != "" && deviceDataItem != "device" { // Check if not whitespace or "device" adb log
				deviceData = append(deviceData, deviceDataItem)
			}
		}

		if len(deviceData) > 0 {
			devices = append(devices, NewDevice(deviceData[0], adb))
		}
	}

	return devices
}
