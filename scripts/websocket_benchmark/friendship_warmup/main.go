package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type config struct {
	baseURL        string
	tokensFile     string
	recipientsFile string
	pairs          int
	workers        int
	timeout        time.Duration
	content        string
	maxRetries     int
	retryBackoff   time.Duration
}

type pairTask struct {
	index         int
	senderToken   string
	receiverToken string
	senderID      string
	receiverID    string
}

type pairResult struct {
	index  int
	status string
	detail string
	err    error
}

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type listFriendshipsData struct {
	Friendships []friendshipDTO `json:"friendships"`
}

type friendshipDTO struct {
	UserID1 string `json:"user_id_1"`
	UserID2 string `json:"user_id_2"`
}

type listFriendRequestsData struct {
	Requests []friendRequestDTO `json:"requests"`
}

type friendRequestDTO struct {
	ID    string `json:"id"`
	From  string `json:"from"`
	To    string `json:"to"`
	State string `json:"state"`
}

func main() {
	cfg := parseFlags()
	if err := validate(cfg); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	tokens, err := loadLines(cfg.tokensFile, true)
	if err != nil {
		log.Fatalf("load tokens failed: %v", err)
	}
	recipients, err := loadLines(cfg.recipientsFile, false)
	if err != nil {
		log.Fatalf("load recipients failed: %v", err)
	}

	pairCount, err := resolvePairCount(cfg, len(tokens), len(recipients))
	if err != nil {
		log.Fatalf("resolve pairs failed: %v", err)
	}

	log.Printf("friendship warmup start base_url=%s pairs=%d workers=%d", cfg.baseURL, pairCount, cfg.workers)

	tasks := make([]pairTask, 0, pairCount)
	for i := 0; i < pairCount; i++ {
		tasks = append(tasks, pairTask{
			index:         i,
			receiverToken: tokens[i],
			senderToken:   tokens[pairCount+i],
			receiverID:    recipients[i],
			senderID:      recipients[pairCount+i],
		})
	}

	transport := &http.Transport{
		MaxIdleConns:        cfg.workers * 4,
		MaxIdleConnsPerHost: cfg.workers * 2,
		MaxConnsPerHost:     cfg.workers * 2,
		IdleConnTimeout:     60 * time.Second,
	}
	client := &http.Client{Transport: transport, Timeout: cfg.timeout}
	defer transport.CloseIdleConnections()

	results := runWorkers(cfg, client, tasks)

	var warmed, already, failed int64
	for _, r := range results {
		switch r.status {
		case "warmed":
			warmed++
		case "already":
			already++
		default:
			failed++
		}
	}

	log.Printf("friendship warmup complete")
	log.Printf("pairs total:   %d", len(results))
	log.Printf("warmed:        %d", warmed)
	log.Printf("already_friend:%d", already)
	log.Printf("failed:        %d", failed)

	if failed > 0 {
		sampled := 0
		for _, r := range results {
			if r.status == "failed" {
				log.Printf("sample failed pair=%d detail=%s err=%v", r.index, r.detail, r.err)
				sampled++
				if sampled >= 20 {
					break
				}
			}
		}
		os.Exit(1)
	}
}

func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.baseURL, "base-url", "http://localhost:8080/api/v1", "API base URL, no trailing slash")
	flag.StringVar(&cfg.tokensFile, "tokens-file", "scripts/websocket_benchmark/tokens.txt", "tokens file path")
	flag.StringVar(&cfg.recipientsFile, "recipients-file", "scripts/websocket_benchmark/recipients.txt", "recipient userID file path")
	flag.IntVar(&cfg.pairs, "pairs", 0, "pair count, default min(len(tokens), len(recipients))/2")
	flag.IntVar(&cfg.workers, "workers", 20, "concurrent workers")
	flag.DurationVar(&cfg.timeout, "timeout", 10*time.Second, "single HTTP request timeout")
	flag.StringVar(&cfg.content, "content", "warmup friend request", "friend request content")
	flag.IntVar(&cfg.maxRetries, "max-retries", 1, "max retries per pair")
	flag.DurationVar(&cfg.retryBackoff, "retry-backoff", 250*time.Millisecond, "retry backoff base duration")
	flag.Parse()
	cfg.baseURL = strings.TrimRight(strings.TrimSpace(cfg.baseURL), "/")
	return cfg
}

func validate(cfg config) error {
	if cfg.baseURL == "" {
		return errors.New("base-url is required")
	}
	if strings.TrimSpace(cfg.tokensFile) == "" {
		return errors.New("tokens-file is required")
	}
	if strings.TrimSpace(cfg.recipientsFile) == "" {
		return errors.New("recipients-file is required")
	}
	if cfg.workers <= 0 {
		return errors.New("workers must be > 0")
	}
	if cfg.timeout <= 0 {
		return errors.New("timeout must be > 0")
	}
	if cfg.maxRetries < 0 {
		return errors.New("max-retries must be >= 0")
	}
	if cfg.retryBackoff < 0 {
		return errors.New("retry-backoff must be >= 0")
	}
	return nil
}

func resolvePairCount(cfg config, tokenCount, recipientCount int) (int, error) {
	if tokenCount < 2 || recipientCount < 2 {
		return 0, fmt.Errorf("need at least 2 tokens and 2 recipients")
	}
	available := min(tokenCount, recipientCount) / 2
	if available <= 0 {
		return 0, fmt.Errorf("insufficient paired data: tokens=%d recipients=%d", tokenCount, recipientCount)
	}
	if cfg.pairs == 0 {
		return available, nil
	}
	if cfg.pairs < 0 {
		return 0, fmt.Errorf("pairs must be >= 0")
	}
	if cfg.pairs > available {
		return 0, fmt.Errorf("pairs=%d exceeds available=%d from tokens/recipients", cfg.pairs, available)
	}
	return cfg.pairs, nil
}

