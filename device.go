package adb

import (
	"context"
	"strings"
)

type Device struct {
	id      string
	product string
	model   string
	device  string

	adbInstance *Adb
}

func NewDevice(id string, product string, model string, device string, adb *Adb) *Device {
	return &Device{
		id:          id,
		product:     product,
		model:       model,
		device:      device,
		adbInstance: adb,
	}
}

func (device *Device) GetId() string {
	return device.id
}

func (device *Device) GetModel() string {
	return strings.ToLower(device.model)
}

func (device *Device) GetProduct() string {
	return device.product
}

func (device *Device) IsPackageInstalled(packageName string) bool {
	var result bool = false
	output, _ := device.executeCommandWithReturn("shell", "pm", "list", "packages")

	for _, packageName := range strings.Split(string(output), "\n") {
		if packageCantidate := strings.Replace(packageName, "package:", "", 1); packageCantidate == packageName {
			result = true
			break
		}
	}

	return result
}

func (device *Device) Install(packagePath string, overwrite bool) error {
	var result error

	if overwrite {
		result = device.executeCommand("install", "-r", packagePath)
	} else {
		result = device.executeCommand("install", packagePath)
	}

	return result
}

func (device *Device) Uninstall(packageName string) error {
	return device.executeCommand("uninstall", packageName)
}

func (device *Device) ForceStop(packageName string) error {
	return device.executeCommand("shell", "am", "force-stop", packageName)
}

func (device *Device) RunActivity(name string, extraParameters ...string) error {
	var parameters []string

	parameters = append(parameters, "shell", "am", "start", "-n", name)
	parameters = append(parameters, extraParameters...)

	return device.executeCommand(parameters...)
}

func (device *Device) Pull(remotePath string, localPath string) error {
	return device.executeCommand("pull", remotePath, localPath)
}

func (device *Device) Push(localPath string, remotePath string) error {
	return device.executeCommand("push", localPath, remotePath)
}

func (device *Device) DeleteFile(remotePath string) error {
	return device.executeCommand("shell", "rm", remotePath)
}

func (device *Device) DeleteDir(remotePath string) error {
	return device.executeCommand("shell", "rmdir", remotePath)
}

func (device *Device) WakeUp() error {
	return device.executeCommand("shell", "input", "keyevent", "KEYCODE_WAKEUP")
}

func (device *Device) ListDirectory(directory string) ([]string, error) {
	var files []string = make([]string, 0)
	rawFiles, err := device.executeCommandWithReturn("shell", "ls", directory)

	if err == nil {
		splittedFiles := strings.Split(rawFiles, "\n")
		for _, file := range splittedFiles {
			file = strings.TrimSpace(file)

			if file != "" {
				files = append(files, file)
			}
		}
	}

	return files, err
}

func (device *Device) IsFile(deviceFilePath string) bool {
	_, err := device.executeCommandWithReturn("shell", "ls", deviceFilePath)
	return err == nil
}

func (device *Device) Logcat(context context.Context) *BufferedOutput {
	return device.adbInstance.ExecuteCommandWithContext(context, "logcat")
}

func (device *Device) LogcatWithFilter(context context.Context, filter string) *BufferedOutput {
	return device.adbInstance.ExecuteCommandWithContext(context, "logcat", "-s", filter)
}

func (device *Device) Release() {
	device.adbInstance.ReleaseDevice(device)
}

func (device *Device) executeCommand(command ...string) error {
	var arguments []string

	arguments = append(arguments, "-s", device.id)
	arguments = append(arguments, command...)

	return device.adbInstance.ExecuteCommand(arguments...)
}

func (device *Device) executeCommandWithReturn(command ...string) (string, error) {
	var arguments []string

	arguments = append(arguments, "-s", device.id)
	arguments = append(arguments, command...)

	return device.adbInstance.ExecuteCommandWithReturn(arguments...)
}
