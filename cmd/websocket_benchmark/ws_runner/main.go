package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	actionSendPrivateMessage = "send_private_message"
	actionSendRoomMessage    = "send_room_message"
)

type upstreamEnvelope struct {
	ClientMessageID string          `json:"client_message_id"`
	Action          string          `json:"action"`
	Payload         json.RawMessage `json:"payload"`
}

type privatePayload struct {
	RecipientID string `json:"recipient_id"`
	Content     string `json:"content"`
}

type roomPayload struct {
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
}

type ackEnvelope struct {
	ClientMessageID string `json:"client_message_id"`
	AckType         string `json:"ack_type"`
	Error           string `json:"error,omitempty"`
}

type downstreamEnvelope struct {
	Action          string          `json:"action"`
	Payload         json.RawMessage `json:"payload"`
	ClientMessageID string          `json:"client_message_id,omitempty"`
	Result          string          `json:"result,omitempty"`
	ErrorMessage    string          `json:"error_message,omitempty"`
}

type benchmarkMessagePayload struct {
	ID        string `json:"id,omitempty"`
	MessageID string `json:"message_id,omitempty"`
	Content   string `json:"content"`
}

type pairLine struct {
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	SendStatus  string `json:"send_status"`
	AgreeStatus string `json:"agree_status"`
}

type metrics struct {
	connected       int64
	connectFail     int64
	sent            int64
	sendFail        int64
	recv            int64
	ackMatched      int64
	ackUnmatched    int64
	ackTimedOut     int64
	ackReceived     int64
	ackError        int64
	ackLatencyNs    int64
	ackMaxLatencyNs int64
	pendingCurrent  int64
	readFail        int64
	writeFail       int64
}

type pendingMessage struct {
	sentAt    time.Time
	clientIdx int
	seq       int
}

type latencySample struct {
	Timestamp string  `json:"ts"`
	ClientID  int     `json:"client_id"`
	Sequence  int     `json:"sequence"`
	MessageID string  `json:"client_message_id"`
	AckType   string  `json:"ack_type"`
	LatencyMS float64 `json:"latency_ms"`
	SendMode  string  `json:"send_mode"`
}

type e2eLatencySample struct {
	Timestamp         string  `json:"ts"`
	ClientID          int     `json:"client_id"`
	Sequence          int     `json:"sequence"`
	ClientMessageID   string  `json:"client_message_id,omitempty"`
	BusinessMessageID string  `json:"business_message_id,omitempty"`
	CorrelationMethod string  `json:"correlation_method"`
	SentAtUnixNano    int64   `json:"sent_at_unix_nano"`
	RecvAtUnixNano    int64   `json:"recv_at_unix_nano"`
	LatencyMS         float64 `json:"latency_ms"`
	Content           string  `json:"content"`
	Action            string  `json:"action,omitempty"`
}

type clientPlan struct {
	ClientIdx    int
	Token        string
	SenderID     string
	ActiveSender bool
}

type sentMessageMeta struct {
	clientIdx       int
	seq             int
	sentAtUnixNano  int64
	clientMessageID string
}

type correlationStore struct {
	mu                  sync.RWMutex
	byClientMessageID   map[string]sentMessageMeta
	byBusinessMessageID map[string]sentMessageMeta
}

type latencyPoint struct {
	at      time.Time
	valueMS float64
}

type percentileSnapshot struct {
	Count int
	AvgMS float64
	MaxMS float64
	P50MS float64
	P90MS float64
	P95MS float64
	P99MS float64
}

type latencyStats struct {
	mu        sync.Mutex
	ackAll    []float64
	e2eAll    []float64
	ackPoints []latencyPoint
	e2ePoints []latencyPoint
}

