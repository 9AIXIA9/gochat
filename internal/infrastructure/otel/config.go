package otel

// Config holds OpenTelemetry client settings.
type Config struct {
	Endpoint    string `mapstructure:"Endpoint" `
	ServiceName string `mapstructure:"ServiceName" `
	Environment string `mapstructure:"Environment" `
	Insecure    bool   `mapstructure:"Insecure" `
}
