package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type config struct {
	url            string
	origin         string
	clients        int
	duration       time.Duration
	connectRate    int
	connectTimeout time.Duration
	maxThreads     int
	pingInterval   time.Duration
	sendInterval   time.Duration
	sendMode       string
	payload        string
	content        string
	latencyField   string
	reportInterval time.Duration
	authHeader     string
	cookie         string
	tokensFile     string
	recipientID    string
	recipientsFile string
	chaosDropRatio float64
	chaosDropAfter time.Duration
	// paired mode: "paired" enables sender/receiver pairing where
	// the first N tokens/recipients are receivers and the next N are senders.
	mode  string
	pairs int
}

type metrics struct {
	attempted        atomic.Int64
	connected        atomic.Int64
	connectFailed    atomic.Int64
	active           atomic.Int64
	peakActive       atomic.Int64
	retained         atomic.Int64
	unexpectedClosed atomic.Int64
	chaosDropped     atomic.Int64
	pingsSent        atomic.Int64
	pingFailed       atomic.Int64
	messagesSent     atomic.Int64
	messageFailed    atomic.Int64
	messagesReceived atomic.Int64
	latencyMatched   atomic.Int64
	readErrors       atomic.Int64
	connectLatencyNs atomic.Int64
	messageLatencyNs atomic.Int64

	mu              sync.Mutex
	connectLatencyS []float64
	messageLatencyS []float64

	// Paired-mode delivery metrics
	deliverySuccess atomic.Int64
	deliveryFailed  atomic.Int64
	duplicates      atomic.Int64
	orderErrors     atomic.Int64
}

func main() {
	cfg := parseFlags()
	tokens, err := loadTokens(cfg.tokensFile)
	if err != nil {
		log.Fatalf("load tokens failed: %v", err)
	}
	recipients, err := loadStringLines(cfg.recipientsFile)
	if err != nil {
		log.Fatalf("load recipients failed: %v", err)
	}
	if len(tokens) > 0 {
		log.Printf("loaded %d tokens from %s", len(tokens), cfg.tokensFile)
	}
	if len(recipients) > 0 {
		log.Printf("loaded %d recipients from %s", len(recipients), cfg.recipientsFile)
	}
	if err := validateConfig(cfg, len(tokens), len(recipients)); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	applyRuntimeTuning(cfg)

	var m metrics
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), cfg.duration)
	defer cancel()

	log.Printf(
		"start websocket benchmark url=%s clients=%d duration=%s connect_rate=%d/s max_threads=%d ping_interval=%s send_interval=%s",
		cfg.url,
		cfg.clients,
		cfg.duration,
		cfg.connectRate,
		cfg.maxThreads,
		cfg.pingInterval,
		cfg.sendInterval,
	)

	start := time.Now()

	if cfg.reportInterval > 0 {
		go reportLoop(ctx, cfg.reportInterval, &m)
	}

	launchClients(ctx, cfg, tokens, recipients, &wg, &m)
	wg.Wait()
	elapsed := time.Since(start)

	printSummary(cfg, &m, elapsed)
	// generate markdown report file
	if err := writeReport(cfg, &m, elapsed); err != nil {
		log.Printf("write report failed: %v", err)
	}
}

func writeReport(cfg config, m *metrics, elapsed time.Duration) error {
	ts := time.Now().UTC().Format("20060102-150405")
	reportPath := fmt.Sprintf("websocket_benchmark_report_%s.md", ts)

	attempted := m.attempted.Load()
	connected := m.connected.Load()
	msgSent := m.messagesSent.Load()
	msgRecv := m.messagesReceived.Load()
	deliverySucc := m.deliverySuccess.Load()
	deliveryFail := m.deliveryFailed.Load()
	dup := m.duplicates.Load()
	orderErr := m.orderErrors.Load()

	connP50, connP95, connP99 := connectLatencyPercentiles(m)
	msgP50, msgP95, msgP99 := messageLatencyPercentiles(m)

	content := fmt.Sprintf(`# GoChat 压测报告

测试时间: %s UTC
测试时长: %s
URL: %s
Clients: %d
ConnectRate: %d/s

核心指标
--------
- attempted connections: %d
- successful connections: %d
- messages sent: %d
- messages received: %d
- delivery success/fail: %d / %d
- duplicate messages: %d
- order errors: %d

延迟指标
--------
- connect latency p50/p95/p99: %.2fms / %.2fms / %.2fms
- message latency p50/p95/p99: %.2fms / %.2fms / %.2fms

`, ts, elapsed, cfg.url, cfg.clients, cfg.connectRate, attempted, connected, msgSent, msgRecv, deliverySucc, deliveryFail, dup, orderErr, connP50*1000, connP95*1000, connP99*1000, msgP50*1000, msgP95*1000, msgP99*1000)

	if err := os.WriteFile(reportPath, []byte(content), 0o644); err != nil {
		return err
	}
	log.Printf("report written: %s", reportPath)
	return nil
}