var (
	wsURL = flag.String("url", "ws://localhost:8080/api/v1/ws/", "WebSocket URL")

	tokensFile     = flag.String("tokens-file", "./websocket_benchmark_data/output/tokens.txt", "tokens file")
	recipientsFile = flag.String("recipients-file", "./websocket_benchmark_data/output/recipients.txt", "recipients(user_id) file")
	roomsFile      = flag.String("rooms-file", "./websocket_benchmark_data/output/rooms.txt", "rooms file")
	pairsFile      = flag.String("pairs-file", "./websocket_benchmark_data/output/friendship_pairs.jsonl", "friendship pairs file (optional)")

	clients      = flag.Int("clients", 200, "number of ws clients")
	duration     = flag.Duration("duration", 60*time.Second, "benchmark duration")
	sendInterval = flag.Duration("send-interval", 1500*time.Millisecond, "send interval per client (0 means no send)")
	sendMode     = flag.String("send-mode", "private", "send mode: private|room|mixed")
	roomRatio    = flag.Float64("room-ratio", 0.5, "room traffic ratio when send-mode=mixed")

	connectRate       = flag.Int("connect-rate", 0, "connections per second, 0 means burst")
	origin            = flag.String("origin", "", "optional Origin header")
	ackTimeout        = flag.Duration("ack-timeout", 5*time.Second, "pending ack timeout; 0 disables timeout tracking")
	jsonOutFile       = flag.String("json-out-file", "./websocket_benchmark_data/output/ws_metrics.jsonl", "JSONL metrics output file (appended); empty disables")
	latencyOutFile    = flag.String("latency-out-file", "./websocket_benchmark_data/output/ws_latency.jsonl", "JSONL latency samples output file (appended); empty disables")
	e2eLatencyOutFile = flag.String("e2e-latency-out-file", "./websocket_benchmark_data/output/ws_e2e_latency.jsonl", "JSONL end-to-end latency samples output file (appended); empty disables")

	contentPrefix = flag.String("content-prefix", "bench", "message content prefix")
	seed          = flag.Int64("seed", 0, "rng seed, 0 means now")
	reportEvery   = flag.Duration("report-interval", 5*time.Second, "periodic report interval, 0 disables")
	windowEvery   = flag.Duration("window-interval", 5*time.Second, "latency percentile window size")
)

func newCorrelationStore() *correlationStore {
	return &correlationStore{
		byClientMessageID:   make(map[string]sentMessageMeta, 4096),
		byBusinessMessageID: make(map[string]sentMessageMeta, 4096),
	}
}

func (s *correlationStore) trackByClientMessageID(clientMessageID string, meta sentMessageMeta) {
	if clientMessageID == "" {
		return
	}
	s.mu.Lock()
	s.byClientMessageID[clientMessageID] = meta
	s.mu.Unlock()
}

func (s *correlationStore) bindBusinessMessageID(clientMessageID string, businessMessageID string) bool {
	if clientMessageID == "" || businessMessageID == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	meta, ok := s.byClientMessageID[clientMessageID]
	if !ok {
		return false
	}
	s.byBusinessMessageID[businessMessageID] = meta
	return true
}

func (s *correlationStore) resolveByClientMessageID(clientMessageID string) (sentMessageMeta, bool) {
	if clientMessageID == "" {
		return sentMessageMeta{}, false
	}
	s.mu.RLock()
	meta, ok := s.byClientMessageID[clientMessageID]
	s.mu.RUnlock()
	return meta, ok
}

func (s *correlationStore) resolveByBusinessMessageID(businessMessageID string) (sentMessageMeta, bool) {
	if businessMessageID == "" {
		return sentMessageMeta{}, false
	}
	s.mu.RLock()
	meta, ok := s.byBusinessMessageID[businessMessageID]
	s.mu.RUnlock()
	return meta, ok
}

func newLatencyStats() *latencyStats {
	return &latencyStats{
		ackAll:    make([]float64, 0, 4096),
		e2eAll:    make([]float64, 0, 4096),
		ackPoints: make([]latencyPoint, 0, 4096),
		e2ePoints: make([]latencyPoint, 0, 4096),
	}
}

func (s *latencyStats) addAck(latencyMS float64) {
	now := time.Now()
	s.mu.Lock()
	s.ackAll = append(s.ackAll, latencyMS)
	s.ackPoints = append(s.ackPoints, latencyPoint{at: now, valueMS: latencyMS})
	s.mu.Unlock()
}

func (s *latencyStats) addE2E(latencyMS float64) {
	now := time.Now()
	s.mu.Lock()
	s.e2eAll = append(s.e2eAll, latencyMS)
	s.e2ePoints = append(s.e2ePoints, latencyPoint{at: now, valueMS: latencyMS})
	s.mu.Unlock()
}

