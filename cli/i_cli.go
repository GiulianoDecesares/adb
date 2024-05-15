package cli

import "context"

type ICommandLineTool interface {
	Check() error

	Run(command ...string) (string, error)
	RunWithContext(context context.Context, command ...string) *BufferedOutput
}