func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.url, "url", "ws://localhost:8080/api/v1/ws/", "websocket endpoint")
	flag.StringVar(&cfg.origin, "origin", "", "Origin header value for websocket handshake, e.g. 'http://localhost:5173'")
	flag.IntVar(&cfg.clients, "clients", 1000, "total clients to create")
	flag.DurationVar(&cfg.duration, "duration", 60*time.Second, "test duration")
	flag.IntVar(&cfg.connectRate, "connect-rate", 200, "new connections per second (0 means no ramp-up)")
	flag.DurationVar(&cfg.connectTimeout, "connect-timeout", 5*time.Second, "dial timeout")
	flag.IntVar(&cfg.maxThreads, "max-threads", 50000, "Go runtime max OS threads for this process; increase for very high concurrency benchmark")
	flag.StringVar(&cfg.mode, "mode", "", "optional mode: paired")
	flag.IntVar(&cfg.pairs, "pairs", 0, "number of sender/receiver pairs when -mode paired; defaults to clients/2 if 0")
	flag.DurationVar(&cfg.pingInterval, "ping-interval", 5*time.Second, "ping interval (0 disables ping)")
	flag.DurationVar(&cfg.sendInterval, "send-interval", 0, "text message send interval (0 disables send)")
	flag.StringVar(&cfg.sendMode, "send-mode", "raw", "send mode: raw or private")
	flag.StringVar(&cfg.payload, "payload", `{"type":"benchmark","bench_ts_ns":"{{ts_unix_nano}}","bench_client":"{{client_id}}","bench_seq":"{{seq}}","content":"ping"}`, "text payload when send-interval > 0; supports placeholders {{ts_unix_nano}}/{{client_id}}/{{seq}}")
	flag.StringVar(&cfg.content, "content", "bench_ts_ns={{ts_unix_nano}};bench_client={{client_id}};bench_seq={{seq}}", "message content template for protocol send mode; supports placeholders {{ts_unix_nano}}/{{client_id}}/{{seq}}")
	flag.StringVar(&cfg.latencyField, "latency-field", "bench_ts_ns", "json field name used to compute end-to-end latency from received message")
	flag.DurationVar(&cfg.reportInterval, "report-interval", 5*time.Second, "periodic metrics print interval (0 disables)")
	flag.StringVar(&cfg.authHeader, "auth", "", "Authorization header value, e.g. 'Bearer <token>'")
	flag.StringVar(&cfg.cookie, "cookie", "", "Cookie header value, e.g. 'access_token=<token>'")
	flag.StringVar(&cfg.tokensFile, "tokens-file", "", "path to token file, one token per line (supports '<token>' or 'Bearer <token>')")
	flag.StringVar(&cfg.recipientID, "recipient-id", "", "recipient user id for private send mode")
	flag.StringVar(&cfg.recipientsFile, "recipients-file", "", "path to recipient user id file, one id per line for private send mode")
	flag.Float64Var(&cfg.chaosDropRatio, "chaos-drop-ratio", 0, "ratio of clients that will drop tcp abruptly after chaos-drop-after (0~1)")
	flag.DurationVar(&cfg.chaosDropAfter, "chaos-drop-after", 0, "delay after connect before chaos drop (0 disables)")
	flag.Parse()
	return cfg
}

