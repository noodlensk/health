// Package checks handles health checks suites
package checks

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/noodlensk/health/pkg/exec"
	"github.com/noodlensk/health/pkg/report"
	"github.com/noodlensk/health/pkg/reporters"
)

// Checks holds information about health checks
type Checks struct {
	suites   map[string]Suite
	executor exec.Executor
	reporter reporters.Reporters
	logger   *zap.SugaredLogger
}

// Params holds params for new check
type Params struct {
	Config   Config
	Executor exec.Executor
	Reporter reporters.Reporters
	Logger   *zap.SugaredLogger
}

// New returns new checks
func New(params Params) Checks {
	s := Checks{
		suites:   map[string]Suite{},
		executor: params.Executor,
		reporter: params.Reporter,
		logger:   params.Logger,
	}

	for _, st := range params.Config.Suites {
		s.suites[st.Name] = st
	}

	return s
}

// Exec health check
func (c Checks) Exec(ctx context.Context, name string) (*report.Report, error) {
	st, ok := c.suites[name]
	if !ok {
		return nil, errors.Errorf("unknown Checks %q", name)
	}

	if err := c.executor.Exec(ctx, st.Command()); err != nil {
		return nil, errors.Wrap(err, "exec command")
	}

	r, err := report.Parse(st.ReportPath)
	if err != nil {
		return nil, errors.Wrap(err, "parse report")
	}

	return r, nil
}

// Serve executes all checks according to to it configuration
func (c Checks) Serve(ctx context.Context) error {
	wg := sync.WaitGroup{}

	for _, st := range c.suites {
		wg.Add(1)

		go func(suite Suite) {
			defer wg.Done()

			ticker := time.NewTicker(suite.RunEvery)

			for {
				select {
				case <-ticker.C:
					if err := c.executor.Exec(ctx, suite.Command()); err != nil {
						c.logger.With("error", err.Error()).Errorf("failed to execute %q suite", suite.Name)
						// TODO: add metrics about this

						continue
					}

					r, err := report.Parse(suite.ReportPath)
					if err != nil {
						c.logger.With("error", err.Error()).Errorf("failed to parse report for %q suite", suite.Name)

						continue
					}

					if err := r.Validate(); err != nil {
						c.logger.With("error", err.Error()).Errorf("report validation error for %q suite", suite.Name)

						continue
					}

					if err := c.reporter.Report(ctx, suite.Name, r); err != nil {
						c.logger.With("error", err.Error()).Errorf("failed to report for %q suite", suite.Name)
					}
				case <-ctx.Done():
					if ctx.Err() != context.Canceled {
						c.logger.With("error", ctx.Err()).Errorf("context canceled while executing suite %q", suite.Name)
					}

					c.logger.Infof("context canceled while executing suite %q", suite.Name)

					return
				}
			}
		}(st)
	}

	wg.Wait()

	return nil
}
