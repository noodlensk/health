// Package prometheus contains prometheus reporter
package prometheus

import (
	"context"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/noodlensk/health/internal/common/server/http"
	"github.com/noodlensk/health/pkg/report"
)

// Reporter holds set of metrics for Prometheus
type Reporter struct {
	healthChecksPassedCount     *prom.GaugeVec
	healthChecksFailedCount     *prom.GaugeVec
	healthChecksTotalCount      *prom.GaugeVec
	healthChecksDurationSeconds *prom.GaugeVec
}

// Options for reporter
type Options struct {
	ProjectName string
	HTTPServer  *http.Server
}

// New returns new metrics
func New(opts Options) *Reporter {
	s := &Reporter{
		healthChecksPassedCount: prom.NewGaugeVec(
			prom.GaugeOpts{
				Name:        "health_checks_passed_count",
				Help:        "Health checks passed count",
				ConstLabels: map[string]string{"project": opts.ProjectName},
			},
			[]string{"suite", "group", "check"},
		),
		healthChecksFailedCount: prom.NewGaugeVec(
			prom.GaugeOpts{
				Name:        "health_checks_failed_count",
				Help:        "Health checks failed count",
				ConstLabels: map[string]string{"project": opts.ProjectName},
			},
			[]string{"suite", "group", "check"},
		), healthChecksTotalCount: prom.NewGaugeVec(
			prom.GaugeOpts{
				Name:        "health_checks_total_count",
				Help:        "Health checks total count",
				ConstLabels: map[string]string{"project": opts.ProjectName},
			},
			[]string{"suite", "group", "check"},
		),
		healthChecksDurationSeconds: prom.NewGaugeVec(
			prom.GaugeOpts{
				Name:        "health_checks_duration_seconds",
				Help:        "Health checks duration seconds",
				ConstLabels: map[string]string{"project": opts.ProjectName},
			},
			[]string{"suite", "group", "check"},
		),
	}

	prom.MustRegister(s.healthChecksPassedCount)
	prom.MustRegister(s.healthChecksFailedCount)
	prom.MustRegister(s.healthChecksTotalCount)
	prom.MustRegister(s.healthChecksDurationSeconds)

	opts.HTTPServer.Handle("GET", "/metrics", promhttp.Handler())

	return s
}

// Report update metric exposed to prometheus
func (r *Reporter) Report(ctx context.Context, suite string, rep *report.Report) error {
	for _, check := range rep.Checks {
		r.healthChecksPassedCount.WithLabelValues(suite, check.Group, check.Name).Set(0)
		r.healthChecksFailedCount.WithLabelValues(suite, check.Group, check.Name).Set(0)

		if check.Status == report.CheckStatusPassed {
			r.healthChecksPassedCount.WithLabelValues(suite, check.Group, check.Name).Inc()
		} else {
			r.healthChecksFailedCount.WithLabelValues(suite, check.Group, check.Name).Inc()
		}

		r.healthChecksDurationSeconds.WithLabelValues(suite, check.Group, check.Name).Set(check.Took.Seconds())
	}

	return nil
}