func validateConfig(cfg config, tokenCount int, recipientCount int) error {
	if cfg.clients <= 0 {
		return fmt.Errorf("clients must be > 0")
	}
	if cfg.duration <= 0 {
		return fmt.Errorf("duration must be > 0")
	}
	if cfg.connectRate < 0 {
		return fmt.Errorf("connect-rate must be >= 0")
	}
	if cfg.connectTimeout <= 0 {
		return fmt.Errorf("connect-timeout must be > 0")
	}
	if cfg.maxThreads <= 0 {
		return fmt.Errorf("max-threads must be > 0")
	}
	if cfg.pingInterval < 0 || cfg.sendInterval < 0 || cfg.reportInterval < 0 {
		return fmt.Errorf("intervals must be >= 0")
	}
	if cfg.sendMode != "raw" && cfg.sendMode != "private" {
		return fmt.Errorf("send-mode must be one of: raw, private")
	}
	if recipientCount > 0 && strings.TrimSpace(cfg.recipientID) != "" {
		return fmt.Errorf("recipients-file cannot be used together with -recipient-id")
	}
	if cfg.sendInterval > 0 && cfg.sendMode == "private" && recipientCount == 0 && strings.TrimSpace(cfg.recipientID) == "" {
		return fmt.Errorf("private send mode requires -recipient-id or -recipients-file")
	}
	if cfg.chaosDropRatio < 0 || cfg.chaosDropRatio > 1 {
		return fmt.Errorf("chaos-drop-ratio must be in [0, 1]")
	}
	if cfg.chaosDropAfter < 0 {
		return fmt.Errorf("chaos-drop-after must be >= 0")
	}
	if cfg.chaosDropRatio > 0 && cfg.chaosDropAfter == 0 {
		return fmt.Errorf("chaos-drop-after must be > 0 when chaos-drop-ratio > 0")
	}
	if tokenCount > 0 {
		if tokenCount < cfg.clients {
			return fmt.Errorf("tokens-file has %d tokens but clients=%d; need at least one token per client for multi-user test", tokenCount, cfg.clients)
		}
		if strings.TrimSpace(cfg.authHeader) != "" || strings.TrimSpace(cfg.cookie) != "" {
			return fmt.Errorf("tokens-file cannot be used together with -auth or -cookie")
		}
	}
	if tokenCount == 0 && strings.TrimSpace(cfg.authHeader) == "" && strings.TrimSpace(cfg.cookie) == "" {
		return fmt.Errorf("provide one of -tokens-file, -auth, or -cookie")
	}
	// paired mode specific checks
	if strings.EqualFold(cfg.mode, "paired") {
		if recipientCount == 0 {
			return fmt.Errorf("paired mode requires -recipients-file with userIDs")
		}
		if cfg.pairs < 0 {
			return fmt.Errorf("pairs must be >= 0")
		}
		if cfg.pairs == 0 {
			// require clients to be even
			if cfg.clients%2 != 0 {
				return fmt.Errorf("clients must be even when using paired mode and pairs not set")
			}
		} else {
			if cfg.clients != cfg.pairs*2 {
				return fmt.Errorf("clients must equal pairs*2 when pairs is set")
			}
		}
	}
	return nil
}

func applyRuntimeTuning(cfg config) {
	prev := debug.SetMaxThreads(cfg.maxThreads)
	log.Printf("runtime max threads set: %d (previous: %d)", cfg.maxThreads, prev)
}

func launchClients(ctx context.Context, cfg config, tokens []string, recipients []string, wg *sync.WaitGroup, m *metrics) {
	launch := func(id int) {
		wg.Add(1)
		go runClient(ctx, id, cfg, tokens, recipients, m, wg)
	}

	if cfg.connectRate == 0 {
		for i := 0; i < cfg.clients; i++ {
			launch(i)
		}
		return
	}

	interval := time.Second / time.Duration(cfg.connectRate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for i := 0; i < cfg.clients; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			launch(i)
		}
	}
}