func collectWindowValues(points []latencyPoint, cutoff time.Time) []float64 {
	if len(points) == 0 {
		return nil
	}
	out := make([]float64, 0, len(points))
	for i := len(points) - 1; i >= 0; i-- {
		if points[i].at.Before(cutoff) {
			break
		}
		out = append(out, points[i].valueMS)
	}
	return out
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if len(sorted) == 1 {
		return sorted[0]
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 1 {
		return sorted[len(sorted)-1]
	}
	pos := p * float64(len(sorted)-1)
	low := int(math.Floor(pos))
	high := int(math.Ceil(pos))
	if low == high {
		return sorted[low]
	}
	frac := pos - float64(low)
	return sorted[low] + (sorted[high]-sorted[low])*frac
}

func summarizePercentiles(values []float64) percentileSnapshot {
	if len(values) == 0 {
		return percentileSnapshot{}
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	total := 0.0
	for _, v := range sorted {
		total += v
	}
	return percentileSnapshot{
		Count: len(sorted),
		AvgMS: total / float64(len(sorted)),
		MaxMS: sorted[len(sorted)-1],
		P50MS: percentile(sorted, 0.50),
		P90MS: percentile(sorted, 0.90),
		P95MS: percentile(sorted, 0.95),
		P99MS: percentile(sorted, 0.99),
	}
}

func (s *latencyStats) snapshots(window time.Duration) (ackAll percentileSnapshot, e2eAll percentileSnapshot, ackWindow percentileSnapshot, e2eWindow percentileSnapshot) {
	if window <= 0 {
		window = 5 * time.Second
	}
	cutoff := time.Now().Add(-window)

	s.mu.Lock()
	ackAllValues := append([]float64(nil), s.ackAll...)
	e2eAllValues := append([]float64(nil), s.e2eAll...)
	ackWindowValues := collectWindowValues(s.ackPoints, cutoff)
	e2eWindowValues := collectWindowValues(s.e2ePoints, cutoff)
	s.mu.Unlock()

	return summarizePercentiles(ackAllValues), summarizePercentiles(e2eAllValues), summarizePercentiles(ackWindowValues), summarizePercentiles(e2eWindowValues)
}

func extractBusinessMessageID(payload benchmarkMessagePayload) string {
	if strings.TrimSpace(payload.MessageID) != "" {
		return strings.TrimSpace(payload.MessageID)
	}
	return strings.TrimSpace(payload.ID)
}

func extractBusinessMessageIDFromRaw(raw json.RawMessage) string {
	return extractBusinessMessageIDFromRawDepth(raw, 0)
}

func extractBusinessMessageIDFromRawDepth(raw json.RawMessage, depth int) string {
	if depth > 2 || len(raw) == 0 {
		return ""
	}

	var payload benchmarkMessagePayload
	if err := json.Unmarshal(raw, &payload); err == nil {
		if id := extractBusinessMessageID(payload); id != "" {
			return id
		}
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ""
	}

	for _, key := range []string{"id", "message_id", "messageId", "msg_id"} {
		if v, ok := obj[key]; ok {
			var s string
			if err := json.Unmarshal(v, &s); err == nil && strings.TrimSpace(s) != "" {
				return strings.TrimSpace(s)
			}
		}
	}

	for _, key := range []string{"message", "data", "result", "payload"} {
		if v, ok := obj[key]; ok {
			if id := extractBusinessMessageIDFromRawDepth(v, depth+1); id != "" {
				return id
			}
		}
	}

	return ""
}

func main() {
	flag.Parse()

	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(*seed))
	corrStore := newCorrelationStore()
	latStats := newLatencyStats()

	tokens, err := readLines(*tokensFile)
	if err != nil {
		log.Fatalf("read tokens file failed: %v", err)
	}
	if len(tokens) == 0 {
		log.Fatalf("empty tokens file: %s", *tokensFile)
	}

	recipients, _ := readLines(*recipientsFile)
	rooms, _ := readLines(*roomsFile)
	pairMap, _ := readPairs(*pairsFile)

	plans := buildClientPlans(tokens, recipients, pairMap)
	if len(plans) == 0 {
		log.Fatalf("no runnable clients after pairing/filtering")
	}

	n := *clients
	if n > len(plans) {
		n = len(plans)
	}
	if n <= 0 {
		log.Fatalf("invalid clients: %d", *clients)
	}

	if *sendMode == "private" || *sendMode == "mixed" {
		if len(pairMap) == 0 && len(recipients) == 0 {
			log.Fatalf("private or mixed mode requires recipients-file or pairs-file")
		}
	}
	if *sendMode == "room" || *sendMode == "mixed" {
		if len(rooms) == 0 {
			log.Fatalf("room or mixed mode requires rooms-file")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()

	var m metrics
	var clientWG sync.WaitGroup
	var auxWG sync.WaitGroup
	var latencyCh chan latencySample
	var e2eLatencyCh chan e2eLatencySample
	if *latencyOutFile != "" {
		latencyCh = make(chan latencySample, 4096)
		auxWG.Add(1)
		go func() {
			defer auxWG.Done()
			if err := writeLatencySamples(*latencyOutFile, latencyCh); err != nil {
				log.Printf("write latency samples failed: %v", err)
			}
		}()
	}
	if *e2eLatencyOutFile != "" {
		e2eLatencyCh = make(chan e2eLatencySample, 4096)
		auxWG.Add(1)
		go func() {
			defer auxWG.Done()
			if err := writeE2ELatencySamples(*e2eLatencyOutFile, e2eLatencyCh); err != nil {
				log.Printf("write e2e latency samples failed: %v", err)
			}
		}()
	}

	if *reportEvery > 0 {
		auxWG.Add(1)
		go func() {
			defer auxWG.Done()
			t := time.NewTicker(*reportEvery)
			defer t.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					ackAll, e2eAll, ackWindow, e2eWindow := latStats.snapshots(*windowEvery)
					printMetrics("progress", &m, ackAll, e2eAll, ackWindow, e2eWindow)
					if *jsonOutFile != "" {
						if err := dumpMetricsJSON(*jsonOutFile, "progress", &m, ackAll, e2eAll, ackWindow, e2eWindow); err != nil {
							log.Printf("dump metrics json failed: %v", err)
						}
					}
				}
			}
		}()
	}

	connectTicker := (*time.Ticker)(nil)
	if *connectRate > 0 {
		connectTicker = time.NewTicker(time.Second / time.Duration(*connectRate))
		defer connectTicker.Stop()
	}

	for i := 0; i < n; i++ {
		if connectTicker != nil {
			<-connectTicker.C
		}
		plan := plans[i]
		clientWG.Add(1)
		go func(p clientPlan) {
			defer clientWG.Done()
			runClient(ctx, p.ClientIdx, p.Token, p.SenderID, p.ActiveSender, recipients, rooms, pairMap, rng.Int63(), corrStore, latStats, &m, latencyCh, e2eLatencyCh)
		}(plan)
	}

	<-ctx.Done()
	clientWG.Wait()
	if *latencyOutFile != "" {
		close(latencyCh)
	}
	if *e2eLatencyOutFile != "" {
		close(e2eLatencyCh)
	}
	auxWG.Wait()
	ackAll, e2eAll, ackWindow, e2eWindow := latStats.snapshots(*windowEvery)
	printMetrics("final", &m, ackAll, e2eAll, ackWindow, e2eWindow)
	if *jsonOutFile != "" {
		if err := dumpMetricsJSON(*jsonOutFile, "final", &m, ackAll, e2eAll, ackWindow, e2eWindow); err != nil {
			log.Printf("dump final metrics json failed: %v", err)
		}
	}
}

func dumpMetricsJSON(path string, prefix string, m *metrics, ackAll percentileSnapshot, e2eAll percentileSnapshot, ackWindow percentileSnapshot, e2eWindow percentileSnapshot) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	matched := atomic.LoadInt64(&m.ackMatched)
	avgLatencyMs := float64(0)
	if matched > 0 {
		avgLatencyMs = float64(atomic.LoadInt64(&m.ackLatencyNs)) / float64(time.Millisecond) / float64(matched)
	}

	obj := map[string]interface{}{
		"ts":                 time.Now().UTC().Format(time.RFC3339Nano),
		"prefix":             prefix,
		"connected":          atomic.LoadInt64(&m.connected),
		"connect_fail":       atomic.LoadInt64(&m.connectFail),
		"sent":               atomic.LoadInt64(&m.sent),
		"send_fail":          atomic.LoadInt64(&m.sendFail),
		"recv":               atomic.LoadInt64(&m.recv),
		"ack_matched":        matched,
		"ack_unmatched":      atomic.LoadInt64(&m.ackUnmatched),
		"ack_timed_out":      atomic.LoadInt64(&m.ackTimedOut),
		"pending":            atomic.LoadInt64(&m.pendingCurrent),
		"ack_received":       atomic.LoadInt64(&m.ackReceived),
		"ack_error":          atomic.LoadInt64(&m.ackError),
		"ack_avg_latency_ms": avgLatencyMs,
		"ack_max_latency_ms": float64(atomic.LoadInt64(&m.ackMaxLatencyNs)) / float64(time.Millisecond),
		"read_fail":          atomic.LoadInt64(&m.readFail),
		"write_fail":         atomic.LoadInt64(&m.writeFail),
		"ack_count":          ackAll.Count,
		"ack_p50_ms":         ackAll.P50MS,
		"ack_p90_ms":         ackAll.P90MS,
		"ack_p95_ms":         ackAll.P95MS,
		"ack_p99_ms":         ackAll.P99MS,
		"e2e_count":          e2eAll.Count,
		"e2e_p50_ms":         e2eAll.P50MS,
		"e2e_p90_ms":         e2eAll.P90MS,
		"e2e_p95_ms":         e2eAll.P95MS,
		"e2e_p99_ms":         e2eAll.P99MS,
		"ack_window_count":   ackWindow.Count,
		"ack_window_p50_ms":  ackWindow.P50MS,
		"ack_window_p90_ms":  ackWindow.P90MS,
		"ack_window_p95_ms":  ackWindow.P95MS,
		"ack_window_p99_ms":  ackWindow.P99MS,
		"e2e_window_count":   e2eWindow.Count,
		"e2e_window_p50_ms":  e2eWindow.P50MS,
		"e2e_window_p90_ms":  e2eWindow.P90MS,
		"e2e_window_p95_ms":  e2eWindow.P95MS,
		"e2e_window_p99_ms":  e2eWindow.P99MS,
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(b, '\n')); err != nil {
		return err
	}
	return nil
}

