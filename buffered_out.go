package adb

import "bytes"

type BufferedOutput struct {
	Error error

	Out bytes.Buffer
	Err bytes.Buffer
}
