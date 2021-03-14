// Package opsgenie contains opsgenie reporter
package opsgenie

// Config holds info about config for OpsGenie
type Config struct {
	APIKey                 string            `yaml:"APIKey"`
	OpsGenieAPIURL         string            `yaml:"OpsGenieAPIURL"`
	DefaultAlertProperties map[string]string `yaml:"DefaultAlertProperties"`
}