func writeLatencySamples(path string, samples <-chan latencySample) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	enc := json.NewEncoder(f)
	for sample := range samples {
		if err := enc.Encode(sample); err != nil {
			return err
		}
	}
	return nil
}

func writeE2ELatencySamples(path string, samples <-chan e2eLatencySample) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	enc := json.NewEncoder(f)
	for sample := range samples {
		if err := enc.Encode(sample); err != nil {
			return err
		}
	}
	return nil
}

func parseBenchmarkContent(content string) (clientIdx int, sequence int, sentAtUnixNano int64, ok bool) {
	fields := strings.Fields(content)
	if len(fields) != 4 {
		return 0, 0, 0, false
	}
	if fields[0] != *contentPrefix {
		return 0, 0, 0, false
	}
	if !strings.HasPrefix(fields[1], "c") || !strings.HasPrefix(fields[2], "s") || !strings.HasPrefix(fields[3], "t") {
		return 0, 0, 0, false
	}
	clientIdx64, err := strconv.ParseInt(strings.TrimPrefix(fields[1], "c"), 10, 64)
	if err != nil {
		return 0, 0, 0, false
	}
	sequence64, err := strconv.ParseInt(strings.TrimPrefix(fields[2], "s"), 10, 64)
	if err != nil {
		return 0, 0, 0, false
	}
	sentAtUnixNano, err = strconv.ParseInt(strings.TrimPrefix(fields[3], "t"), 10, 64)
	if err != nil {
		return 0, 0, 0, false
	}
	return int(clientIdx64), int(sequence64), sentAtUnixNano, true
}

