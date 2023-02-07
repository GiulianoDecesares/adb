package adb

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.globant.com/Facebook/hwtp-oculus-launcher/lib/utils"
)

const (
	attachedDevicesString = "List of devices attached"
	adbVersionString      = "Android Debug Bridge version "
)

type AdbConfig struct {
	AdbPath string `yaml:"Path"`
}

type Adb struct {
	config AdbConfig
}

func New(config AdbConfig) *Adb {
	return &Adb{
		config: config,
	}
}

func (adb *Adb) Check() error {
	var result error = nil
	var isFile bool = false

	if isFile, result = utils.IsFile(adb.config.AdbPath); result == nil {
		if isFile {
			if !utils.FileExists(adb.config.AdbPath) {
				result = fmt.Errorf("File %s doesn't exists", adb.config.AdbPath)
			}
		} else {
			result = fmt.Errorf("%s is not ADB executable file", adb.config.AdbPath)
		}
	}

	return result
}

// Start starts the ADB server
func (adb *Adb) Start() error {
	return adb.ExecuteCommand("start-server")
}

// Stop kills the ADB server
func (adb *Adb) Stop() error {
	return adb.ExecuteCommand("kill-server")
}

func (adb *Adb) GetVersion() (string, error) {
	versionString, err := adb.ExecuteCommandWithReturn("version")
	var splitted []string = strings.Split(versionString, "\n")

	// Ugly
	if len(splitted) > 0 {
		versionString = strings.ReplaceAll(splitted[0], adbVersionString, "")
	}

	return versionString, err
}

func (adb *Adb) Devices() ([]Device, error) {
	var result error
	var rawDevices string
	var devices []Device

	if rawDevices, result = adb.ExecuteCommandWithReturn("devices", "-l"); result == nil {
		devices = adb.parseDevicesString(rawDevices)
	}

	return devices, result
}

func (adb *Adb) ExecuteCommand(command ...string) error {
	output, err := adb.ExecuteCommandWithReturn(command...)

	if err != nil {
		err = fmt.Errorf("%s -> %s", err.Error(), output)
	}

	return err
}

func (adb *Adb) ExecuteCommandWithContext(context context.Context, command ...string) *BufferedOutput {
	var result BufferedOutput

	executableCommand := exec.CommandContext(context, adb.config.AdbPath, command...)
	executableCommand.Stdout = &result.Out
	executableCommand.Stderr = &result.Err

	result.Error = executableCommand.Start()

	return &result
}

func (adb *Adb) ExecuteCommandWithReturn(command ...string) (string, error) {
	rawOutput, result := exec.Command(adb.config.AdbPath, command...).CombinedOutput()
	return string(rawOutput), result
}

func (adb *Adb) parseDevicesString(rawDevices string) []Device {
	var devices []Device = make([]Device, 0)
	var devicesStrings []string = strings.Split(rawDevices, "\n")

	if len(devicesStrings) > 0 {
		// Trim all whitespace
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
			if len(deviceString) > 0 {
				var deviceData []string

				// Parse device data row
				for _, deviceDataItem := range strings.Split(deviceString, " ") {
					if strings.TrimSpace(deviceDataItem) != "" && deviceDataItem != "device" { // Check if not whitespace or "device" adb log
						deviceData = append(deviceData, deviceDataItem)
					}
				}

				if len(deviceData) > 0 {
					var deviceName string = deviceData[0]
					var deviceProduct string
					var deviceModel string
					var deviceStr string

					// Get some more information
					if len(deviceData) > 1 {
						deviceData = deviceData[1:] // Skip device name

						for _, data := range deviceData {
							var dataItems []string = strings.Split(data, ":")

							if len(dataItems) > 1 {
								var dataType string = dataItems[0]
								var data string = dataItems[1]

								switch dataType {
								case "product":
									deviceProduct = data

								case "model":
									deviceModel = data

								case "device":
									deviceStr = data
								}
							}
						}
					}

					var device *Device = NewDevice(deviceName, deviceProduct, deviceModel, deviceStr, adb)
					devices = append(devices, *device)
				}
			}
		}
	}

	return devices
}