func runClient(ctx context.Context, id int, cfg config, tokens []string, recipients []string, m *metrics, wg *sync.WaitGroup) {
	defer wg.Done()
	m.attempted.Add(1)

	dialer := websocket.Dialer{HandshakeTimeout: cfg.connectTimeout}
	headers := make(http.Header)
	authHeader := pickAuthHeader(cfg, tokens, id)
	if strings.TrimSpace(cfg.origin) != "" {
		headers.Set("Origin", strings.TrimSpace(cfg.origin))
	}
	if authHeader != "" {
		headers.Set("Authorization", authHeader)
	}
	if strings.TrimSpace(cfg.cookie) != "" {
		headers.Set("Cookie", cfg.cookie)
	}

	connStart := time.Now()
	conn, _, err := dialer.DialContext(ctx, cfg.url, headers)
	if err != nil {
		m.connectFailed.Add(1)
		log.Printf("client=%d connect failed: %v", id, err)
		return
	}
	defer func() { _ = conn.Close() }()

	lat := time.Since(connStart)
	m.connected.Add(1)
	activeNow := m.active.Add(1)
	updatePeak(m, activeNow)
	m.connectLatencyNs.Add(lat.Nanoseconds())

	m.mu.Lock()
	m.connectLatencyS = append(m.connectLatencyS, lat.Seconds())
	m.mu.Unlock()

	// Paired mode detection and local state for receivers
	isPaired := strings.EqualFold(cfg.mode, "paired")
	pairs := cfg.pairs
	if pairs <= 0 {
		pairs = cfg.clients / 2
	}
	isReceiver := false
	isSender := false
	var expectedSenderUserID string
	if isPaired {
		if id < pairs {
			isReceiver = true
			senderIdx := pairs + id
			if senderIdx >= 0 && senderIdx < len(recipients) {
				expectedSenderUserID = recipients[senderIdx]
			}
		} else {
			isSender = true
		}
	}

	var closeOnce sync.Once
	closeConn := func() {
		closeOnce.Do(func() {
			m.active.Add(-1)
			_ = conn.Close()
		})
	}
	defer closeConn()

	readDone := make(chan struct{})
	var readUnexpected atomic.Bool
	go func() {
		defer close(readDone)
		// receiver-side state
		lastSeq := make(map[string]int64)
		seenMsg := make(map[string]struct{})
		for {
			msgType, payload, readErr := conn.ReadMessage()
			if readErr != nil {
				// Context cancellation and normal close should not be counted as read errors.
				if ctx.Err() == nil {
					m.readErrors.Add(1)
					readUnexpected.Store(true)
				}
				return
			}
			if msgType == websocket.TextMessage || msgType == websocket.BinaryMessage {
				m.messagesReceived.Add(1)

				// try to parse as JSON message wrapper
				var decoded any
				if err := json.Unmarshal(payload, &decoded); err == nil {
					if docMap, ok := decoded.(map[string]any); ok {
						if topicRaw, ok := docMap["topic"]; ok {
							if topicStr, ok := topicRaw.(string); ok {
								// server push topic for private messages
								if topicStr == "chat.notify_private_message" {
									if body, ok := docMap["body"].(map[string]any); ok {
										msgID, _ := body["id"].(string)
										senderID, _ := body["sender_id"].(string)
										content, _ := body["content"].(string)

										if isReceiver {
											// verify sender mapping in paired mode
											if expectedSenderUserID != "" && senderID != expectedSenderUserID {
												m.deliveryFailed.Add(1)
												continue
											}

											if msgID != "" {
												if _, seen := seenMsg[msgID]; seen {
													m.duplicates.Add(1)
													continue
												}
												seenMsg[msgID] = struct{}{}
											}

											// order check from content template fields: bench_seq and sender identity
											seqVal, hasSeq := parseKVIntFromString(content, "bench_seq")
											if hasSeq {
												last := lastSeq[senderID]
												if last != 0 && seqVal != last+1 {
													m.orderErrors.Add(1)
												}
												lastSeq[senderID] = seqVal
											}

											if tsNs, ok := parseKVIntFromString(content, "bench_ts_ns"); ok {
												deltaNs := time.Now().UnixNano() - tsNs
												if deltaNs >= 0 {
													m.latencyMatched.Add(1)
													m.messageLatencyNs.Add(deltaNs)
													m.mu.Lock()
													m.messageLatencyS = append(m.messageLatencyS, float64(deltaNs)/float64(time.Second))
													m.mu.Unlock()
												}
											}

											m.deliverySuccess.Add(1)
										} else {
											// for non-receiver clients, still collect latency if possible
											if tsNs, ok := parseKVIntFromString(content, "bench_ts_ns"); ok {
												deltaNs := time.Now().UnixNano() - tsNs
												if deltaNs >= 0 {
													m.latencyMatched.Add(1)
													m.messageLatencyNs.Add(deltaNs)
													m.mu.Lock()
													m.messageLatencyS = append(m.messageLatencyS, float64(deltaNs)/float64(time.Second))
													m.mu.Unlock()
												}
											}
										}
									}
								}
							}
						}
					}
				} else {
					// not JSON; try to extract latency field as before
					if tsNs, ok := extractLatencyTimestamp(payload, cfg.latencyField); ok {
						deltaNs := time.Now().UnixNano() - tsNs
						if deltaNs >= 0 {
							m.latencyMatched.Add(1)
							m.messageLatencyNs.Add(deltaNs)
							m.mu.Lock()
							m.messageLatencyS = append(m.messageLatencyS, float64(deltaNs)/float64(time.Second))
							m.mu.Unlock()
						}
					}
				}
			}
		}
	}()

	var pingTicker *time.Ticker
	if cfg.pingInterval > 0 {
		pingTicker = time.NewTicker(cfg.pingInterval)
		defer pingTicker.Stop()
	}

	var sendTicker *time.Ticker
	if cfg.sendInterval > 0 {
		// in paired mode only sender side should send
		if !isPaired || isSender {
			sendTicker = time.NewTicker(cfg.sendInterval)
			defer sendTicker.Stop()
		}
	}

	var chaosTimer *time.Timer
	if cfg.chaosDropRatio > 0 && cfg.chaosDropAfter > 0 {
		chaosTimer = time.NewTimer(cfg.chaosDropAfter)
		defer chaosTimer.Stop()
	}

	var seq int64

	for {
		select {
		case <-ctx.Done():
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "benchmark end"), time.Now().Add(time.Second))
			select {
			case <-readDone:
			case <-time.After(2 * time.Second):
			}
			m.retained.Add(1)
			return
		case <-readDone:
			if readUnexpected.Load() {
				m.unexpectedClosed.Add(1)
			}
			return
		case <-tickChanTimer(chaosTimer):
			if shouldDropByRatio(id, cfg.chaosDropRatio) {
				m.chaosDropped.Add(1)
				_ = conn.Close()
				return
			}
			chaosTimer = nil
		case <-tickChan(pingTicker):
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second)); err != nil {
				m.pingFailed.Add(1)
				m.unexpectedClosed.Add(1)
				return
			}
			m.pingsSent.Add(1)
		case <-tickChan(sendTicker):
			payload := buildSendPayload(cfg, id, seq, time.Now(), recipients)
			seq++
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				m.messageFailed.Add(1)
				m.unexpectedClosed.Add(1)
				return
			}
			m.messagesSent.Add(1)
		}
	}
}