func runClient(ctx context.Context, clientIdx int, token string, senderID string, activeSender bool, recipients []string, rooms []string, pairMap map[string]string, seed int64, corrStore *correlationStore, latStats *latencyStats, m *metrics, latencyCh chan<- latencySample, e2eLatencyCh chan<- e2eLatencySample) {
	d := websocket.Dialer{HandshakeTimeout: 8 * time.Second}
	h := http.Header{}
	h.Set("Authorization", normalizeAuth(token))
	if *origin != "" {
		h.Set("Origin", *origin)
	}

	conn, _, err := d.Dial(*wsURL, h)
	if err != nil {
		atomic.AddInt64(&m.connectFail, 1)
		return
	}
	defer func() { _ = conn.Close() }()
	atomic.AddInt64(&m.connected, 1)

	clientRng := rand.New(rand.NewSource(seed))
	var pendingMu sync.Mutex
	pending := make(map[string]pendingMessage)
	stopTimeoutSweep := make(chan struct{})
	if *ackTimeout > 0 {
		sweepEvery := *ackTimeout / 2
		if sweepEvery < 500*time.Millisecond {
			sweepEvery = 500 * time.Millisecond
		}
		if sweepEvery > *ackTimeout {
			sweepEvery = *ackTimeout
		}
		go func() {
			ticker := time.NewTicker(sweepEvery)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-stopTimeoutSweep:
					return
				case <-ticker.C:
					deadline := time.Now().Add(-*ackTimeout)
					pendingMu.Lock()
					timedOut := 0
					for messageID, msg := range pending {
						if msg.sentAt.After(deadline) {
							continue
						}
						delete(pending, messageID)
						timedOut++
					}
					pendingMu.Unlock()
					if timedOut > 0 {
						atomic.AddInt64(&m.ackTimedOut, int64(timedOut))
						atomic.AddInt64(&m.pendingCurrent, -int64(timedOut))
					}
				}
			}
		}()
	}
	removePending := func(messageID string) {
		pendingMu.Lock()
		if _, ok := pending[messageID]; ok {
			delete(pending, messageID)
			pendingMu.Unlock()
			atomic.AddInt64(&m.pendingCurrent, -1)
			return
		}
		pendingMu.Unlock()
	}
	resolvePending := func(messageID string) (pendingMessage, time.Duration, bool) {
		pendingMu.Lock()
		msg, ok := pending[messageID]
		if ok {
			delete(pending, messageID)
		}
		pendingMu.Unlock()
		if !ok {
			return pendingMessage{}, 0, false
		}
		atomic.AddInt64(&m.pendingCurrent, -1)
		return msg, time.Since(msg.sentAt), true
	}
	defer close(stopTimeoutSweep)

	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		handleMessage := func(data []byte) {
			var env downstreamEnvelope
			envErr := json.Unmarshal(data, &env)
			if envErr == nil && env.ClientMessageID != "" {
				if businessMessageID := extractBusinessMessageIDFromRaw(env.Payload); businessMessageID != "" {
					corrStore.bindBusinessMessageID(env.ClientMessageID, businessMessageID)
				}
			}

			var ack ackEnvelope
			if err := json.Unmarshal(data, &ack); err == nil {
				if ack.AckType != "" || ack.Error != "" {
					if ack.ClientMessageID != "" {
						if pendingMsg, latency, ok := resolvePending(ack.ClientMessageID); ok {
							atomic.AddInt64(&m.ackMatched, 1)
							atomic.AddInt64(&m.ackLatencyNs, latency.Nanoseconds())
							if latStats != nil {
								latStats.addAck(float64(latency) / float64(time.Millisecond))
							}
							for {
								currentMax := atomic.LoadInt64(&m.ackMaxLatencyNs)
								if latency.Nanoseconds() <= currentMax {
									break
								}
								if atomic.CompareAndSwapInt64(&m.ackMaxLatencyNs, currentMax, latency.Nanoseconds()) {
									break
								}
							}
							if latencyCh != nil {
								latencyCh <- latencySample{
									Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
									ClientID:  pendingMsg.clientIdx,
									Sequence:  pendingMsg.seq,
									MessageID: ack.ClientMessageID,
									AckType:   strings.ToLower(ack.AckType),
									LatencyMS: float64(latency) / float64(time.Millisecond),
									SendMode:  strings.ToLower(strings.TrimSpace(*sendMode)),
								}
							}
						} else {
							atomic.AddInt64(&m.ackUnmatched, 1)
						}
					}
					switch strings.ToLower(ack.AckType) {
					case "received":
						atomic.AddInt64(&m.ackReceived, 1)
					case "error":
						atomic.AddInt64(&m.ackError, 1)
					}
					return
				}
			}

			if envErr != nil {
				return
			}
			if strings.ToLower(env.Action) != "push_notification" {
				return
			}

			if len(env.Payload) == 0 {
				return
			}

			var payload benchmarkMessagePayload
			if err := json.Unmarshal(env.Payload, &payload); err != nil {
				return
			}

			businessMessageID := extractBusinessMessageID(payload)

			correlationMethod := ""
			clientMessageID := ""
			clientIdx0 := 0
			seq := 0
			sentAtUnixNano := int64(0)

			if businessMessageID != "" {
				if meta, ok := corrStore.resolveByBusinessMessageID(businessMessageID); ok {
					correlationMethod = "business_message_id+client_message_id"
					clientMessageID = meta.clientMessageID
					clientIdx0 = meta.clientIdx
					seq = meta.seq
					sentAtUnixNano = meta.sentAtUnixNano
				}
			}

			if sentAtUnixNano == 0 && env.ClientMessageID != "" {
				if meta, ok := corrStore.resolveByClientMessageID(env.ClientMessageID); ok {
					correlationMethod = "client_message_id"
					clientMessageID = meta.clientMessageID
					clientIdx0 = meta.clientIdx
					seq = meta.seq
					sentAtUnixNano = meta.sentAtUnixNano
				}
			}

			if sentAtUnixNano == 0 {
				parsedClientIdx, parsedSeq, parsedSentAtUnixNano, ok := parseBenchmarkContent(payload.Content)
				if !ok {
					return
				}
				correlationMethod = "content_fallback"
				clientIdx0 = parsedClientIdx
				seq = parsedSeq
				sentAtUnixNano = parsedSentAtUnixNano
			}

			recvAtUnixNano := time.Now().UTC().UnixNano()
			latencyMs := float64(recvAtUnixNano-sentAtUnixNano) / float64(time.Millisecond)
			if latencyMs < 0 {
				return
			}
			if latStats != nil {
				latStats.addE2E(latencyMs)
			}
			if e2eLatencyCh != nil {
				e2eLatencyCh <- e2eLatencySample{
					Timestamp:         time.Now().UTC().Format(time.RFC3339Nano),
					ClientID:          clientIdx0,
					Sequence:          seq,
					ClientMessageID:   clientMessageID,
					BusinessMessageID: businessMessageID,
					CorrelationMethod: correlationMethod,
					SentAtUnixNano:    sentAtUnixNano,
					RecvAtUnixNano:    recvAtUnixNano,
					LatencyMS:         latencyMs,
					Content:           payload.Content,
					Action:            env.Action,
				}
			}
		}
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				if ctx.Err() == nil {
					atomic.AddInt64(&m.readFail, 1)
				}
				return
			}
			atomic.AddInt64(&m.recv, 1)

			dec := json.NewDecoder(bytes.NewReader(data))
			for {
				var raw json.RawMessage
				if err := dec.Decode(&raw); err != nil {
					break
				}
				handleMessage(raw)
			}
		}
	}()

	if !activeSender || *sendInterval <= 0 {
		<-ctx.Done()
		return
	}
	sendTicker := time.NewTicker(*sendInterval)
	defer sendTicker.Stop()

	seq := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-readDone:
			return
		case <-sendTicker.C:
			seq++
			env, err := buildEnvelope(clientIdx, seq, senderID, recipients, rooms, pairMap, clientRng)
			if err != nil {
				atomic.AddInt64(&m.sendFail, 1)
				continue
			}
			now := time.Now()
			pendingMu.Lock()
			pending[env.ClientMessageID] = pendingMessage{sentAt: now, clientIdx: clientIdx, seq: seq}
			pendingMu.Unlock()
			corrStore.trackByClientMessageID(env.ClientMessageID, sentMessageMeta{
				clientIdx:       clientIdx,
				seq:             seq,
				sentAtUnixNano:  now.UTC().UnixNano(),
				clientMessageID: env.ClientMessageID,
			})
			atomic.AddInt64(&m.pendingCurrent, 1)
			payload, err := json.Marshal(env)
			if err != nil {
				removePending(env.ClientMessageID)
				atomic.AddInt64(&m.sendFail, 1)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				removePending(env.ClientMessageID)
				atomic.AddInt64(&m.writeFail, 1)
				return
			}
			atomic.AddInt64(&m.sent, 1)
		}
	}
}

