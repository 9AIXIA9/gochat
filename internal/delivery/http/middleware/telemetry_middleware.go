package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// NewTelemetryMiddleware creates spans and attaches trace & span IDs into gin context for later logging.
func NewTelemetryMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName + "/http")
	return func(c *gin.Context) {
		start := time.Now().UTC()
		ctx, span := tracer.Start(c.Request.Context(), c.Request.Method+" "+c.FullPath(), trace.WithSpanKind(trace.SpanKindServer))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		dur := time.Since(start).Seconds()
		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", path),
			attribute.Int("http.status_code", status),
			attribute.Float64("http.duration_seconds", dur),
		)
		if status >= 500 && len(c.Errors) > 0 {
			span.RecordError(c.Errors.Last())
		}
		span.End()
		// Attach IDs for logger middleware enrichment
		c.Set("trace_id", span.SpanContext().TraceID().String())
		c.Set("span_id", span.SpanContext().SpanID().String())
	}
}
