// Package config holds configs for the app
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/noodlensk/health/internal/common/server/http"
	"github.com/noodlensk/health/pkg/checks"
	"github.com/noodlensk/health/pkg/reporters/opsgenie"
)

// Config is app config
type Config struct {
	Name string
	Env  string

	Server struct {
		HTTP http.Config
	}

	Reporters struct {
		OpsGenie struct {
			Enabled bool
			Config  opsgenie.Config
		} `yaml:"OpsGenie"`
	}

	Checks checks.Config
}

// Validate config file
func (c *Config) Validate() error {
	var errMsgs []string

	if c.Name == "" {
		errMsgs = append(errMsgs, "empty Name")
	}

	if c.Env == "" {
		errMsgs = append(errMsgs, "empty Env")
	}

	if c.Server.HTTP.Address == "" {
		errMsgs = append(errMsgs, "empty Server.HTTP.Address")
	}

	if len(c.Checks.Suites) == 0 {
		errMsgs = append(errMsgs, "suites number should be > 0")
	}

	for i, s := range c.Checks.Suites {
		if s.Name == "" {
			errMsgs = append(errMsgs, fmt.Sprintf("suites[%d]: empty Name", i))
		}

		if s.ReportPath == "" {
			errMsgs = append(errMsgs, fmt.Sprintf("suites[%d]: empty ReportPath", i))
		}

		if len(s.ExecCommand) == 0 {
			errMsgs = append(errMsgs, fmt.Sprintf("suites[%d]: empty ExecCommand", i))
		}

		if s.RunEvery.Nanoseconds() == 0 {
			errMsgs = append(errMsgs, fmt.Sprintf("suites[%d]: empty RunEvery", i))
		}
	}

	if c.Reporters.OpsGenie.Enabled {
		if c.Reporters.OpsGenie.Config.APIKey == "" {
			errMsgs = append(errMsgs, "empty Reporters.OpsGenie.Config.APIKey")
		}
	}

	if len(errMsgs) > 0 {
		return errors.Errorf("errors: %s", strings.Join(errMsgs, ", "))
	}

	return nil
}

// Parse file into config
func Parse(filepath string) (*Config, error) {
	f, err := os.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	cfg := &Config{}

	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, errors.Wrap(err, "yaml unmarshal")
	}

	return cfg, nil
}
