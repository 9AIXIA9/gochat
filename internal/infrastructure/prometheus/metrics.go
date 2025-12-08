package prometheus

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var runtimeCollectorsOnce sync.Once

// Metrics centralizes Prometheus collectors used across the app.
// It registers with the provided Registerer (or the default one when nil).
// Labels are carefully chosen to avoid high cardinality (e.g. HTTP path is the gin route pattern).
type Metrics struct {
	registry prom.Registerer

	// HTTP
	HTTPInFlight   prom.Gauge
	HTTPRequestDur *prom.HistogramVec
	HTTPRequests   *prom.CounterVec

	// Kafka
	KafkaProduced  *prom.CounterVec   // labels: topic, status
	KafkaConsumed  *prom.CounterVec   // labels: topic, status
	KafkaHandleDur *prom.HistogramVec // labels: topic

	// UnitOfWork (DB transaction scope)
	UOWDur   *prom.HistogramVec // labels: result
	UOWTotal *prom.CounterVec   // labels: result

	// WebSocket
	WSConnections prom.Gauge         // current connection count
	WSMessagesIn  *prom.CounterVec   // labels: type (ping|request|other)
	WSMessagesOut *prom.CounterVec   // labels: type (pong|response|broadcast)
	WSRouteDur    *prom.HistogramVec // labels: topic
	WSWriteErrors *prom.CounterVec   // labels: op (text|ping)
	WSReadErrors  *prom.CounterVec   // labels: kind (close|other)

	// Outbox (CDC)
	OutboxPolled    prom.Counter       // total binlog row events polled (insert on unpublished_events)
	OutboxBatchSize prom.Summary       // batch size of events read from outbox
	OutboxPublishes *prom.CounterVec   // labels: status (success|error)
	OutboxDur       *prom.HistogramVec // labels: status
}