func loadTokens(path string) ([]string, error) {
	if strings.TrimSpace(path) == "" {
		return nil, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	tokens := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		tokens = append(tokens, normalizeBearer(line))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("tokens file is empty: %s", path)
	}
	return tokens, nil
}

func loadStringLines(path string) ([]string, error) {
	if strings.TrimSpace(path) == "" {
		return nil, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	items := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		items = append(items, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("file is empty: %s", path)
	}
	return items, nil
}

func normalizeBearer(token string) string {
	trimmed := strings.TrimSpace(token)
	if strings.HasPrefix(strings.ToLower(trimmed), "bearer ") {
		return "Bearer " + strings.TrimSpace(trimmed[len("bearer "):])
	}
	return "Bearer " + trimmed
}

func pickAuthHeader(cfg config, tokens []string, clientID int) string {
	if len(tokens) > 0 {
		return tokens[clientID]
	}
	if strings.TrimSpace(cfg.authHeader) == "" {
		return ""
	}
	return normalizeBearer(cfg.authHeader)
}

func tickChan(t *time.Ticker) <-chan time.Time {
	if t == nil {
		return nil
	}
	return t.C
}

func tickChanTimer(t *time.Timer) <-chan time.Time {
	if t == nil {
		return nil
	}
	return t.C
}

func updatePeak(m *metrics, candidate int64) {
	for {
		peak := m.peakActive.Load()
		if candidate <= peak {
			return
		}
		if m.peakActive.CompareAndSwap(peak, candidate) {
			return
		}
	}
}

func shouldDropByRatio(clientID int, ratio float64) bool {
	if ratio <= 0 {
		return false
	}
	if ratio >= 1 {
		return true
	}
	seed := int64(clientID)*1103515245 + 12345
	if seed < 0 {
		seed = -seed
	}
	v := float64(seed%1000000) / 1000000.0
	return v < ratio
}

func buildPayload(template string, clientID int, seq int64, now time.Time) []byte {
	s := template
	s = strings.ReplaceAll(s, "{{ts_unix_nano}}", fmt.Sprintf("%d", now.UnixNano()))
	s = strings.ReplaceAll(s, "{{client_id}}", fmt.Sprintf("%d", clientID))
	s = strings.ReplaceAll(s, "{{seq}}", fmt.Sprintf("%d", seq))
	return []byte(s)
}

func buildSendPayload(cfg config, clientID int, seq int64, now time.Time, recipients []string) []byte {
	if cfg.sendMode != "private" {
		return buildPayload(cfg.payload, clientID, seq, now)
	}

	// determine recipient: in paired mode we map sender clients to receivers deterministically
	var recipient string
	if strings.EqualFold(cfg.mode, "paired") {
		// number of pairs
		pairs := cfg.pairs
		if pairs <= 0 {
			pairs = cfg.clients / 2
		}
		if clientID >= pairs {
			// sender side: map to receiver index
			idx := clientID - pairs
			if idx >= 0 && idx < len(recipients) {
				recipient = recipients[idx]
			}
		}
		// if recipient still empty, fall back to default picker
		if recipient == "" {
			recipient = pickRecipient(cfg, clientID, seq, recipients)
		}
	} else {
		recipient = pickRecipient(cfg, clientID, seq, recipients)
	}

	// build payload with metadata for paired verification
	body := map[string]any{
		"topic": "chat.send_private_message",
		"payload": map[string]any{
			"recipient_id":  recipient,
			"sender_client": clientID,
			"seq":           seq,
			"message_id":    fmt.Sprintf("%d-%d-%d", clientID, seq, now.UnixNano()),
			"ts_ns":         now.UnixNano(),
			"content":       string(buildPayload(cfg.content, clientID, seq, now)),
		},
	}
	b, err := json.Marshal(body)
	if err != nil {
		// marshal failure on simple map should be impossible; keep a safe fallback.
		return []byte(`{"topic":"chat.send_private_message","payload":{}}`)
	}
	return b
}

func pickRecipient(cfg config, clientID int, seq int64, recipients []string) string {
	if len(recipients) == 0 {
		return strings.TrimSpace(cfg.recipientID)
	}
	idx := (int64(clientID) + seq) % int64(len(recipients))
	return recipients[idx]
}

func extractLatencyTimestamp(payload []byte, field string) (int64, bool) {
	if strings.TrimSpace(field) == "" {
		return 0, false
	}
	var decoded any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return 0, false
	}
	if ts, ok := findLatencyField(decoded, field); ok {
		return ts, true
	}
	return findLatencyInStrings(decoded, field)
}

