package adb

import (
	"regexp"
	"strings"

	"github.com/GiulianoDecesares/adb/cli"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
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
	adb := &Adb{CommandlineTool: cli.NewCLIWithLogger(binPath, logger)}

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

	if err != nil {
		return "", errors.Wrap(err, "error getting adb version")
	}

	re := regexp.MustCompile(`^Android Debug Bridge version (\S+)`)
	matches := re.FindStringSubmatch(versionString)

	if len(matches) < 2 {
		return "", errors.New("failed to get adb version: invalid output")
	}

	return matches[1], nil
}

func (adb *Adb) Devices() ([]*AdbDevice, error) {
	rawDevices, err := adb.Run("devices", "-l")

	if err != nil {
		return nil, err
	}

	return adb.parseDevicesString(rawDevices), nil
}

func (adb *Adb) parseDevicesString(rawDevices string) []*AdbDevice {
	var devices = make([]*AdbDevice, 0)
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
