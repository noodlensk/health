// Package reporters contains list of reporters
package reporters

import (
	"context"

	"go.uber.org/zap"

	"github.com/noodlensk/health/pkg/report"
)

// Reporter defines interface for creating reporter
type Reporter interface {
	Report(ctx context.Context, suite string, rep *report.Report) error
}

// Reporters holds multiple reporters
type Reporters struct {
	reporters map[string]Reporter
	logger    *zap.SugaredLogger
}

// Options holds options for new Reporters
type Options struct {
	Reporters map[string]Reporter
	Logger    *zap.SugaredLogger
}

// New returns new reporters
func New(opts Options) Reporters {
	r := Reporters{
		reporters: opts.Reporters,
		logger:    opts.Logger,
	}

	return r
}

// Report reports health checks status
func (r *Reporters) Report(ctx context.Context, suite string, rep *report.Report) error {
	for name, reporter := range r.reporters {
		if err := reporter.Report(ctx, suite, rep); err != nil {
			r.logger.With("error", err.Error()).Errorf("failed execture reporter for %q", name)
		}
	}

	return nil
}