func findLatencyField(v any, field string) (int64, bool) {
	switch typed := v.(type) {
	case map[string]any:
		if raw, ok := typed[field]; ok {
			return toInt64(raw)
		}
		for _, nested := range typed {
			if ts, ok := findLatencyField(nested, field); ok {
				return ts, true
			}
		}
	case []any:
		for _, nested := range typed {
			if ts, ok := findLatencyField(nested, field); ok {
				return ts, true
			}
		}
	}
	return 0, false
}

func toInt64(v any) (int64, bool) {
	switch typed := v.(type) {
	case float64:
		return int64(typed), true
	case int64:
		return typed, true
	case int:
		return int64(typed), true
	case string:
		var parsed int64
		if _, err := fmt.Sscanf(strings.TrimSpace(typed), "%d", &parsed); err == nil {
			return parsed, true
		}
	}
	return 0, false
}

func findLatencyInStrings(v any, field string) (int64, bool) {
	switch typed := v.(type) {
	case map[string]any:
		for _, nested := range typed {
			if ts, ok := findLatencyInStrings(nested, field); ok {
				return ts, true
			}
		}
	case []any:
		for _, nested := range typed {
			if ts, ok := findLatencyInStrings(nested, field); ok {
				return ts, true
			}
		}
	case string:
		if ts, ok := parseLatencyFromString(typed, field); ok {
			return ts, true
		}
	}
	return 0, false
}

