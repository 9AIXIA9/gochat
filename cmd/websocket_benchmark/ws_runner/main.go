package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	actionSendPrivateMessage = "send_private_message_command"
	actionSendRoomMessage    = "send_room_message_command"
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

type pairLine struct {
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	SendStatus  string `json:"send_status"`
	AgreeStatus string `json:"agree_status"`
}

type metrics struct {
	connected   int64
	connectFail int64
	sent        int64
	sendFail    int64
	recv        int64
	ackReceived int64
	ackError    int64
	readFail    int64
	writeFail   int64
}

type clientPlan struct {
	ClientIdx int
	Token     string
	SenderID  string
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

	connectRate = flag.Int("connect-rate", 0, "connections per second, 0 means burst")
	origin      = flag.String("origin", "", "optional Origin header")

	contentPrefix = flag.String("content-prefix", "bench", "message content prefix")
	seed          = flag.Int64("seed", 0, "rng seed, 0 means now")
	reportEvery   = flag.Duration("report-interval", 5*time.Second, "periodic report interval, 0 disables")
)

func main() {
	flag.Parse()

	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(*seed))

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
	wg := sync.WaitGroup{}

	if *reportEvery > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t := time.NewTicker(*reportEvery)
			defer t.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					printMetrics("progress", &m)
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
		wg.Add(1)
		go func(p clientPlan) {
			defer wg.Done()
			runClient(ctx, p.ClientIdx, p.Token, p.SenderID, recipients, rooms, pairMap, rng.Int63(), &m)
		}(plan)
	}

	<-ctx.Done()
	wg.Wait()
	printMetrics("final", &m)
}

func runClient(ctx context.Context, clientIdx int, token string, senderID string, recipients []string, rooms []string, pairMap map[string]string, seed int64, m *metrics) {
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

	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				if ctx.Err() == nil {
					atomic.AddInt64(&m.readFail, 1)
				}
				return
			}
			atomic.AddInt64(&m.recv, 1)

			var ack ackEnvelope
			if err := json.Unmarshal(data, &ack); err == nil {
				switch strings.ToLower(ack.AckType) {
				case "received":
					atomic.AddInt64(&m.ackReceived, 1)
				case "error":
					atomic.AddInt64(&m.ackError, 1)
				}
			}
		}
	}()

	if *sendInterval <= 0 {
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
			payload, err := json.Marshal(env)
			if err != nil {
				atomic.AddInt64(&m.sendFail, 1)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
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

func printMetrics(prefix string, m *metrics) {
	log.Printf("[%s] connected=%d connect_fail=%d sent=%d send_fail=%d recv=%d ack_received=%d ack_error=%d read_fail=%d write_fail=%d",
		prefix,
		atomic.LoadInt64(&m.connected),
		atomic.LoadInt64(&m.connectFail),
		atomic.LoadInt64(&m.sent),
		atomic.LoadInt64(&m.sendFail),
		atomic.LoadInt64(&m.recv),
		atomic.LoadInt64(&m.ackReceived),
		atomic.LoadInt64(&m.ackError),
		atomic.LoadInt64(&m.readFail),
		atomic.LoadInt64(&m.writeFail),
	)
}

func buildClientPlans(tokens, recipients []string, pairMap map[string]string) []clientPlan {
	maxN := len(tokens)
	if len(recipients) < maxN {
		maxN = len(recipients)
	}

	mode := strings.ToLower(strings.TrimSpace(*sendMode))
	needPairBinding := (mode == "private" || mode == "mixed") && len(pairMap) > 0

	plans := make([]clientPlan, 0, maxN)
	for i := 0; i < maxN; i++ {
		senderID := recipients[i]
		if needPairBinding {
			if _, ok := pairMap[senderID]; !ok {
				continue
			}
		}
		plans = append(plans, clientPlan{
			ClientIdx: i,
			Token:     tokens[i],
			SenderID:  senderID,
		})
	}

	if needPairBinding {
		log.Printf("pair-bound mode enabled: usable senders=%d (from %d)", len(plans), maxN)
	}
	return plans
}
