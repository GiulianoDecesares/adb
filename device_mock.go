package adb

import (
	"context"
	"strings"
)

type MockDevice struct {
	id      string
	product string
	model   string
	device  string

	adbInstance *Adb
}

func NewMockDevice() *Device {
	return &Device{
		id:          "mock-id",
		product:     "mock-product",
		model:       "mock-model",
		device:      "mock-device",
		adbInstance: nil,
	}
}

func (device *MockDevice) GetId() string {
	return device.id
}

func (device *MockDevice) GetModel() string {
	return strings.ToLower(device.model)
}

func (device *MockDevice) GetProduct() string {
	return device.product
}

func (device *MockDevice) IsPackageInstalled(packageName string) bool {
	return true
}

func (device *MockDevice) Install(packagePath string, overwrite bool) error {
	return nil
}

func (device *MockDevice) Uninstall(packageName string) error {
	return nil
}

func (device *MockDevice) ForceStop(packageName string) error {
	return nil
}

func (device *MockDevice) RunActivity(name string, extraParameters ...string) error {
	return nil
}

func (device *MockDevice) Pull(remotePath string, localPath string) error {
	return nil
}

func (device *MockDevice) Push(localPath string, remotePath string) error {
	return nil
}

func (device *MockDevice) DeleteFile(remotePath string) error {
	return nil
}

func (device *MockDevice) DeleteDir(remotePath string) error {
	return nil
}

func (device *MockDevice) WakeUp() error {
	return nil
}

func (device *MockDevice) ListDirectory(directory string) ([]string, error) {
	return nil, nil
}

func (device *MockDevice) IsFile(deviceFilePath string) bool {
	return true
}

func (device *MockDevice) Logcat(context context.Context) *BufferedOutput {
	return nil // NOTE :: Careful here
}

func (device *MockDevice) LogcatWithFilter(context context.Context, filter string) *BufferedOutput {
	return nil // NOTE :: Careful here
}

func (device *MockDevice) Release() {}

func (device *MockDevice) executeCommand(command ...string) error {
	return nil
}

func (device *MockDevice) executeCommandWithReturn(command ...string) (string, error) {
	return "", nil
}
