package otel

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"
)

type TelemetryConfig struct {
	Enabled        bool    `mapstructure:"TelemetryEnabled"`
	TraceEnabled   bool    `mapstructure:"TraceEnabled"`
	MetricsEnabled bool    `mapstructure:"MetricsEnabled"`
	SampleRatio    float64 `mapstructure:"SampleRatio"`
	OTLPEndpoint   string  `mapstructure:"OTLPEndpoint"`
	ServiceVersion string  `mapstructure:"ServiceVersion"`
	Environment    string  `mapstructure:"Environment"`
}

func (c *TelemetryConfig) Validate() error {
	if c == nil || !c.Enabled {
		return nil
	}

	if c.SampleRatio < 0 || c.SampleRatio > 1 {
		return fmt.Errorf("Telemetry.SampleRatio: %w: must be in [0,1], got %f", myErrors.ErrInvalidNumber, c.SampleRatio)
	}
	if c.TraceEnabled && c.OTLPEndpoint == "" {
		return fmt.Errorf("Telemetry.OTLPEndpoint: %w", myErrors.ErrEmptyInput)
	}

	return nil
}
