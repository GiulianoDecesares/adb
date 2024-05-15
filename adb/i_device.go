package adb

import (
	"context"
	"github.com/GiulianoDecesares/adb/cli"
	"github.com/GiulianoDecesares/adb/fastboot"
)

type IDevice interface {
	GetSerial() string

	GetProduct() (string, error)
	GetModel() (string, error)

	GetOsVersion() (string, error)

	IsPackageInstalled(packageName string) bool
	Install(packagePath string, overwrite bool) error
	Uninstall(packageName string) error

	ForceStop(packageName string) error
	RunActivity(name string, extraParameters ...string) error

	Pull(remotePath string, localPath string) error
	Push(localPath string, remotePath string) error

	DeleteFile(remotePath string) error
	DeleteDir(remotePath string) error
	CreateDir(remotePath string) error

	WakeUp() error

	ListDirectory(directory string) ([]string, error)
	IsFile(deviceFilePath string) bool

	Logcat(context context.Context) *cli.BufferedOutput
	LogcatWithFilter(context context.Context, filter string) *cli.BufferedOutput

	SetPermission(grant bool, packageName string, permission string) error
	SetGps(enabled bool) error

	SetRoot(root bool) error
	Mount(remotePath string) error
	Chmod(path string, mod string, recursive bool) error

	Run(command ...string) (string, error)

	Fastboot(fastbootCli *fastboot.Fastboot) (*fastboot.Device, error)
}