func parseLatencyFromString(s string, field string) (int64, bool) {
	for _, token := range strings.FieldsFunc(s, func(r rune) bool {
		return r == ';' || r == ',' || r == '|' || r == ' ' || r == '\n' || r == '\t'
	}) {
		parts := strings.SplitN(strings.TrimSpace(token), "=", 2)
		if len(parts) != 2 {
			continue
		}
		if strings.TrimSpace(parts[0]) != field {
			continue
		}
		return toInt64(strings.TrimSpace(parts[1]))
	}
	return 0, false
}

func parseKVIntFromString(s, field string) (int64, bool) {
	for _, token := range strings.FieldsFunc(s, func(r rune) bool {
		return r == ';' || r == ',' || r == '|' || r == ' ' || r == '\n' || r == '\t'
	}) {
		parts := strings.SplitN(strings.TrimSpace(token), "=", 2)
		if len(parts) != 2 {
			continue
		}
		if strings.TrimSpace(parts[0]) != field {
			continue
		}
		return toInt64(strings.TrimSpace(parts[1]))
	}
	return 0, false
}

func reportLoop(ctx context.Context, interval time.Duration, m *metrics) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			attempted := m.attempted.Load()
			connected := m.connected.Load()
			active := m.active.Load()
			peak := m.peakActive.Load()
			sendOK := m.messagesSent.Load()
			sendFail := m.messageFailed.Load()
			recv := m.messagesReceived.Load()
			log.Printf("progress attempted=%d connected=%d active=%d peak=%d retained=%d unexpected=%d msg_send_ok=%d msg_send_fail=%d msg_recv=%d ping_ok=%d ping_fail=%d",
				attempted, connected, active, peak, m.retained.Load(), m.unexpectedClosed.Load(), sendOK, sendFail, recv, m.pingsSent.Load(), m.pingFailed.Load())
		}
	}
}

