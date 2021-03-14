// Package opsgenie contains opsgenie reporter
package opsgenie

import (
	"context"

	// nolint:gosec
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/noodlensk/health/pkg/report"
)

// Reporter is a reporter
type Reporter struct {
	opsgenieClient    *alert.Client
	defaultProperties map[string]string
}

// Options represents options for reporter
type Options struct {
	Config Config
}

// New returns new reporter
func New(opts Options) (*Reporter, error) {
	r := &Reporter{defaultProperties: opts.Config.DefaultAlertProperties}

	logrusLogger := logrus.New()
	logrusLogger.Out = io.Discard

	alertClient, err := alert.NewClient(&client.Config{
		ApiKey:         opts.Config.APIKey,
		OpsGenieAPIURL: client.ApiUrl(opts.Config.OpsGenieAPIURL),
		Logger:         logrusLogger,
	})
	if err != nil {
		return nil, errors.Wrap(err, "ops gennie client")
	}

	r.opsgenieClient = alertClient

	return r, nil
}

// Report reports creates/closes alerts in OpsGenie
func (r Reporter) Report(ctx context.Context, suite string, rep *report.Report) error {
	for _, c := range rep.Checks {
		// TODO: close alert
		if c.Status == report.CheckStatusFailed {
			if err := r.createAlert(ctx, suite, c); err != nil {
				return errors.Wrapf(err, "create alert for %s:%s:%s", suite, c.Group, c.Name)
			}
		}
	}

	return nil
}

func (r Reporter) createAlert(ctx context.Context, suite string, ch report.Check) error {
	_, err := r.opsgenieClient.Create(ctx, &alert.CreateAlertRequest{
		Message:     fmt.Sprintf("Health check failed: %s %s", ch.Group, ch.Name),
		Alias:       alias(suite, ch),
		Description: fmt.Sprintf("Error: %q", ch.Error),
		Priority:    priority(ch.Priority),
		Details:     r.defaultProperties,
	})

	return err
}

func alias(suite string, c report.Check) string {
	hash := md5.Sum([]byte(suite + c.Group + c.Name)) //nolint:gosec

	return hex.EncodeToString(hash[:])
}

func priority(rep report.CheckPriority) alert.Priority {
	m := map[report.CheckPriority]alert.Priority{
		report.CheckPriorityCritical:      alert.P1,
		report.CheckPriorityHigh:          alert.P2,
		report.CheckPriorityModerate:      alert.P3,
		report.CheckPriorityLow:           alert.P4,
		report.CheckPriorityInformational: alert.P5,
	}

	if v, ok := m[rep]; ok {
		return v
	}

	return alert.P1
}
