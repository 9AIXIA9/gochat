package main

import (
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
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type config struct {
	baseURL          string
	count            int
	startIndex       int
	password         string
	emailPrefix      string
	emailDomain      string
	outFile          string
	detailFile       string
	workers          int
	timeout          time.Duration
	maxRetries       int
	retryBackoff     time.Duration
	progressInterval time.Duration
	append           bool
	dryRun           bool
}

type signUpResp struct {
	Data struct {
		UserNumber string `json:"user_number"`
	} `json:"data"`
}

type loginResp struct {
	Data struct {
		AccessToken string `json:"access_token"`
	} `json:"data"`
}

type task struct {
	index int
	idx   int
}

type result struct {
	index       int
	idx         int
	email       string
	userNumber  string
	accessToken string
	err         error
}

type detailLine struct {
	Index       int    `json:"index"`
	Email       string `json:"email"`
	UserNumber  string `json:"user_number"`
	AccessToken string `json:"access_token"`
}

func main() {
	cfg := parseFlags()
	if err := validateConfig(cfg); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("get cwd failed: %v", err)
	}

	outPath := resolvePath(root, cfg.outFile)
	detailPath := resolvePath(root, cfg.detailFile)
	if err := ensureParentDirs(outPath, detailPath); err != nil {
		log.Fatalf("prepare output directories failed: %v", err)
	}

	log.Printf("start generating token pool: count=%d workers=%d dry_run=%v append=%v", cfg.count, cfg.workers, cfg.dryRun, cfg.append)
	log.Printf("base_url=%s", cfg.baseURL)

	res := run(cfg)

	if err := writeOutputs(outPath, detailPath, res, cfg.append); err != nil {
		log.Fatalf("write output failed: %v", err)
	}

	failCount := 0
	for _, r := range res {
		if r.err != nil {
			failCount++
		}
	}
	okCount := len(res) - failCount

	log.Printf("done")
	log.Printf("success: %d", okCount)
	log.Printf("failed:  %d", failCount)
	log.Printf("token file:  %s", outPath)
	log.Printf("detail file: %s", detailPath)

	if failCount > 0 {
		sampled := 0
		for _, r := range res {
			if r.err != nil {
				log.Printf("sample fail idx=%d email=%s err=%v", r.idx, r.email, r.err)
				sampled++
				if sampled >= 20 {
					break
				}
			}
		}
	}
}

func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.baseURL, "base-url", "http://localhost:8080/api/v1", "API base URL, no trailing slash")
	flag.IntVar(&cfg.count, "count", 100, "number of accounts/tokens to generate")
	flag.IntVar(&cfg.startIndex, "start-index", 1, "start index for generated account sequence")
	flag.StringVar(&cfg.password, "password", "123456", "password used for sign-up/login")
	flag.StringVar(&cfg.emailPrefix, "email-prefix", "bench", "email prefix")
	flag.StringVar(&cfg.emailDomain, "email-domain", "example.com", "email domain")
	flag.StringVar(&cfg.outFile, "out-file", "scripts/websocket_benchmark/tokens.txt", "output token file path")
	flag.StringVar(&cfg.detailFile, "detail-file", "scripts/websocket_benchmark/token_pool.jsonl", "output detail jsonl path")
	flag.BoolVar(&cfg.append, "append", true, "append results to output files; set -append=false to overwrite")
	flag.IntVar(&cfg.workers, "workers", 200, "number of concurrent workers")
	flag.DurationVar(&cfg.timeout, "timeout", 12*time.Second, "single HTTP request timeout")
	flag.IntVar(&cfg.maxRetries, "max-retries", 2, "max retries per account on failure")
	flag.DurationVar(&cfg.retryBackoff, "retry-backoff", 250*time.Millisecond, "base backoff between retries")
	flag.DurationVar(&cfg.progressInterval, "progress-interval", 2*time.Second, "progress print interval")
	flag.BoolVar(&cfg.dryRun, "dry-run", false, "generate fake tokens without calling server")
	flag.Parse()

	cfg.baseURL = strings.TrimRight(strings.TrimSpace(cfg.baseURL), "/")
	return cfg
}

func validateConfig(cfg config) error {
	if cfg.baseURL == "" {
		return errors.New("base-url is required")
	}
	if cfg.count <= 0 {
		return errors.New("count must be > 0")
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
	if cfg.progressInterval <= 0 {
		return errors.New("progress-interval must be > 0")
	}
	if strings.TrimSpace(cfg.password) == "" {
		return errors.New("password is required")
	}
	if strings.TrimSpace(cfg.emailPrefix) == "" {
		return errors.New("email-prefix is required")
	}
	if strings.TrimSpace(cfg.emailDomain) == "" {
		return errors.New("email-domain is required")
	}
	return nil
}

func run(cfg config) []result {
	results := make([]result, cfg.count)
	var done atomic.Int64
	var ok atomic.Int64
	var fail atomic.Int64

	tasks := make(chan task, cfg.workers*2)
	resCh := make(chan result, cfg.workers*2)

	transport := &http.Transport{
		MaxIdleConns:        cfg.workers * 4,
		MaxIdleConnsPerHost: cfg.workers * 2,
		MaxConnsPerHost:     cfg.workers * 2,
		IdleConnTimeout:     90 * time.Second,
	}
	client := &http.Client{Transport: transport, Timeout: cfg.timeout}

	var wg sync.WaitGroup
	for w := 0; w < cfg.workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for t := range tasks {
				r := processOne(cfg, client, t)
				resCh <- r
			}
		}(w)
	}

	go func() {
		for i := 0; i < cfg.count; i++ {
			tasks <- task{index: i, idx: cfg.startIndex + i}
		}
		close(tasks)
	}()

	progressDone := make(chan struct{})
	go func() {
		ticker := time.NewTicker(cfg.progressInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				finished := done.Load()
				pct := float64(finished) / float64(cfg.count) * 100
				log.Printf("progress finished=%d/%d (%.2f%%) ok=%d fail=%d", finished, cfg.count, pct, ok.Load(), fail.Load())
			case <-progressDone:
				return
			}
		}
	}()

	for i := 0; i < cfg.count; i++ {
		r := <-resCh
		results[r.index] = r
		done.Add(1)
		if r.err != nil {
			fail.Add(1)
		} else {
			ok.Add(1)
		}
	}
	close(progressDone)

	wg.Wait()
	close(resCh)
	transport.CloseIdleConnections()

	return results
}