func printSummary(cfg config, m *metrics, elapsed time.Duration) {
	attempted := m.attempted.Load()
	connected := m.connected.Load()
	retained := m.retained.Load()

	successRate := 0.0
	if attempted > 0 {
		successRate = float64(connected) / float64(attempted) * 100
	}

	retainedRate := 0.0
	if connected > 0 {
		retainedRate = float64(retained) / float64(connected) * 100
	}

	avgConnMs := 0.0
	if connected > 0 {
		avgConnMs = float64(m.connectLatencyNs.Load()) / float64(connected) / float64(time.Millisecond)
	}

	elapsedSec := elapsed.Seconds()
	if elapsedSec <= 0 {
		elapsedSec = 1
	}

	msgSendTPS := float64(m.messagesSent.Load()) / elapsedSec
	msgRecvTPS := float64(m.messagesReceived.Load()) / elapsedSec

	sendTotal := m.messagesSent.Load() + m.messageFailed.Load()
	sendSuccessRate := 0.0
	if sendTotal > 0 {
		sendSuccessRate = float64(m.messagesSent.Load()) / float64(sendTotal) * 100
	}

	receiveDeliveryRate := 0.0
	if m.messagesSent.Load() > 0 {
		receiveDeliveryRate = float64(m.messagesReceived.Load()) / float64(m.messagesSent.Load()) * 100
	}

	latencyMatchRate := 0.0
	if m.messagesReceived.Load() > 0 {
		latencyMatchRate = float64(m.latencyMatched.Load()) / float64(m.messagesReceived.Load()) * 100
	}

	connP50, connP95, connP99 := connectLatencyPercentiles(m)
	msgP50, msgP95, msgP99 := messageLatencyPercentiles(m)

	avgMsgLatencyMs := 0.0
	if m.latencyMatched.Load() > 0 {
		avgMsgLatencyMs = float64(m.messageLatencyNs.Load()) / float64(m.latencyMatched.Load()) / float64(time.Millisecond)
	}

	log.Printf("benchmark complete")
	log.Printf("url: %s", cfg.url)
	log.Printf("duration: %s", elapsed)
	log.Printf("attempted connections: %d", attempted)
	log.Printf("successful connections: %d", connected)
	log.Printf("connection success rate: %.2f%%", successRate)
	log.Printf("active peak connections: %d", m.peakActive.Load())
	log.Printf("retained connections: %d", retained)
	log.Printf("retained connection rate: %.2f%%", retainedRate)
	log.Printf("unexpected closed connections: %d", m.unexpectedClosed.Load())
	log.Printf("chaos dropped connections: %d", m.chaosDropped.Load())
	log.Printf("avg connect latency: %.2fms", avgConnMs)
	log.Printf("connect latency p50/p95/p99: %.2fms / %.2fms / %.2fms", connP50*1000, connP95*1000, connP99*1000)
	log.Printf("messages sent(success/fail): %d / %d", m.messagesSent.Load(), m.messageFailed.Load())
	log.Printf("message send success rate: %.2f%%", sendSuccessRate)
	log.Printf("messages received: %d", m.messagesReceived.Load())
	log.Printf("message send/recv throughput(tps): %.2f / %.2f", msgSendTPS, msgRecvTPS)
	log.Printf("message receive delivery rate(received/sent): %.2f%%", receiveDeliveryRate)
	log.Printf("latency matched messages: %d (match rate %.2f%%, field=%s)", m.latencyMatched.Load(), latencyMatchRate, cfg.latencyField)
	if m.latencyMatched.Load() > 0 {
		log.Printf("avg end-to-end message latency: %.2fms", avgMsgLatencyMs)
		log.Printf("message latency p50/p95/p99: %.2fms / %.2fms / %.2fms", msgP50*1000, msgP95*1000, msgP99*1000)
	}
	// Paired-mode delivery metrics (if applicable)
	if m.deliverySuccess.Load()+m.deliveryFailed.Load() > 0 {
		totalDelivery := m.deliverySuccess.Load() + m.deliveryFailed.Load()
		deliveryRate := float64(m.deliverySuccess.Load()) / float64(totalDelivery) * 100
		log.Printf("delivery success/fail: %d / %d (success_rate=%.2f%%)", m.deliverySuccess.Load(), m.deliveryFailed.Load(), deliveryRate)
		log.Printf("duplicate messages: %d", m.duplicates.Load())
		log.Printf("order errors: %d", m.orderErrors.Load())
	}
	log.Printf("pings sent(success/fail): %d / %d", m.pingsSent.Load(), m.pingFailed.Load())
	log.Printf("read errors: %d", m.readErrors.Load())
}

func connectLatencyPercentiles(m *metrics) (p50, p95, p99 float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.connectLatencyS) == 0 {
		return 0, 0, 0
	}
	vals := make([]float64, len(m.connectLatencyS))
	copy(vals, m.connectLatencyS)
	sort.Float64s(vals)
	return percentile(vals, 0.50), percentile(vals, 0.95), percentile(vals, 0.99)
}

func messageLatencyPercentiles(m *metrics) (p50, p95, p99 float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.messageLatencyS) == 0 {
		return 0, 0, 0
	}
	vals := make([]float64, len(m.messageLatencyS))
	copy(vals, m.messageLatencyS)
	sort.Float64s(vals)
	return percentile(vals, 0.50), percentile(vals, 0.95), percentile(vals, 0.99)
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 1 {
		return sorted[len(sorted)-1]
	}
	idx := int(math.Ceil(float64(len(sorted))*p)) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
