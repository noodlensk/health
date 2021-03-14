// Package report contains report information about executed health checks
package report

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Check represents health check structure
type Check struct {
	Name     string
	Group    string
	Status   CheckStatus
	Priority CheckPriority
	Error    string
	Took     time.Duration
}

// CheckPriority describes priority of check
type CheckPriority string

// possible check priority
const (
	CheckPriorityCritical      CheckPriority = "CRITICAL"
	CheckPriorityHigh          CheckPriority = "HIGH"
	CheckPriorityModerate      CheckPriority = "MODERATE"
	CheckPriorityLow           CheckPriority = "LOW"
	CheckPriorityInformational CheckPriority = "INFORMATIONAL"
)

// UnmarshalJSON is needed since time.Duration doesn't work out of the box
func (c *Check) UnmarshalJSON(data []byte) error {
	type alias struct {
		Name   string
		Group  string
		Status CheckStatus
		Error  string
		Took   string
	}

	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	t, err := time.ParseDuration(tmp.Took)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' to time.Duration: %v", tmp.Took, err)
	}

	c.Name = tmp.Name
	c.Group = tmp.Group
	c.Status = tmp.Status
	c.Error = tmp.Error
	c.Took = t

	return nil
}

// CheckStatus is a status of a check
type CheckStatus string

// Possible statuses of a check
const (
	CheckStatusPassed  CheckStatus = "PASSED"
	CheckStatusFailed  CheckStatus = "FAILED"
	CheckStatusUnknown CheckStatus = "UNKNOWN"
)

// Report holds information about executed set of health checks
type Report struct {
	SuitName string
	Created  time.Time
	Took     time.Duration
	Status   Status
	Checks   []Check
}

// UnmarshalJSON is needed since time.Duration doesn't work out of the box
func (r *Report) UnmarshalJSON(data []byte) error {
	type alias struct {
		SuitName string
		Created  time.Time
		Took     string
		Status   Status
		Checks   []Check
	}

	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	t, err := time.ParseDuration(tmp.Took)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' to time.Duration: %v", tmp.Took, err)
	}

	r.SuitName = tmp.SuitName
	r.Created = tmp.Created
	r.Status = tmp.Status
	r.Checks = tmp.Checks
	r.Took = t

	return nil
}

// Status represents status of report
type Status string

// Possible statuses of a report
const (
	StatusPassed  Status = "PASSED"
	StatusFailed  Status = "FAILED"
	StatusUnknown Status = "UNKNOWN"
)

// Validate validates report
func (r *Report) Validate() error {
	var errMsgs []string

	if len(r.Checks) == 0 {
		errMsgs = append(errMsgs, "empty checks")
	}

	if len(errMsgs) > 0 {
		return errors.Errorf("errors: %s", strings.Join(errMsgs, ", "))
	}

	return nil
}

// Parse parses file into report
func Parse(path string) (*Report, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}

	r := &Report{}

	if err := json.Unmarshal(f, &r); err != nil {
		return r, err
	}

	return r, nil
}