func runWorkers(cfg config, client *http.Client, tasks []pairTask) []pairResult {
	results := make([]pairResult, len(tasks))
	jobs := make(chan pairTask, len(tasks))
	var finished atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < cfg.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range jobs {
				results[t.index] = processPair(cfg, client, t)
				finished.Add(1)
			}
		}()
	}

	for _, t := range tasks {
		jobs <- t
	}
	close(jobs)
	wg.Wait()
	return results
}

func processPair(cfg config, client *http.Client, p pairTask) pairResult {
	var lastErr error
	for attempt := 0; attempt <= cfg.maxRetries; attempt++ {
		if already, err := isFriend(client, cfg.baseURL, p.senderToken, p.receiverID); err == nil && already {
			return pairResult{index: p.index, status: "already", detail: "already friends"}
		}

		_ = sendFriendRequest(client, cfg.baseURL, p.senderToken, p.receiverID, cfg.content)
		requestID, err := findPendingRequestID(client, cfg.baseURL, p.receiverToken, p.senderID, p.receiverID)
		if err != nil {
			lastErr = err
			sleepBackoff(cfg.retryBackoff, attempt)
			continue
		}
		if requestID != "" {
			if err := agreeFriendRequest(client, cfg.baseURL, p.receiverToken, requestID); err != nil {
				lastErr = err
				sleepBackoff(cfg.retryBackoff, attempt)
				continue
			}
		}

		if ok, err := isFriend(client, cfg.baseURL, p.senderToken, p.receiverID); err == nil && ok {
			return pairResult{index: p.index, status: "warmed", detail: "friendship ready"}
		}
		lastErr = fmt.Errorf("friendship not ready after agree")
		sleepBackoff(cfg.retryBackoff, attempt)
	}
	return pairResult{index: p.index, status: "failed", detail: "warmup failed", err: lastErr}
}

func isFriend(client *http.Client, baseURL, token, targetUserID string) (bool, error) {
	resp, _, err := doRequestJSON(client, http.MethodGet, baseURL+"/friendships?limit=100", token, nil)
	if err != nil {
		return false, err
	}
	if resp.Code != 200 {
		return false, nil
	}
	var data listFriendshipsData
	if len(resp.Data) == 0 {
		return false, nil
	}
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return false, err
	}
	for _, f := range data.Friendships {
		if f.UserID1 == targetUserID || f.UserID2 == targetUserID {
			return true, nil
		}
	}
	return false, nil
}

func sendFriendRequest(client *http.Client, baseURL, senderToken, receiverID, content string) error {
	payload := map[string]string{
		"to_id":   receiverID,
		"content": content,
	}
	resp, _, err := doRequestJSON(client, http.MethodPost, baseURL+"/friendship-requests", senderToken, payload)
	if err != nil {
		return err
	}
	if resp.Code == 200 || resp.Code == 300 {
		return nil
	}
	return fmt.Errorf("send friend request response code=%d message=%s", resp.Code, resp.Message)
}

func findPendingRequestID(client *http.Client, baseURL, receiverToken, senderID, receiverID string) (string, error) {
	resp, _, err := doRequestJSON(client, http.MethodGet, baseURL+"/friendship-requests?limit=100", receiverToken, nil)
	if err != nil {
		return "", err
	}
	if resp.Code != 200 {
		return "", fmt.Errorf("list friend requests response code=%d message=%s", resp.Code, resp.Message)
	}
	var data listFriendRequestsData
	if len(resp.Data) == 0 {
		return "", nil
	}
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}
	for _, req := range data.Requests {
		if req.From == senderID && req.To == receiverID && strings.EqualFold(req.State, "pending") {
			return req.ID, nil
		}
	}
	return "", nil
}

func agreeFriendRequest(client *http.Client, baseURL, receiverToken, requestID string) error {
	resp, _, err := doRequestJSON(client, http.MethodPut, baseURL+"/friendship-requests/"+requestID+"/agree", receiverToken, map[string]any{})
	if err != nil {
		return err
	}
	if resp.Code == 200 || resp.Code == 300 {
		return nil
	}
	return fmt.Errorf("agree friend request response code=%d message=%s", resp.Code, resp.Message)
}

func doRequestJSON(client *http.Client, method, url, bearerToken string, body any) (*apiResponse, int, error) {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, url, reader)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(bearerToken) != "" {
		req.Header.Set("Authorization", normalizeBearer(bearerToken))
	}

	httpResp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = httpResp.Body.Close() }()

	bodyBytes, err := io.ReadAll(io.LimitReader(httpResp.Body, 128*1024))
	if err != nil {
		return nil, httpResp.StatusCode, err
	}

	var out apiResponse
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &out); err != nil {
			return nil, httpResp.StatusCode, fmt.Errorf("decode response failed: %w body=%s", err, strings.TrimSpace(string(bodyBytes)))
		}
	}
	return &out, httpResp.StatusCode, nil
}

func loadLines(path string, normalizeToken bool) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if normalizeToken {
			line = normalizeBearer(line)
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("file is empty: %s", path)
	}
	return lines, nil
}

func normalizeBearer(token string) string {
	trimmed := strings.TrimSpace(token)
	if strings.HasPrefix(strings.ToLower(trimmed), "bearer ") {
		return "Bearer " + strings.TrimSpace(trimmed[len("bearer "):])
	}
	return "Bearer " + trimmed
}

func sleepBackoff(base time.Duration, attempt int) {
	if base <= 0 {
		return
	}
	factor := 1 << attempt
	time.Sleep(time.Duration(factor) * base)
}
