package cli

import (
	"context"
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CommandlineTool struct {
	binPath string
	logger  *pterm.Logger
}

func NewCLI(binPath string) *CommandlineTool {
	return &CommandlineTool{
		binPath: binPath,
		logger:  nil,
	}
}

func NewCLIWithLogger(binPath string, logger *pterm.Logger) *CommandlineTool {
	return &CommandlineTool{
		binPath: binPath,
		logger:  logger,
	}
}

func (cli *CommandlineTool) Check() error {
	info, err := os.Stat(cli.binPath)

	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("%s is not an executable file", cli.binPath)
	}

	return nil
}

func (cli *CommandlineTool) Run(command ...string) (string, error) {
	cli.tryLog(command...)
	rawOutput, result := cli.buildCommand(command...).CombinedOutput()
	return string(rawOutput), result
}

func (cli *CommandlineTool) RunWithContext(context context.Context, command ...string) *BufferedOutput {
	cli.tryLog(command...)
	var result BufferedOutput

	executableCommand := cli.buildCommandWithContext(context, command...)
	executableCommand.Stdout = &result.Out
	executableCommand.Stderr = &result.Err

	result.Error = executableCommand.Start()
	return &result
}

func (cli *CommandlineTool) tryLog(command ...string) {
	if cli.logger != nil {
		binName := filepath.Base(cli.binPath)
		cli.logger.Trace(fmt.Sprintf("%s %s", binName, strings.Trim(fmt.Sprint(command), "[]")))
	}
}

func (cli *CommandlineTool) buildCommand(command ...string) *exec.Cmd {
	return exec.Command(cli.binPath, command...)
}

func (cli *CommandlineTool) buildCommandWithContext(ctx context.Context, command ...string) *exec.Cmd {
	return exec.CommandContext(ctx, cli.binPath, command...)
}
