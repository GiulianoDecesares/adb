package adb

import "context"

type IDevice interface {
	GetId() string
	GetModel() string
	GetProduct() string

	IsPackageInstalled(packageName string) bool
	Install(packagePath string, overwrite bool) error
	Uninstall(packageName string) error

	ForceStop(packageName string) error
	RunActivity(name string, extraParameters ...string) error

	Pull(remotePath string, localPath string) error
	Push(localPath string, remotePath string) error

	DeleteFile(remotePath string) error
	DeleteDir(remotePath string) error

	WakeUp() error

	ListDirectory(directory string) ([]string, error)
	IsFile(deviceFilePath string) bool

	Logcat(context context.Context) *BufferedOutput
	LogcatWithFilter(context context.Context, filter string) *BufferedOutput

	Release()
}