func processOne(cfg config, client *http.Client, t task) result {
	if cfg.dryRun {
		token := fmt.Sprintf("dry-token-%d", t.idx)
		email := fmt.Sprintf("%s.dryrun.%d@%s", cfg.emailPrefix, t.idx, cfg.emailDomain)
		return result{
			index:       t.index,
			idx:         t.idx,
			email:       email,
			userNumber:  fmt.Sprintf("dryrun-%d", t.idx),
			accessToken: token,
		}
	}

	var lastErr error
	for attempt := 0; attempt <= cfg.maxRetries; attempt++ {
		email := buildEmail(cfg.emailPrefix, cfg.emailDomain, t.idx, attempt)
		userNumber, err := signUp(client, cfg.baseURL, email, cfg.password)
		if err != nil {
			lastErr = fmt.Errorf("sign-up failed: %w", err)
			sleepBackoff(cfg.retryBackoff, attempt)
			continue
		}

		accessToken, err := login(client, cfg.baseURL, userNumber, cfg.password)
		if err != nil {
			lastErr = fmt.Errorf("login failed: %w", err)
			sleepBackoff(cfg.retryBackoff, attempt)
			continue
		}

		return result{
			index:       t.index,
			idx:         t.idx,
			email:       email,
			userNumber:  userNumber,
			accessToken: accessToken,
		}
	}

	return result{
		index: t.index,
		idx:   t.idx,
		email: buildEmail(cfg.emailPrefix, cfg.emailDomain, t.idx, cfg.maxRetries),
		err:   fmt.Errorf("all retries failed: %w", lastErr),
	}
}

func buildEmail(prefix, domain string, idx, attempt int) string {
	stamp := time.Now().UTC().Format("20060102150405")
	if attempt == 0 {
		return fmt.Sprintf("%s.%s.%d@%s", prefix, stamp, idx, domain)
	}
	return fmt.Sprintf("%s.%s.%d.r%d@%s", prefix, stamp, idx, attempt, domain)
}

func signUp(client *http.Client, baseURL, email, password string) (string, error) {
	body := map[string]string{
		"email":    email,
		"password": password,
	}
	var out signUpResp
	if err := postJSON(client, baseURL+"/auth/sign-up", body, &out); err != nil {
		return "", err
	}
	if strings.TrimSpace(out.Data.UserNumber) == "" {
		return "", errors.New("user_number is empty")
	}
	return out.Data.UserNumber, nil
}

func login(client *http.Client, baseURL, userNumber, password string) (string, error) {
	body := map[string]string{
		"number":   userNumber,
		"password": password,
	}
	var out loginResp
	if err := postJSON(client, baseURL+"/auth/login", body, &out); err != nil {
		return "", err
	}
	if strings.TrimSpace(out.Data.AccessToken) == "" {
		return "", errors.New("access_token is empty")
	}
	return out.Data.AccessToken, nil
}

func postJSON(client *http.Client, url string, payload any, out any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}

	if err := json.Unmarshal(bodyBytes, out); err != nil {
		return fmt.Errorf("decode response failed: %w, body=%s", err, strings.TrimSpace(string(bodyBytes)))
	}
	return nil
}

func sleepBackoff(base time.Duration, attempt int) {
	if base <= 0 {
		return
	}
	factor := 1 << attempt
	time.Sleep(time.Duration(factor) * base)
}

func writeOutputs(tokenPath, detailPath string, results []result, appendMode bool) error {
	tokenFlags := os.O_CREATE | os.O_WRONLY
	detailFlags := os.O_CREATE | os.O_WRONLY
	if appendMode {
		tokenFlags |= os.O_APPEND
		detailFlags |= os.O_APPEND
	} else {
		tokenFlags |= os.O_TRUNC
		detailFlags |= os.O_TRUNC
	}

	tokenFile, err := os.OpenFile(tokenPath, tokenFlags, 0o644)
	if err != nil {
		return err
	}
	defer func() { _ = tokenFile.Close() }()

	detailFile, err := os.OpenFile(detailPath, detailFlags, 0o644)
	if err != nil {
		return err
	}
	defer func() { _ = detailFile.Close() }()

	enc := json.NewEncoder(detailFile)
	for _, r := range results {
		if r.err != nil {
			continue
		}
		if _, err := tokenFile.WriteString("Bearer " + r.accessToken + "\n"); err != nil {
			return err
		}
		if err := enc.Encode(detailLine{
			Index:       r.idx,
			Email:       r.email,
			UserNumber:  r.userNumber,
			AccessToken: r.accessToken,
		}); err != nil {
			return err
		}
	}
	return nil
}

func resolvePath(root, p string) string {
	if filepath.IsAbs(p) {
		return filepath.Clean(p)
	}
	return filepath.Clean(filepath.Join(root, p))
}

func ensureParentDirs(paths ...string) error {
	for _, p := range paths {
		d := filepath.Dir(p)
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	return nil
}
