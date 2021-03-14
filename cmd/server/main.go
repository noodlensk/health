// Package server is a server for executing health checks
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/noodlensk/health/internal/common/config"
	"github.com/noodlensk/health/internal/common/logger"
	"github.com/noodlensk/health/internal/common/server/http"
	"github.com/noodlensk/health/pkg/checks"
	"github.com/noodlensk/health/pkg/exec"
	"github.com/noodlensk/health/pkg/reporters"
	"github.com/noodlensk/health/pkg/reporters/opsgenie"
	"github.com/noodlensk/health/pkg/reporters/prometheus"
)

var configPath string

func init() { //nolint:gochecknoinits
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Println(err) // nolint:forbidigo
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Parse(configPath)
	if err != nil {
		return errors.Wrap(err, "parse config")
	}

	if err := cfg.Validate(); err != nil {
		return errors.Wrap(err, "validate config")
	}

	log, err := logger.New()
	if err != nil {
		return errors.Wrap(err, "create logger")
	}

	exc := exec.New(log)

	httpSrv := http.NewServer(cfg.Server.HTTP.Address)

	reps := map[string]reporters.Reporter{
		"prometheus": prometheus.New(prometheus.Options{
			ProjectName: cfg.Name,
			HTTPServer:  httpSrv,
		}),
	}

	if cfg.Reporters.OpsGenie.Enabled {
		opsResp, err := opsgenie.New(opsgenie.Options{Config: cfg.Reporters.OpsGenie.Config})
		if err != nil {
			return errors.Wrap(err, "init opsgenie responder")
		}

		reps["opsgenie"] = opsResp
	}

	rep := reporters.New(reporters.Options{
		Reporters: reps,
		Logger:    log,
	})

	st := checks.New(checks.Params{
		Config:   cfg.Checks,
		Executor: exc,
		Reporter: rep,
		Logger:   log,
	})

	ctx := context.Background()

	go func() { log.Fatal(httpSrv.Serve()) }()

	return st.Serve(ctx)
}
