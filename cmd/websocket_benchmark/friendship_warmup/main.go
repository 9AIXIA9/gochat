package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	sendFriendRequestPath = "/friendship-requests"
	agreeFriendReqFmtPath = "/friendship-requests/%s/agree"
)

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type listFriendReqResponseData struct {
	Requests []friendRequest `json:"requests"`
}

type friendRequest struct {
	ID    string `json:"id"`
	From  string `json:"from"`
	To    string `json:"to"`
	State string `json:"state"`
}

type pairTask struct {
	Index       int    `json:"index"`
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	SenderToken string `json:"-"`
	RecvToken   string `json:"-"`
}

type pairResult struct {
	Index       int    `json:"index"`
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	SendStatus  string `json:"send_status"`
	AgreeStatus string `json:"agree_status"`
	Message     string `json:"message,omitempty"`
}

type counters struct {
	planned      int64
	sendOK       int64
	sendSkipped  int64
	sendFailed   int64
	agreeOK      int64
	agreeSkipped int64
	agreeFailed  int64
}

var (
	baseURL = flag.String("base-url", "http://localhost:8080/api/v1", "Backend API Base URL")

	tokensFile     = flag.String("tokens-file", "./websocket_benchmark_data/output/tokens.txt", "Token file path")
	recipientsFile = flag.String("recipients-file", "./websocket_benchmark_data/output/recipients.txt", "Recipients(user_id) file path")
	outMapFile     = flag.String("out-map-file", "./websocket_benchmark_data/output/friendship_pairs.jsonl", "Output pair map result file")

	pairs          = flag.Int("pairs", 0, "How many sender->recipient pairs to warmup. 0 means auto = min(len(tokens), len(recipients))/2")
	workers        = flag.Int("workers", 30, "Number of concurrent workers")
	requestContent = flag.String("request-content", "benchmark warmup", "Friend request content")

	listLimit    = flag.Int("list-limit", 100, "Page size when querying /friendship-requests")
	maxPages     = flag.Int("max-pages", 20, "Max pages to scan when searching request ID")
	maxRetries   = flag.Int("max-retries", 2, "Max retries per HTTP operation")
	retryBackoff = flag.Duration("retry-backoff", 250*time.Millisecond, "Retry backoff")
	strict       = flag.Bool("strict", false, "Exit non-zero when any hard failure occurs")

	client = &http.Client{Timeout: 10 * time.Second}
)

func main() {
	flag.Parse()

	tokens, err := readLines(*tokensFile)
	if err != nil {
		log.Fatalf("failed to read tokens file: %v", err)
	}
	userIDs, err := readLines(*recipientsFile)
	if err != nil {
		log.Fatalf("failed to read recipients file: %v", err)
	}

	n := min(len(tokens), len(userIDs))
	if n < 2 {
		log.Fatalf("not enough accounts: tokens=%d recipients=%d", len(tokens), len(userIDs))
	}
	tokens = tokens[:n]
	userIDs = userIDs[:n]

	targetPairs := *pairs
	if targetPairs <= 0 {
		targetPairs = n / 2
	}
	if targetPairs <= 0 {
		log.Fatalf("calculated pairs is 0 (n=%d)", n)
	}
	if targetPairs*2 > n {
		log.Fatalf("pairs=%d requires at least %d users, but got %d", targetPairs, targetPairs*2, n)
	}

	tasks := buildPairedTasks(tokens, userIDs, targetPairs)
	results := make([]pairResult, len(tasks))
	var cnt counters
	atomic.StoreInt64(&cnt.planned, int64(len(tasks)))

	log.Printf("friendship_warmup start: users=%d pairs=%d workers=%d", n, len(tasks), *workers)

	runWorkerPool(tasks, *workers, func(t pairTask) {
		r := pairResult{Index: t.Index, SenderID: t.SenderID, RecipientID: t.RecipientID}
		err := retry(*maxRetries, *retryBackoff, func() error {
			status, message, e := sendFriendRequest(t.SenderToken, t.RecipientID, *requestContent)
			if e != nil {
				return e
			}

			if isSuccess(status) {
				r.SendStatus = "ok"
				atomic.AddInt64(&cnt.sendOK, 1)
				return nil
			}

			if isIdempotentFriendshipMessage(message) {
				r.SendStatus = "skipped"
				r.Message = message
				atomic.AddInt64(&cnt.sendSkipped, 1)
				return nil
			}

			return fmt.Errorf("send status=%d message=%s", status, message)
		})
		if err != nil {
			r.SendStatus = "failed"
			r.Message = err.Error()
			atomic.AddInt64(&cnt.sendFailed, 1)
			results[t.Index] = r
			return
		}

		if r.SendStatus == "" {
			r.SendStatus = "ok"
		}

		agreeErr := retry(*maxRetries, *retryBackoff, func() error {
			requestID, e := findPendingRequestID(t.RecvToken, t.SenderID, *listLimit, *maxPages)
			if e != nil {
				return e
			}
			if requestID == "" {
				r.AgreeStatus = "skipped"
				atomic.AddInt64(&cnt.agreeSkipped, 1)
				return nil
			}

			status, message, e := agreeFriendRequest(t.RecvToken, requestID)
			if e != nil {
				return e
			}
			if isSuccess(status) {
				r.AgreeStatus = "ok"
				atomic.AddInt64(&cnt.agreeOK, 1)
				return nil
			}
			if isIdempotentFriendshipMessage(message) {
				r.AgreeStatus = "skipped"
				r.Message = message
				atomic.AddInt64(&cnt.agreeSkipped, 1)
				return nil
			}
			return fmt.Errorf("agree status=%d message=%s", status, message)
		})

		if agreeErr != nil {
			r.AgreeStatus = "failed"
			if r.Message == "" {
				r.Message = agreeErr.Error()
			}
			atomic.AddInt64(&cnt.agreeFailed, 1)
		} else if r.AgreeStatus == "" {
			r.AgreeStatus = "ok"
		}

		results[t.Index] = r
	})

	if err := writeJSONL(*outMapFile, results); err != nil {
		log.Fatalf("failed to write pair map: %v", err)
	}

	log.Printf("friendship_warmup done: planned=%d send(ok=%d skipped=%d failed=%d) agree(ok=%d skipped=%d failed=%d) out=%s",
		cnt.planned, cnt.sendOK, cnt.sendSkipped, cnt.sendFailed, cnt.agreeOK, cnt.agreeSkipped, cnt.agreeFailed, *outMapFile)

	if *strict && (cnt.sendFailed > 0 || cnt.agreeFailed > 0) {
		os.Exit(1)
	}
}

