// Package exec responsible for execution of commands
package exec

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Command represents command to execute
type Command struct {
	Name string
	Args []string
	Env  map[string]string
}

// New returns new executor
func New(logger *zap.SugaredLogger) Executor {
	return Executor{logger: logger}
}

// Executor executes commands
type Executor struct {
	logger *zap.SugaredLogger
}

// Exec executes command
func (e Executor) Exec(ctx context.Context, command Command) error {
	cmd := exec.CommandContext(ctx, command.Name, command.Args...) //nolint:gosec

	cmd.Env = os.Environ()

	if command.Env != nil {
		for k, v := range command.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	var wg sync.WaitGroup

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "get stdout pipe")
	}

	stdoutScanner := bufio.NewScanner(stdoutPipe)

	wg.Add(1)

	go func() {
		defer wg.Done()

		for stdoutScanner.Scan() {
			e.logger.With("output", stdoutScanner.Text()).Infow("command Stdout")
		}
	}()

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "get stderr pipe")
	}

	stderrScanner := bufio.NewScanner(stderrPipe)

	wg.Add(1)

	go func() {
		defer wg.Done()

		for stderrScanner.Scan() {
			e.logger.With("output", stderrScanner.Text()).Infow("command Stderr")
		}
	}()

	if err := cmd.Start(); err != nil {
		return err
	}

	err = cmd.Wait()

	wg.Wait()

	if err != nil {
		return errors.Wrap(err, "exec cmd")
	}

	return nil
}
