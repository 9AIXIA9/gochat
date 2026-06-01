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
const loginURL = "/auth/login/user_number"
const profileMeURL = "/profiles/me"

type Account struct {
	UserID     string `json:"user_id"`
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
	baseURL              = flag.String("base-url", "http://localhost:8080/api/v1", "Backend API Base URL")
	targetCount          = flag.Int("count", 100, "Target number of accounts to maintain")
	workers              = flag.Int("workers", 50, "Number of concurrent workers")
	stateFile            = flag.String("state-file", "./websocket_benchmark_data/state/accounts.json", "Account state file")
	tokensFile           = flag.String("tokens-file", "./websocket_benchmark_data/output/tokens.txt", "Output tokens file")
	recipientsFile       = flag.String("recipients-file", "./websocket_benchmark_data/output/recipients.txt", "Output recipients(user_id) file")
	profileMaxRetries    = flag.Int("profile-max-retries", 8, "Max retries for /profiles/me when user profile is eventually consistent")
	profileRetryInterval = flag.Duration("profile-retry-interval", 300*time.Millisecond, "Retry interval for /profiles/me")
	clients              = &http.Client{Timeout: 10 * time.Second}
)

type LoginResult struct {
	UserNumber string
	Token      string
	UserID     string
}

type getMyProfileResponseData struct {
	Profile struct {
		ID string `json:"id"`
	} `json:"profile"`
}

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
	if err := os.MkdirAll(filepath.Dir(*recipientsFile), 0755); err != nil {
		log.Fatalf("Failed to create recipients directory: %v", err)
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

	// Login and extract token + user_id(profile)
	results := loginAccounts(accountsToLogin)
	if len(results) == 0 {
		log.Fatalf("No login result produced. Please check /auth/login and /profiles/me availability")
	}

	// Merge user_id back into state by user_number for reusability
	userIDByNumber := make(map[string]string, len(results))
	tokens := make([]string, 0, len(results))
	recipients := make([]string, 0, len(results))
	for _, r := range results {
		userIDByNumber[r.UserNumber] = r.UserID
		tokens = append(tokens, r.Token)
		recipients = append(recipients, r.UserID)
	}
	for i := range state.Accounts {
		if uid, ok := userIDByNumber[state.Accounts[i].UserNumber]; ok {
			state.Accounts[i].UserID = uid
		}
	}
	saveState(*stateFile, state)

	// Save tokens
	err := writeLines(*tokensFile, tokens)
	if err != nil {
		log.Fatalf("Failed to write tokens: %v", err)
	}
	err = writeLines(*recipientsFile, recipients)
	if err != nil {
		log.Fatalf("Failed to write recipients: %v", err)
	}
	log.Printf("Successfully saved %d tokens to %s", len(tokens), *tokensFile)
	log.Printf("Successfully saved %d recipients(user_id) to %s", len(recipients), *recipientsFile)
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

func loginAccounts(accounts []Account) []LoginResult {
	var mu sync.Mutex
	var results []LoginResult
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

				resp, err := clients.Post(*baseURL+loginURL, "application/json", bytes.NewBuffer(body))
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
						userID, profileErr := fetchMyUserIDWithRetry(token)
						if profileErr != nil {
							log.Printf("Get /profiles/me failed after retries for user_number=%s: %v", acc.UserNumber, profileErr)
							continue
						}

						mu.Lock()
						results = append(results, LoginResult{
							UserNumber: acc.UserNumber,
							Token:      token,
							UserID:     userID,
						})
						mu.Unlock()
					} else {
						log.Printf("No access_token in login response for user_number=%s", acc.UserNumber)
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
	return results
}

func fetchMyUserIDWithRetry(token string) (string, error) {
	var lastErr error
	for i := 0; i < *profileMaxRetries; i++ {
		userID, err := fetchMyUserID(token)
		if err == nil && userID != "" {
			return userID, nil
		}
		lastErr = err
		time.Sleep(*profileRetryInterval)
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("user id is empty")
	}
	return "", lastErr
}

func fetchMyUserID(token string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, *baseURL+profileMeURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := clients.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", err
	}

	var data getMyProfileResponseData
	if err := json.Unmarshal(apiResp.Data, &data); err != nil {
		return "", err
	}

	if data.Profile.ID == "" {
		return "", fmt.Errorf("profile.id is empty")
	}

	return data.Profile.ID, nil
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
