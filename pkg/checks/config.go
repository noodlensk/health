// Package checks handles health checks suites
package checks

import (
	"time"

	"github.com/noodlensk/health/pkg/exec"
)

// Suite represents set of health checks
type Suite struct {
	Name        string            `yaml:"name"`
	ExecCommand []string          `yaml:"execCommand"`
	Env         map[string]string `yaml:"env"`
	RunEvery    time.Duration     `yaml:"runEvery"`
	ReportPath  string            `yaml:"reportPath"`
}

// Command return Command to pass it to exec
func (s Suite) Command() exec.Command {
	cmd := exec.Command{}

	if len(s.ExecCommand) > 0 {
		cmd.Name = s.ExecCommand[0]
	}

	if len(s.ExecCommand) > 1 {
		cmd.Args = s.ExecCommand[1:]
	}

	cmd.Env = s.Env

	return cmd
}

// Config is a config :obvious:
type Config struct {
	Suites []Suite `yaml:"suites"`
}
