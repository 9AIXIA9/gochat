package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const signUpURL = "/auth/sign-up"

type Account struct {
	UserNumber string `json:"user_number"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type GlobalState struct {
	Accounts []Account `json:"accounts"`
}

type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

var (
	baseURL     = flag.String("base-url", "http://localhost:8080/api/v1", "Backend API Base URL")
	targetCount = flag.Int("count", 100, "Target number of accounts to maintain")
	workers     = flag.Int("workers", 50, "Number of concurrent workers")
	stateFile   = flag.String("state-file", "./websocket_benchmark_data/state/accounts.json", "Account state file")
	tokensFile  = flag.String("tokens-file", "./websocket_benchmark_data/output/tokens.txt", "Output tokens file")
	clients     = &http.Client{Timeout: 10 * time.Second}
)

func main() {
	flag.Parse()

	// Ensure directories exist
	if err := os.MkdirAll(filepath.Dir(*stateFile), 0755); err != nil {
		log.Fatalf("Failed to create state directory: %v", err)
		return
	}
	if err := os.MkdirAll(filepath.Dir(*tokensFile), 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
		return
	}

	state := loadState(*stateFile)
	log.Printf("Loaded %d existing accounts.", len(state.Accounts))

	// Register up to targetCount
	if len(state.Accounts) < *targetCount {
		needed := *targetCount - len(state.Accounts)
		log.Printf("Need to register %d new accounts...", needed)
		newAccounts := registerAccounts(needed, len(state.Accounts))
		state.Accounts = append(state.Accounts, newAccounts...)
		saveState(*stateFile, state)
	}

	log.Printf("Current accounts pool size: %d. Starting bulk login...", len(state.Accounts))

	// Safely slice the accounts
	accountsToLogin := state.Accounts
	if len(accountsToLogin) > *targetCount {
		accountsToLogin = accountsToLogin[:*targetCount]
	}

	if len(accountsToLogin) == 0 {
		log.Fatalf("No accounts available to login! Did registration fail?")
	}

	// Login and extract tokens
	tokens := loginAccounts(accountsToLogin)

	// Save tokens
	err := writeLines(*tokensFile, tokens)
	if err != nil {
		log.Fatalf("Failed to write tokens: %v", err)
	}
	log.Printf("Successfully saved %d tokens to %s", len(tokens), *tokensFile)
}

func loadState(path string) GlobalState {
	var state GlobalState
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return GlobalState{Accounts: []Account{}}
		}
		log.Fatalf("Failed to read state file: %v", err)
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return GlobalState{}
	}
	return state
}

func saveState(path string, state GlobalState) {
	data, _ := json.MarshalIndent(state, "", "  ")
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatalf("Failed to save state file: %v", err)
	}
}

func registerAccounts(count, startIndex int) []Account {
	var mu sync.Mutex
	var newAccs []Account
	workCh := make(chan int, count)

	for i := 0; i < count; i++ {
		workCh <- startIndex + i
	}
	close(workCh)

	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range workCh {
				acc := Account{
					Email:    fmt.Sprintf("bench_%06d_%d@gochat.com", idx, time.Now().Unix()),
					Password: "password123",
				}

				payload := map[string]string{
					"email":    acc.Email,
					"password": acc.Password,
				}
				body, _ := json.Marshal(payload)

				resp, err := clients.Post(*baseURL+signUpURL, "application/json", bytes.NewBuffer(body))
				if err != nil {
					log.Printf("Register request failed (Network error): %v", err)
					continue
				}

				if resp.StatusCode == 200 {
					var apiResp APIResponse
					if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
						log.Printf("Failed to decode register response: %v", err)
						return
					}
					if err := resp.Body.Close(); err != nil {
						log.Printf("Failed to close response body: %v", err)
						return
					}

					var data map[string]string
					if err := json.Unmarshal(apiResp.Data, &data); err != nil {
						log.Printf("Failed to unmarshal data field: %v", err)
						return
					}

					if num, ok := data["user_number"]; ok {
						acc.UserNumber = num
						mu.Lock()
						newAccs = append(newAccs, acc)
						mu.Unlock()
					} else {
						log.Printf("No user_number in response. Body: %v", apiResp.Data)
					}
				} else {
					respBody, _ := io.ReadAll(resp.Body)
					log.Printf("Register failed with status %d: %s", resp.StatusCode, string(respBody))
					if err := resp.Body.Close(); err != nil {
						log.Printf("Failed to close response body: %v", err)
						return
					}
				}
			}
		}()
	}
	wg.Wait()
	return newAccs
}

func loginAccounts(accounts []Account) []string {
	var mu sync.Mutex
	var tokens []string
	workCh := make(chan Account, len(accounts))

	for _, acc := range accounts {
		workCh <- acc
	}
	close(workCh)

	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for acc := range workCh {
				payload := map[string]string{
					"number":   acc.UserNumber,
					"password": acc.Password,
				}
				body, _ := json.Marshal(payload)

				resp, err := clients.Post(*baseURL+"/auth/login", "application/json", bytes.NewBuffer(body))
				if err != nil {
					log.Printf("Login request failed (Network error): %v", err)
					continue
				}

				if resp.StatusCode == 200 {
					var apiResp APIResponse
					if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
						log.Printf("Failed to decode login response: %v", err)
						return
					}
					if err := resp.Body.Close(); err != nil {
						log.Printf("Failed to close response body: %v", err)
						return
					}

					var data map[string]string
					if err := json.Unmarshal(apiResp.Data, &data); err != nil {
						log.Printf("Failed to unmarshal data field: %v", err)
						return
					}

					if token, ok := data["access_token"]; ok {
						mu.Lock()
						tokens = append(tokens, token)
						mu.Unlock()
					}
				} else {
					respBody, _ := io.ReadAll(resp.Body)
					log.Printf("Login failed with status %d: %s", resp.StatusCode, string(respBody))
					if err := resp.Body.Close(); err != nil {
						log.Printf("Failed to close response body: %v", err)
						return
					}
				}
			}
		}()
	}
	wg.Wait()
	return tokens
}

func writeLines(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
			return
		}
	}()

	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}
