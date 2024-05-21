package adb

import (
	"context"
	"fmt"
	"github.com/GiulianoDecesares/adb/cli"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type AdbDevice struct {
	serialNo string
	adbCli   *Adb
}

func NewDevice(serialNo string, adb *Adb) *AdbDevice {
	return &AdbDevice{
		serialNo: serialNo,
		adbCli:   adb,
	}
}

func (device *AdbDevice) Reboot(ctx context.Context) error {
	_, err := device.Run("reboot")

	if err != nil {
		return err
	}

	return device.WaitUntilReady(ctx)
}

func (device *AdbDevice) GetSerial() string {
	return device.serialNo
}

func (device *AdbDevice) GetProduct() (string, error) {
	output, err := device.Run("shell", "getprop", "ro.product.device")

	if err != nil {
		return "", errors.Wrap(err, "error getting code name from device")
	}

	return strings.TrimSpace(output), nil
}

func (device *AdbDevice) GetModel() (string, error) {
	output, err := device.Run("shell", "getprop", "ro.product.model")

	if err != nil {
		return "", errors.Wrap(err, "error getting model from device")
	}

	return strings.TrimSpace(output), nil
}

func (device *AdbDevice) GetOsVersion() (string, error) {
	output, err := device.Run("shell", "getprop", "ro.build.version.release")

	if err != nil {
		return "", errors.Wrap(err, "error getting os version from device")
	}

	return strings.TrimSpace(output), nil
}

func (device *AdbDevice) IsPackageInstalled(packageName string) bool {
	var result = false
	output, _ := device.Run("shell", "pm", "list", "packages")

	packageName = strings.TrimSpace(fmt.Sprintf("package:%s", packageName))
	installed := strings.Split(output, "\r\n")

	for _, current := range installed {
		if current == packageName {
			result = true
			break
		}
	}

	return result
}

func (device *AdbDevice) Install(packagePath string, overwrite bool) error {
	var command = make([]string, 0)
	command = append(command, "install")

	if overwrite {
		command = append(command, "-r")
	}

	command = append(command, packagePath)
	_, err := device.Run(command...)
	return err
}

func (device *AdbDevice) Uninstall(packageName string) error {
	_, err := device.Run("uninstall", packageName)
	return err
}

func (device *AdbDevice) ForceStop(packageName string) error {
	_, err := device.Run("shell", "am", "force-stop", packageName)
	return err
}

func (device *AdbDevice) RunActivity(name string, extraParameters ...string) error {
	var parameters []string

	parameters = append(parameters, "shell", "am", "start", "-n", name)
	parameters = append(parameters, extraParameters...)

	_, err := device.Run(parameters...)
	return err
}

func (device *AdbDevice) RunService(name string, extraParameters ...string) error {
	var parameters []string

	parameters = append(parameters, "shell", "am", "startservice", "-n", name)
	parameters = append(parameters, extraParameters...)

	_, err := device.Run(parameters...)
	return err
}

func (device *AdbDevice) Pull(remotePath string, localPath string) error {
	_, err := device.Run("pull", remotePath, localPath)
	return err
}

func (device *AdbDevice) Push(localPath string, remotePath string) error {
	_, err := device.Run("push", localPath, remotePath)
	return err
}

func (device *AdbDevice) DeleteFile(remotePath string) error {
	_, err := device.Run("shell", "rm", remotePath)
	return err
}

func (device *AdbDevice) DeleteDir(remotePath string) error {
	_, err := device.Run("shell", "rmdir", remotePath)
	return err
}

func (device *AdbDevice) CreateDir(remotePath string) error {
	_, err := device.Run("shell", "mkdir", "-p", "\""+remotePath+"\"")
	return err
}

func (device *AdbDevice) WakeUp() error {
	_, err := device.Run("shell", "input", "keyevent", "KEYCODE_WAKEUP")
	return err
}

func (device *AdbDevice) ListDirectory(directory string) ([]string, error) {
	var files = make([]string, 0)
	rawFiles, err := device.Run("shell", "ls", directory)

	if err != nil {
		return nil, err
	}

	splitFiles := strings.Split(rawFiles, "\n")

	for _, file := range splitFiles {
		file = strings.TrimSpace(file)

		if file != "" {
			files = append(files, file)
		}
	}

	return files, err
}

func (device *AdbDevice) IsFile(deviceFilePath string) bool {
	_, err := device.Run("shell", "ls", deviceFilePath)
	return err == nil
}

func (device *AdbDevice) CatFile(deviceFilePath string) (string, error) {
	return device.Run("shell", "cat", deviceFilePath)
}

func (device *AdbDevice) Logcat(context context.Context) *cli.BufferedOutput {
	return device.adbCli.RunWithContext(context, "logcat")
}

func (device *AdbDevice) LogcatWithFilter(context context.Context, filter string) *cli.BufferedOutput {
	return device.adbCli.RunWithContext(context, "logcat", "-s", filter)
}

func (device *AdbDevice) SetPermission(grant bool, packageName string, permission string) error {
	command := "grant"

	if !grant {
		command = "revoke"
	}

	_, err := device.adbCli.Run("shell", "pm", command, packageName, permission)
	return err
}

// This will work only for Android 11.0+
func (device *AdbDevice) SetGps(enabled bool) error {
	enable := "0" // Disabled

	if enabled {
		enable = "3" // Enabled
	}

	_, err := device.adbCli.Run("shell", "settings", "put", "secure", "location_mode", enable)
	return err
}

func (device *AdbDevice) SetRoot(root bool) error {
	command := "root"

	if !root {
		command = "unroot"
	}

	_, err := device.Run(command)
	return err
}

func (device *AdbDevice) Mount(remotePath string) error {
	_, err := device.Run("shell", "service", "call", "mount", "90", "s16", "\""+remotePath+"\"")
	return err
}

func (device *AdbDevice) Chmod(path string, mod string, recursive bool) error {
	var command []string
	command = append(command, "shell", "chmod")

	if recursive {
		command = append(command, "-R")
	}

	command = append(command, mod, path)

	_, err := device.Run(command...)
	return err
}

func (device *AdbDevice) SwitchToFastboot(fastbootCli *Fastboot, ctx context.Context) (*FastbootDevice, error) {
	if fastbootCli == nil || fastbootCli.Check() != nil {
		return nil, errors.New("null or unavailable fastboot CLI")
	}

	if err := device.WaitUntilReady(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to wait for adb device ready")
	}

	_, err := device.Run("reboot", "bootloader")

	if err != nil {
		return nil, errors.Wrap(err, "unable to reboot to fastboot mode")
	}

	fastbootDevice := NewFastbootDevice(device.GetSerial(), fastbootCli)

	if err := fastbootDevice.WaitUntilReady(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to wait for fastboot device ready")
	}

	return fastbootDevice, nil
}

func (device *AdbDevice) Run(command ...string) (string, error) {
	var arguments []string

	arguments = append(arguments, "-s", device.serialNo)
	arguments = append(arguments, command...)

	return device.adbCli.Run(arguments...)
}

func (device *AdbDevice) WaitUntilReady(ctx context.Context) error {
	deviceReady := false

	for !deviceReady {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			time.Sleep(time.Second)

			if err := device.WakeUp(); err == nil {
				deviceReady = true
			}
		}
	}

	return nil
}