func NewMetrics(reg prom.Registerer) *Metrics {
	if reg == nil {
		reg = prom.DefaultRegisterer
	}

	// Register process and Go runtime collectors only once when using the default registry
	if reg == prom.DefaultRegisterer {
		runtimeCollectorsOnce.Do(func() {
			// Use Register instead of MustRegister so we can ignore AlreadyRegisteredError if it ever appears.
			if err := reg.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
				var alreadyRegisteredError prom.AlreadyRegisteredError
				if !errors.As(err, &alreadyRegisteredError) {
					panic(fmt.Errorf("register process collector failed: %w", err))
				}
			}
			if err := reg.Register(collectors.NewGoCollector()); err != nil {
				var alreadyRegisteredError prom.AlreadyRegisteredError
				if !errors.As(err, &alreadyRegisteredError) {
					panic(fmt.Errorf("register go collector failed: %w", err))
				}
			}
		})
	} else {
		// For custom registry, safe to register every time.
		reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		reg.MustRegister(collectors.NewGoCollector())
	}

	m := &Metrics{registry: reg}

	// HTTP
	m.HTTPInFlight = promauto.With(reg).NewGauge(prom.GaugeOpts{
		Namespace: "gochat",
		Subsystem: "http",
		Name:      "in_flight_requests",
		Help:      "Current number of in-flight HTTP requests.",
	})
	m.HTTPRequestDur = promauto.With(reg).NewHistogramVec(prom.HistogramOpts{
		Namespace: "gochat",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "HTTP request duration in seconds.",
		Buckets:   prom.DefBuckets,
	}, []string{"method", "path", "status"})
	m.HTTPRequests = promauto.With(reg).NewCounterVec(prom.CounterOpts{
		Namespace: "gochat",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests.",
	}, []string{"method", "path", "status"})

	// Kafka
	m.KafkaProduced = promauto.With(reg).NewCounterVec(prom.CounterOpts{
		Namespace: "gochat",
		Subsystem: "kafka",
		Name:      "messages_produced_total",
		Help:      "Total number of Kafka messages produced.",
	}, []string{"topic", "status"})
	m.KafkaConsumed = promauto.With(reg).NewCounterVec(prom.CounterOpts{
		Namespace: "gochat",
		Subsystem: "kafka",
		Name:      "messages_consumed_total",
		Help:      "Total number of Kafka messages consumed/handled.",
	}, []string{"topic", "status"})
	m.KafkaHandleDur = promauto.With(reg).NewHistogramVec(prom.HistogramOpts{
		Namespace: "gochat",
		Subsystem: "kafka",
		Name:      "handler_duration_seconds",
		Help:      "Kafka handler duration in seconds.",
		Buckets:   prom.DefBuckets,
	}, []string{"topic"})

	// UOW
	m.UOWDur = promauto.With(reg).NewHistogramVec(prom.HistogramOpts{
		Namespace: "gochat",
		Subsystem: "uow",
		Name:      "duration_seconds",
		Help:      "UnitOfWork execution duration in seconds.",
		Buckets:   prom.DefBuckets,
	}, []string{"result"})
	m.UOWTotal = promauto.With(reg).NewCounterVec(prom.CounterOpts{
		Namespace: "gochat",
		Subsystem: "uow",
		Name:      "total",
		Help:      "Total UnitOfWork executions by result.",
	}, []string{"result"})

	// WebSocket
	m.WSConnections = promauto.With(reg).NewGauge(prom.GaugeOpts{
		Namespace: "gochat",
		Subsystem: "ws",
		Name:      "connections",
		Help:      "Current number of active WebSocket connections.",
	})
	m.WSMessagesIn = promauto.With(reg).NewCounterVec(prom.CounterOpts{
		Namespace: "gochat",
		Subsystem: "ws",
		Name:      "messages_in_total",
		Help:      "Total number of WebSocket messages received.",
	}, []string{"type"})
	m.WSMessagesOut = promauto.With(reg).NewCounterVec(prom.CounterOpts{
		Namespace: "gochat",
		Subsystem: "ws",
		Name:      "messages_out_total",
		Help:      "Total number of WebSocket messages sent.",
	}, []string{"type"})
	m.WSRouteDur = promauto.With(reg).NewHistogramVec(prom.HistogramOpts{
		Namespace: "gochat",
		Subsystem: "ws",
		Name:      "route_duration_seconds",
		Help:      "Duration of WebSocket route handling in seconds.",
		Buckets:   prom.DefBuckets,
	}, []string{"topic"})
	m.WSWriteErrors = promauto.With(reg).NewCounterVec(prom.CounterOpts{
		Namespace: "gochat",
		Subsystem: "ws",
		Name:      "write_errors_total",
		Help:      "Total number of WebSocket write errors.",
	}, []string{"op"})
	m.WSReadErrors = promauto.With(reg).NewCounterVec(prom.CounterOpts{
		Namespace: "gochat",
		Subsystem: "ws",
		Name:      "read_errors_total",
		Help:      "Total number of WebSocket read errors.",
	}, []string{"kind"})

	// Outbox
	m.OutboxPolled = promauto.With(reg).NewCounter(prom.CounterOpts{
		Namespace: "gochat",
		Subsystem: "outbox",
		Name:      "polled_total",
		Help:      "Total number of outbox binlog insert events observed.",
	})
	m.OutboxBatchSize = promauto.With(reg).NewSummary(prom.SummaryOpts{
		Namespace:  "gochat",
		Subsystem:  "outbox",
		Name:       "batch_size",
		Help:       "Batch size of outbox events published.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	m.OutboxPublishes = promauto.With(reg).NewCounterVec(prom.CounterOpts{
		Namespace: "gochat",
		Subsystem: "outbox",
		Name:      "publishes_total",
		Help:      "Total number of outbox publish attempts.",
	}, []string{"status"})
	m.OutboxDur = promauto.With(reg).NewHistogramVec(prom.HistogramOpts{
		Namespace: "gochat",
		Subsystem: "outbox",
		Name:      "publish_duration_seconds",
		Help:      "Duration of outbox publish batches in seconds.",
		Buckets:   prom.DefBuckets,
	}, []string{"status"})

	return m
}

// Handler returns an http.Handler serving Prometheus metrics from the registry.
func (m *Metrics) Handler() http.Handler { return promhttp.Handler() }

// GinMiddleware instruments incoming HTTP requests.
func (m *Metrics) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.HTTPInFlight.Inc()
		start := time.Now().UTC()
		c.Next()
		m.HTTPInFlight.Dec()
		statusCode := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}
		labels := prom.Labels{"method": c.Request.Method, "path": path, "status": fmt.Sprintf("%d", statusCode)}
		dur := time.Since(start).Seconds()
		m.HTTPRequestDur.With(labels).Observe(dur)
		m.HTTPRequests.With(labels).Inc()
	}
}