func buildEnvelope(clientIdx, seq int, senderID string, recipients, rooms []string, pairMap map[string]string, rng *rand.Rand) (*upstreamEnvelope, error) {
	mode := strings.ToLower(strings.TrimSpace(*sendMode))
	if mode == "mixed" {
		if rng.Float64() < *roomRatio {
			mode = "room"
		} else {
			mode = "private"
		}
	}

	msgID, err := newMessageID()
	if err != nil {
		return nil, err
	}
	content := fmt.Sprintf("%s c%d s%d t%d", *contentPrefix, clientIdx, seq, time.Now().UnixNano())

	switch mode {
	case "private":
		recipient := ""
		if senderID != "" {
			recipient = pairMap[senderID]
		}
		if recipient == "" {
			if len(recipients) == 0 {
				return nil, errors.New("no recipients for private mode")
			}
			for i := 0; i < 3; i++ {
				candidate := recipients[rng.Intn(len(recipients))]
				if senderID == "" || candidate != senderID {
					recipient = candidate
					break
				}
			}
			if recipient == "" {
				recipient = recipients[rng.Intn(len(recipients))]
			}
		}

		p, _ := json.Marshal(privatePayload{RecipientID: recipient, Content: content})
		return &upstreamEnvelope{
			ClientMessageID: msgID,
			Action:          actionSendPrivateMessage,
			Payload:         p,
		}, nil
	case "room":
		if len(rooms) == 0 {
			return nil, errors.New("no rooms for room mode")
		}
		roomID := rooms[rng.Intn(len(rooms))]
		p, _ := json.Marshal(roomPayload{RoomID: roomID, Content: content})
		return &upstreamEnvelope{
			ClientMessageID: msgID,
			Action:          actionSendRoomMessage,
			Payload:         p,
		}, nil
	default:
		return nil, fmt.Errorf("unknown send-mode: %s", *sendMode)
	}
}

