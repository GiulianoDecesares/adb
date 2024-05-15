package fastboot

import (
	"github.com/GiulianoDecesares/adb/cli"
	"github.com/pterm/pterm"
	"strings"
)

type Fastboot struct {
	*cli.CommandlineTool
}

func NewFastboot(binPath string) *Fastboot {
	return &Fastboot{
		CommandlineTool: cli.NewCLI(binPath),
	}
}

func NewFastbootWithLogger(binPath string, logger *pterm.Logger) *Fastboot {
	return &Fastboot{
		CommandlineTool: cli.NewCLIWithLogger(binPath, logger),
	}
}

func (fastboot *Fastboot) Devices() ([]*Device, error) {
	var devices = make([]*Device, 0)
	output, err := fastboot.Run("devices")

	if err != nil {
		return devices, err
	}

	devicesStrings := strings.Split(output, "\n")

	if len(devicesStrings) == 0 {
		return devices, nil
	}

	// Trim whitespace
	for index, deviceString := range devicesStrings {
		devicesStrings[index] = strings.TrimSpace(deviceString)
	}

	// Get device serialNo
	for _, deviceString := range devicesStrings {
		if len(deviceString) == 0 {
			continue
		}

		deviceData := strings.Split(deviceString, "\t")
		if len(deviceData) == 0 {
			continue
		}

		deviceId := strings.TrimSpace(deviceData[0])

		if len(deviceId) == 0 {
			continue
		}

		devices = append(devices, NewDevice(deviceId, fastboot))
	}

	return devices, nil
}
