package gorm

import (
	"context"
	"strings"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DBMetrics holds Prometheus collectors specific to database operations.
type DBMetrics struct {
	QueryDur   *prom.HistogramVec // labels: op, table, slow
	QueryTotal *prom.CounterVec   // labels: op, table, slow, status
}

func NewDBMetrics(reg prom.Registerer) *DBMetrics {
	if reg == nil {
		reg = prom.DefaultRegisterer
	}
	return &DBMetrics{
		QueryDur: prom.NewHistogramVec(prom.HistogramOpts{
			Namespace: "gochat",
			Subsystem: "db",
			Name:      "query_duration_seconds",
			Help:      "GORM SQL execution duration.",
			Buckets:   prom.DefBuckets,
		}, []string{"op", "table", "slow"}),
		QueryTotal: prom.NewCounterVec(prom.CounterOpts{
			Namespace: "gochat",
			Subsystem: "db",
			Name:      "queries_total",
			Help:      "Total executed SQL statements.",
		}, []string{"op", "table", "slow", "status"}),
	}
}

// Register registers db metrics with registry.
func (m *DBMetrics) Register(reg prom.Registerer) {
	if reg == nil {
		reg = prom.DefaultRegisterer
	}
	reg.MustRegister(m.QueryDur, m.QueryTotal)
}

// ObservabilityPlugin instruments GORM with tracing and metrics.
type ObservabilityPlugin struct {
	tracer        trace.Tracer
	metrics       *DBMetrics
	slowThreshold time.Duration
}

func NewObservabilityPlugin(metrics *DBMetrics, slowMillis int) *ObservabilityPlugin {
	if slowMillis <= 0 {
		slowMillis = 200 // default 200ms
	}
	return &ObservabilityPlugin{
		tracer:        otel.Tracer("gorm"),
		metrics:       metrics,
		slowThreshold: time.Duration(slowMillis) * time.Millisecond,
	}
}

func (p *ObservabilityPlugin) Name() string { return "observability" }

func (p *ObservabilityPlugin) Initialize(db *gorm.DB) error {
	cb := db.Callback()
	// Register before callbacks
	if err := cb.Create().Before("gorm:create").Register("observe:before_create", p.before("create")); err != nil {
		return err
	}
	if err := cb.Query().Before("gorm:query").Register("observe:before_query", p.before("select")); err != nil {
		return err
	}
	if err := cb.Update().Before("gorm:update").Register("observe:before_update", p.before("update")); err != nil {
		return err
	}
	if err := cb.Delete().Before("gorm:delete").Register("observe:before_delete", p.before("delete")); err != nil {
		return err
	}
	if err := cb.Row().Before("gorm:row").Register("observe:before_row", p.before("row")); err != nil {
		return err
	}
	if err := cb.Raw().Before("gorm:raw").Register("observe:before_raw", p.before("raw")); err != nil {
		return err
	}
	// After callbacks
	if err := cb.Create().After("gorm:create").Register("observe:after_create", p.after("create")); err != nil {
		return err
	}
	if err := cb.Query().After("gorm:query").Register("observe:after_query", p.after("select")); err != nil {
		return err
	}
	if err := cb.Update().After("gorm:update").Register("observe:after_update", p.after("update")); err != nil {
		return err
	}
	if err := cb.Delete().After("gorm:delete").Register("observe:after_delete", p.after("delete")); err != nil {
		return err
	}
	if err := cb.Row().After("gorm:row").Register("observe:after_row", p.after("row")); err != nil {
		return err
	}
	if err := cb.Raw().After("gorm:raw").Register("observe:after_raw", p.after("raw")); err != nil {
		return err
	}
	return nil
}

type ctxKey struct{}

func (p *ObservabilityPlugin) before(op string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		ctx := db.Statement.Context
		if ctx == nil {
			ctx = context.Background()
		}
		spanCtx, span := p.tracer.Start(ctx, "db."+op, trace.WithSpanKind(trace.SpanKindClient))
		span.SetAttributes(attribute.String("db.system", "mysql"), attribute.String("db.operation", op))
		if db.Statement != nil && db.Statement.Table != "" {
			span.SetAttributes(attribute.String("db.table", db.Statement.Table))
		}
		db.Statement.Context = context.WithValue(spanCtx, ctxKey{}, span)
		db.InstanceSet("observe_start", time.Now())
	}
}

func (p *ObservabilityPlugin) after(op string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		startRaw, _ := db.InstanceGet("observe_start")
		start, _ := startRaw.(time.Time)
		dur := time.Since(start)
		slowFlag := dur >= p.slowThreshold
		table := "unknown"
		if db.Statement != nil && db.Statement.Table != "" {
			table = db.Statement.Table
		}
		table = strings.ToLower(table)
		status := "ok"
		if db.Error != nil {
			status = "error"
		}
		if p.metrics != nil {
			labelsH := prom.Labels{"op": op, "table": table, "slow": boolToStr(slowFlag)}
			p.metrics.QueryDur.With(labelsH).Observe(dur.Seconds())
			labelsC := prom.Labels{"op": op, "table": table, "slow": boolToStr(slowFlag), "status": status}
			p.metrics.QueryTotal.With(labelsC).Inc()
		}
		ctx := db.Statement.Context
		spanVal := ctx.Value(ctxKey{})
		span, _ := spanVal.(trace.Span)
		if span != nil {
			span.SetAttributes(attribute.String("db.table", table), attribute.Int64("db.duration_ms", dur.Milliseconds()))
			if slowFlag {
				span.SetAttributes(attribute.Bool("db.slow", true))
				span.AddEvent("slow_query", trace.WithAttributes(attribute.Int64("duration_ms", dur.Milliseconds())))
			}
			if db.Error != nil {
				span.RecordError(db.Error)
				span.SetStatus(codes.Error, db.Error.Error())
			}
			span.End()
		}
		if slowFlag {
			zap.L().Warn("slow SQL", zap.String("op", op), zap.String("table", table), zap.Duration("duration", dur))
		}
	}
}

func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