func buildPairedTasks(tokens, userIDs []string, targetPairs int) []pairTask {
	tasks := make([]pairTask, 0, targetPairs)
	for i := 0; i < targetPairs; i++ {
		recvIdx := i
		senderIdx := i + targetPairs
		tasks = append(tasks, pairTask{
			Index:       i,
			SenderID:    userIDs[senderIdx],
			RecipientID: userIDs[recvIdx],
			SenderToken: tokens[senderIdx],
			RecvToken:   tokens[recvIdx],
		})
	}
	return tasks
}

func sendFriendRequest(token, toID, content string) (int, string, error) {
	body, _ := json.Marshal(map[string]string{
		"to_id":   toID,
		"content": content,
	})
	req, err := http.NewRequest(http.MethodPost, *baseURL+sendFriendRequestPath, bytes.NewBuffer(body))
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", normalizeAuth(token))

	status, apiResp, _, err := do(req)
	if err != nil {
		return 0, "", err
	}
	return status, strings.ToLower(apiResp.Message), nil
}

func agreeFriendRequest(token, requestID string) (int, string, error) {
	uri := fmt.Sprintf(agreeFriendReqFmtPath, url.PathEscape(requestID))
	req, err := http.NewRequest(http.MethodPut, *baseURL+uri, nil)
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Authorization", normalizeAuth(token))

	status, apiResp, _, err := do(req)
	if err != nil {
		return 0, "", err
	}
	return status, strings.ToLower(apiResp.Message), nil
}

func findPendingRequestID(token, senderID string, limit, maxPages int) (string, error) {
	baseID := ""
	for i := 0; i < maxPages; i++ {
		u := fmt.Sprintf("%s%s?limit=%d", *baseURL, sendFriendRequestPath, limit)
		if baseID != "" {
			u += "&base_id=" + url.QueryEscape(baseID)
		}

		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", normalizeAuth(token))

		status, apiResp, rawBody, err := do(req)
		if err != nil {
			return "", err
		}
		if !isSuccess(status) {
			return "", fmt.Errorf("list requests status=%d message=%s", status, apiResp.Message)
		}

		var data listFriendReqResponseData
		if err := json.Unmarshal(apiResp.Data, &data); err != nil {
			return "", fmt.Errorf("unmarshal list data failed: %w body=%s", err, string(rawBody))
		}

		if len(data.Requests) == 0 {
			return "", nil
		}

		for _, reqItem := range data.Requests {
			if reqItem.From == senderID && strings.EqualFold(reqItem.State, "pending") {
				return reqItem.ID, nil
			}
		}

		baseID = data.Requests[len(data.Requests)-1].ID
	}

	return "", nil
}

func do(req *http.Request) (int, apiResponse, []byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return 0, apiResponse{}, nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, apiResponse{}, nil, err
	}

	var apiResp apiResponse
	if len(body) > 0 {
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return resp.StatusCode, apiResponse{}, body, nil
		}
	}

	return resp.StatusCode, apiResp, body, nil
}

func isSuccess(status int) bool {
	return status == http.StatusOK
}

func normalizeAuth(token string) string {
	token = strings.TrimSpace(token)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		return token
	}
	return "Bearer " + token
}

func isIdempotentFriendshipMessage(message string) bool {
	m := strings.ToLower(message)
	keywords := []string{
		"already",
		"pending",
		"exist",
		"friend",
		"processed",
		"agreed",
	}
	for _, k := range keywords {
		if strings.Contains(m, k) {
			return true
		}
	}
	return false
}

func retry(maxRetries int, backoff time.Duration, fn func() error) error {
	var lastErr error
	for i := 0; i <= maxRetries; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err
		if i < maxRetries {
			time.Sleep(backoff)
		}
	}
	if lastErr == nil {
		lastErr = errors.New("retry failed")
	}
	return lastErr
}

func runWorkerPool(tasks []pairTask, workers int, fn func(pairTask)) {
	if workers <= 0 {
		workers = 1
	}
	ch := make(chan pairTask, len(tasks))
	for _, t := range tasks {
		ch <- t
	}
	close(ch)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range ch {
				fn(t)
			}
		}()
	}
	wg.Wait()
}

func readLines(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(string(b), "\n")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return out, nil
}

func writeJSONL(path string, rows []pairResult) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	enc := json.NewEncoder(f)
	for _, row := range rows {
		if err := enc.Encode(row); err != nil {
			return err
		}
	}
	return nil
}