func newMessageID() (string, error) {
	id, err := uuid.NewV7()
	if err == nil {
		return id.String(), nil
	}
	fallback, err2 := uuid.NewRandom()
	if err2 != nil {
		return "", err2
	}
	return fallback.String(), nil
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	out := make([]string, 0, 1024)
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func readPairs(path string) (map[string]string, error) {
	m := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return m, err
	}
	defer func() { _ = f.Close() }()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		var p pairLine
		if err := json.Unmarshal([]byte(line), &p); err != nil {
			continue
		}
		if p.SenderID == "" || p.RecipientID == "" {
			continue
		}
		// only keep success-like pairs; if fields missing, still accept.
		if p.SendStatus != "" && p.SendStatus != "ok" && p.SendStatus != "skipped" {
			continue
		}
		if p.AgreeStatus != "" && p.AgreeStatus != "ok" && p.AgreeStatus != "skipped" {
			continue
		}
		m[p.SenderID] = p.RecipientID
	}
	return m, s.Err()
}

func normalizeAuth(token string) string {
	token = strings.TrimSpace(token)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		return token
	}
	return "Bearer " + token
}

func printMetrics(prefix string, m *metrics, ackAll percentileSnapshot, e2eAll percentileSnapshot, ackWindow percentileSnapshot, e2eWindow percentileSnapshot) {
	matched := atomic.LoadInt64(&m.ackMatched)
	avgLatencyNs := int64(0)
	if matched > 0 {
		avgLatencyNs = atomic.LoadInt64(&m.ackLatencyNs) / matched
	}
	log.Printf("[%s] connected=%d connect_fail=%d sent=%d send_fail=%d recv=%d ack_matched=%d ack_unmatched=%d ack_timed_out=%d pending=%d ack_received=%d ack_error=%d ack_avg_latency_ms=%.2f ack_max_latency_ms=%.2f ack_p50_ms=%.2f ack_p90_ms=%.2f ack_p95_ms=%.2f ack_p99_ms=%.2f ack_window_count=%d ack_window_p95_ms=%.2f e2e_count=%d e2e_p50_ms=%.2f e2e_p90_ms=%.2f e2e_p95_ms=%.2f e2e_p99_ms=%.2f e2e_window_count=%d e2e_window_p95_ms=%.2f read_fail=%d write_fail=%d",
		prefix,
		atomic.LoadInt64(&m.connected),
		atomic.LoadInt64(&m.connectFail),
		atomic.LoadInt64(&m.sent),
		atomic.LoadInt64(&m.sendFail),
		atomic.LoadInt64(&m.recv),
		matched,
		atomic.LoadInt64(&m.ackUnmatched),
		atomic.LoadInt64(&m.ackTimedOut),
		atomic.LoadInt64(&m.pendingCurrent),
		atomic.LoadInt64(&m.ackReceived),
		atomic.LoadInt64(&m.ackError),
		float64(avgLatencyNs)/float64(time.Millisecond),
		float64(atomic.LoadInt64(&m.ackMaxLatencyNs))/float64(time.Millisecond),
		ackAll.P50MS,
		ackAll.P90MS,
		ackAll.P95MS,
		ackAll.P99MS,
		ackWindow.Count,
		ackWindow.P95MS,
		e2eAll.Count,
		e2eAll.P50MS,
		e2eAll.P90MS,
		e2eAll.P95MS,
		e2eAll.P99MS,
		e2eWindow.Count,
		e2eWindow.P95MS,
		atomic.LoadInt64(&m.readFail),
		atomic.LoadInt64(&m.writeFail),
	)
}

func buildClientPlans(tokens, recipients []string, pairMap map[string]string) []clientPlan {
	maxN := len(tokens)
	if len(recipients) < maxN {
		maxN = len(recipients)
	}

	plans := make([]clientPlan, 0, maxN)
	for i := 0; i < maxN; i++ {
		senderID := recipients[i]
		_, activeSender := pairMap[senderID]
		plans = append(plans, clientPlan{
			ClientIdx:    i,
			Token:        tokens[i],
			SenderID:     senderID,
			ActiveSender: activeSender,
		})
	}

	activeSenders := 0
	for _, plan := range plans {
		if plan.ActiveSender {
			activeSenders++
		}
	}
	log.Printf("benchmark clients=%d active_senders=%d receivers=%d", len(plans), activeSenders, len(plans)-activeSenders)
	return plans
}
