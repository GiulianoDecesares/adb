package adb

import "context"

type IDevice interface {
	GetSerial() string
	GetProduct() (string, error)

	Reboot(ctx context.Context) error
	WaitUntilReady(ctx context.Context) error
}
