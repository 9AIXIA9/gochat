package metrics

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	wsOnce sync.Once

	wsHandshakeTotal metric.Int64Counter
	wsActiveConns    metric.Int64UpDownCounter
	wsDisconnects    metric.Int64Counter
	wsReadErrors     metric.Int64Counter
	wsWriteErrors    metric.Int64Counter
	wsSendFailures   metric.Int64Counter
	wsMessageInTotal metric.Int64Counter
	wsMessageOut     metric.Int64Counter
	wsRouteDuration  metric.Float64Histogram

	kafkaOnce sync.Once

	kafkaProduceTotal metric.Int64Counter
	kafkaQueueFull    metric.Int64Counter
	kafkaConsumeTotal metric.Int64Counter
	kafkaCommitFailed metric.Int64Counter
	kafkaConsumerErr  metric.Int64Counter
	kafkaHandleDur    metric.Float64Histogram

	binlogOnce sync.Once

	binlogRunTotal metric.Int64Counter
)

func initWebsocketInstruments() {
	m := otel.Meter("gochat.websocket")
	wsHandshakeTotal, _ = m.Int64Counter("ws_handshake_total")
	wsActiveConns, _ = m.Int64UpDownCounter("ws_active_connections")
	wsDisconnects, _ = m.Int64Counter("ws_disconnect_total")
	wsReadErrors, _ = m.Int64Counter("ws_read_error_total")
	wsWriteErrors, _ = m.Int64Counter("ws_write_error_total")
	wsSendFailures, _ = m.Int64Counter("ws_send_failure_total")
	wsMessageInTotal, _ = m.Int64Counter("ws_message_in_total")
	wsMessageOut, _ = m.Int64Counter("ws_message_out_total")
	wsRouteDuration, _ = m.Float64Histogram("ws_route_duration_seconds")
}

func initKafkaInstruments() {
	m := otel.Meter("gochat.kafka")
	kafkaProduceTotal, _ = m.Int64Counter("kafka_produce_total")
	kafkaQueueFull, _ = m.Int64Counter("kafka_queue_full_total")
	kafkaConsumeTotal, _ = m.Int64Counter("kafka_consume_total")
	kafkaCommitFailed, _ = m.Int64Counter("kafka_commit_failed_total")
	kafkaConsumerErr, _ = m.Int64Counter("kafka_consumer_error_total")
	kafkaHandleDur, _ = m.Float64Histogram("kafka_handle_duration_seconds")
}

func initBinlogInstruments() {
	m := otel.Meter("gochat.binlog")
	binlogRunTotal, _ = m.Int64Counter("binlog_reader_run_total")
}

func WSHandshake(ctx context.Context, result string) {
	wsOnce.Do(initWebsocketInstruments)
	if wsHandshakeTotal == nil {
		return
	}
	wsHandshakeTotal.Add(ctx, 1, metric.WithAttributes(attribute.String("result", result)))
}

func WSConnectionDelta(ctx context.Context, delta int64) {
	wsOnce.Do(initWebsocketInstruments)
	if wsActiveConns == nil {
		return
	}
	wsActiveConns.Add(ctx, delta)
}

func WSDisconnect(ctx context.Context, reason string) {
	wsOnce.Do(initWebsocketInstruments)
	if wsDisconnects == nil {
		return
	}
	wsDisconnects.Add(ctx, 1, metric.WithAttributes(attribute.String("reason", reason)))
}

func WSReadError(ctx context.Context, reason string) {
	wsOnce.Do(initWebsocketInstruments)
	if wsReadErrors == nil {
		return
	}
	wsReadErrors.Add(ctx, 1, metric.WithAttributes(attribute.String("reason", reason)))
}

func WSWriteError(ctx context.Context, reason string) {
	wsOnce.Do(initWebsocketInstruments)
	if wsWriteErrors == nil {
		return
	}
	wsWriteErrors.Add(ctx, 1, metric.WithAttributes(attribute.String("reason", reason)))
}

func WSSendFailure(ctx context.Context, reason string) {
	wsOnce.Do(initWebsocketInstruments)
	if wsSendFailures == nil {
		return
	}
	wsSendFailures.Add(ctx, 1, metric.WithAttributes(attribute.String("reason", reason)))
}

func WSMessageIn(ctx context.Context, source string) {
	wsOnce.Do(initWebsocketInstruments)
	if wsMessageInTotal == nil {
		return
	}
	wsMessageInTotal.Add(ctx, 1, metric.WithAttributes(attribute.String("source", source)))
}

func WSMessageOut(ctx context.Context, source string, result string) {
	wsOnce.Do(initWebsocketInstruments)
	if wsMessageOut == nil {
		return
	}
	wsMessageOut.Add(ctx, 1, metric.WithAttributes(
		attribute.String("source", source),
		attribute.String("result", result),
	))
}

func WSRouteDuration(ctx context.Context, topic string, result string, seconds float64) {
	wsOnce.Do(initWebsocketInstruments)
	if wsRouteDuration == nil {
		return
	}
	wsRouteDuration.Record(ctx, seconds, metric.WithAttributes(
		attribute.String("topic", topic),
		attribute.String("result", result),
	))
}

func KafkaProduce(ctx context.Context, result string, errorType string) {
	kafkaOnce.Do(initKafkaInstruments)
	if kafkaProduceTotal == nil {
		return
	}
	kafkaProduceTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String("result", result),
		attribute.String("error_type", errorType),
	))
}

func KafkaQueueFull(ctx context.Context) {
	kafkaOnce.Do(initKafkaInstruments)
	if kafkaQueueFull == nil {
		return
	}
	kafkaQueueFull.Add(ctx, 1)
}

func KafkaConsume(ctx context.Context, result string, topic string) {
	kafkaOnce.Do(initKafkaInstruments)
	if kafkaConsumeTotal == nil {
		return
	}
	kafkaConsumeTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String("result", result),
		attribute.String("topic", topic),
	))
}

func KafkaCommitFail(ctx context.Context, topic string) {
	kafkaOnce.Do(initKafkaInstruments)
	if kafkaCommitFailed == nil {
		return
	}
	kafkaCommitFailed.Add(ctx, 1, metric.WithAttributes(attribute.String("topic", topic)))
}

func KafkaConsumerError(ctx context.Context, code string) {
	kafkaOnce.Do(initKafkaInstruments)
	if kafkaConsumerErr == nil {
		return
	}
	kafkaConsumerErr.Add(ctx, 1, metric.WithAttributes(attribute.String("code", code)))
}

func KafkaHandleDuration(ctx context.Context, topic string, result string, seconds float64) {
	kafkaOnce.Do(initKafkaInstruments)
	if kafkaHandleDur == nil {
		return
	}
	kafkaHandleDur.Record(ctx, seconds, metric.WithAttributes(
		attribute.String("topic", topic),
		attribute.String("result", result),
	))
}

func BinlogReaderRun(ctx context.Context, result string, errorType string) {
	binlogOnce.Do(initBinlogInstruments)
	if binlogRunTotal == nil {
		return
	}
	binlogRunTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String("result", result),
		attribute.String("error_type", errorType),
	))
}
