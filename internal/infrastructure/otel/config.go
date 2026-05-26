package otel

import "time"

// Config holds OpenTelemetry client settings.
type Config struct {
	Enabled              bool          `mapstructure:"Enabled" `
	Endpoint             string        `mapstructure:"Endpoint" `
	ServiceName          string        `mapstructure:"ServiceName" `
	Environment          string        `mapstructure:"Environment" `
	Insecure             bool          `mapstructure:"Insecure" `
	TraceSampleRatio     float64       `mapstructure:"TraceSampleRatio"`
	MetricExportInterval time.Duration `mapstructure:"MetricExportInterval"`
	MetricExportTimeout  time.Duration `mapstructure:"MetricExportTimeout"`
}

const (
	defaultMetricExportInterval = 60 * time.Second
	defaultMetricExportTimeout  = 30 * time.Second
)

// ApplyDefaults preserves previous behavior when new fields are not configured.
func (c *Config) ApplyDefaults() {
	if c == nil {
		return
	}
	if c.MetricExportInterval == 0 {
		c.MetricExportInterval = defaultMetricExportInterval
	}
	if c.MetricExportTimeout == 0 {
		c.MetricExportTimeout = defaultMetricExportTimeout
	}
}
