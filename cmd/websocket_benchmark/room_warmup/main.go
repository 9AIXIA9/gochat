package main

import (
	"bytes"
	"encoding/json"
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
	createRoomPath         = "/rooms"
	sendMemberRequestPath  = "/rooms/requests"
	listRoomshipsPath      = "/rooms"
	listMemberRequestsPath = "/rooms/requests"
)

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type roomship struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type roomshipsResponseData struct {
	Roomships []roomship `json:"roomships"`
}

type memberRequest struct {
	ID          string `json:"id"`
	State       string `json:"state"`
	ApplicantID string `json:"applicant_id"`
	RoomID      string `json:"room_id"`
	Content     string `json:"content"`
	OperatorID  string `json:"operator_id"`
}

type memberRequestsResponseData struct {
	Requests []memberRequest `json:"requests"`
}

type roomTask struct {
	Index        int      `json:"index"`
	OwnerID      string   `json:"owner_id"`
	RoomID       string   `json:"room_id"`
	MemberIDs    []string `json:"member_ids"`
	CreateStatus string   `json:"create_status"`
	InviteStatus string   `json:"invite_status"`
	Message      string   `json:"message,omitempty"`
	RequestCount int      `json:"request_count,omitempty"`
	AgreedCount  int      `json:"agreed_count,omitempty"`
	SkippedCount int      `json:"skipped_count,omitempty"`
	FailedCount  int      `json:"failed_count,omitempty"`
}

type counters struct {
	planned     int64
	createdOK   int64
	createdSkip int64
	createdFail int64
	inviteOK    int64
	inviteSkip  int64
	inviteFail  int64
	agreeOK     int64
	agreeSkip   int64
	agreeFail   int64
}

var (
	baseURL = flag.String("base-url", "http://localhost:8080/api/v1", "Backend API Base URL")

	tokensFile      = flag.String("tokens-file", "./websocket_benchmark_data/output/tokens.txt", "Token file path")
	recipientsFile  = flag.String("recipients-file", "./websocket_benchmark_data/output/recipients.txt", "Recipients(user_id) file path")
	roomsOutFile    = flag.String("rooms-file", "./websocket_benchmark_data/output/rooms.txt", "Room IDs output file")
	topologyOutFile = flag.String("out-map-file", "./websocket_benchmark_data/output/room_topology.jsonl", "Room topology output file")

	roomsCount      = flag.Int("rooms", 0, "How many rooms to create. 0 means auto by users / room-member-count")
	roomMemberCount = flag.Int("room-member-count", 20, "Total users per room including owner")
	workers         = flag.Int("workers", 10, "Concurrent room workers")
	roomPassword    = flag.String("room-password", "", "Room password. Empty means no password")
	joinContent     = flag.String("join-content", "benchmark warmup join request", "Content of member request")

	listLimit = flag.Int("list-limit", 100, "Page size when querying /rooms and /rooms/requests")
	maxPages  = flag.Int("max-pages", 20, "Max pages to scan when searching a room or member request")
	strict    = flag.Bool("strict", false, "Exit non-zero when any hard failure occurs")

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
		log.Fatalf("not enough users: tokens=%d recipients=%d", len(tokens), len(userIDs))
	}
	tokens = tokens[:n]
	userIDs = userIDs[:n]

	roomSize := *roomMemberCount
	if roomSize < 2 {
		roomSize = 2
	}
	roomCount := *roomsCount
	if roomCount <= 0 {
		roomCount = n / roomSize
	}
	if roomCount <= 0 {
		log.Fatalf("calculated room count is 0: users=%d room-member-count=%d", n, roomSize)
	}
	if roomCount*roomSize > n {
		log.Fatalf("rooms=%d with room-member-count=%d requires at least %d users, but got %d", roomCount, roomSize, roomCount*roomSize, n)
	}

	tasks := buildRoomTasks(tokens, userIDs, roomCount, roomSize)
	results := make([]roomTask, len(tasks))
	var cnt counters
	atomic.StoreInt64(&cnt.planned, int64(len(tasks)))

	log.Printf("room_warmup start: users=%d rooms=%d room-member-count=%d workers=%d", n, len(tasks), roomSize, *workers)

	runWorkerPool(tasks, *workers, func(task roomBuildTask) {
		result := roomTask{Index: task.Index, OwnerID: task.OwnerID, MemberIDs: task.MemberIDs}

		roomID, createStatus, createMsg, err := createRoom(task.OwnerToken, roomSize, *roomPassword)
		if err != nil {
			result.CreateStatus = "failed"
			result.Message = err.Error()
			atomic.AddInt64(&cnt.createdFail, 1)
			results[task.Index] = result
			return
		}
		if createStatus == "skipped" {
			result.CreateStatus = "skipped"
			result.Message = createMsg
			atomic.AddInt64(&cnt.createdSkip, 1)
		} else {
			result.CreateStatus = "ok"
			atomic.AddInt64(&cnt.createdOK, 1)
		}

		if roomID == "" {
			roomID, err = findLatestRoomIDByOwner(task.OwnerToken, *listLimit, *maxPages)
			if err != nil {
				result.CreateStatus = "failed"
				result.Message = err.Error()
				atomic.AddInt64(&cnt.createdFail, 1)
				results[task.Index] = result
				return
			}
		}
		result.RoomID = roomID

		inviteOK, inviteSkip, inviteFail, agreeOK, agreeSkip, agreeFail, inviteMsg := warmupRoomMembers(task, roomID)
		result.InviteStatus = "ok"
		result.RequestCount = inviteOK + inviteSkip + inviteFail
		result.AgreedCount = agreeOK
		result.SkippedCount = inviteSkip + agreeSkip
		result.FailedCount = inviteFail + agreeFail
		if inviteMsg != "" {
			result.Message = inviteMsg
		}

		atomic.AddInt64(&cnt.inviteOK, int64(inviteOK))
		atomic.AddInt64(&cnt.inviteSkip, int64(inviteSkip))
		atomic.AddInt64(&cnt.inviteFail, int64(inviteFail))
		atomic.AddInt64(&cnt.agreeOK, int64(agreeOK))
		atomic.AddInt64(&cnt.agreeSkip, int64(agreeSkip))
		atomic.AddInt64(&cnt.agreeFail, int64(agreeFail))

		results[task.Index] = result
	})

	if err := writeLines(*roomsOutFile, collectRoomIDs(results)); err != nil {
		log.Fatalf("failed to write rooms file: %v", err)
	}
	if err := writeJSONL(*topologyOutFile, results); err != nil {
		log.Fatalf("failed to write topology file: %v", err)
	}

	log.Printf("room_warmup done: planned=%d create(ok=%d skipped=%d failed=%d) invite(ok=%d skipped=%d failed=%d) agree(ok=%d skipped=%d failed=%d) rooms=%s topology=%s",
		cnt.planned, cnt.createdOK, cnt.createdSkip, cnt.createdFail, cnt.inviteOK, cnt.inviteSkip, cnt.inviteFail, cnt.agreeOK, cnt.agreeSkip, cnt.agreeFail, *roomsOutFile, *topologyOutFile)

	if *strict && (cnt.createdFail > 0 || cnt.inviteFail > 0 || cnt.agreeFail > 0) {
		os.Exit(1)
	}
}

type roomBuildTask struct {
	Index      int
	OwnerID    string
	OwnerToken string
	MemberIDs  []string
	MemberToks []string
}

func warmupRoomMembers(task roomBuildTask, roomID string) (inviteOK, inviteSkip, inviteFail, agreeOK, agreeSkip, agreeFail int, msg string) {
	for idx, memberID := range task.MemberIDs {
		memberToken := task.MemberToks[idx]
		reqID, status, reqMsg, err := sendMemberRequest(memberToken, roomID, *joinContent, *roomPassword)
		if err != nil {
			inviteFail++
			msg = err.Error()
			continue
		}
		if status == "skipped" {
			inviteSkip++
			msg = reqMsg
		} else {
			inviteOK++
		}

		if reqID == "" {
			reqID, err = findLatestMemberRequestID(task.OwnerToken, roomID, memberID, *listLimit, *maxPages)
			if err != nil {
				agreeFail++
				msg = err.Error()
				continue
			}
		}
		if reqID == "" {
			agreeSkip++
			continue
		}

		agreeStatus, agreeMsg, err := agreeMemberRequest(task.OwnerToken, reqID)
		if err != nil {
			agreeFail++
			msg = err.Error()
			continue
		}
		if agreeStatus == "skipped" {
			agreeSkip++
			msg = agreeMsg
		} else {
			agreeOK++
		}
	}
	return
}

func createRoom(token string, maxMemberCount int, password string) (roomID string, status string, message string, err error) {
	body, _ := json.Marshal(map[string]any{
		"max_member_count": maxMemberCount,
		"password":         password,
	})
	req, err := http.NewRequest(http.MethodPost, *baseURL+createRoomPath, bytes.NewBuffer(body))
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", normalizeAuth(token))

	statusCode, apiResp, _, err := do(req)
	if err != nil {
		return "", "", "", err
	}
	message = strings.ToLower(apiResp.Message)
	if isSuccess(statusCode) {
		if isIdempotentRoomMessage(message) {
			return "", "skipped", message, nil
		}
		return "", "ok", message, nil
	}
	if isIdempotentRoomMessage(message) {
		return "", "skipped", message, nil
	}
	return "", "failed", message, fmt.Errorf("create room failed: status=%d message=%s", statusCode, apiResp.Message)
}

func findLatestRoomIDByOwner(token string, limit, maxPages int) (string, error) {
	baseID := ""
	for i := 0; i < maxPages; i++ {
		u := fmt.Sprintf("%s%s?limit=%d", *baseURL, listRoomshipsPath, limit)
		if baseID != "" {
			u += "&base_id=" + url.QueryEscape(baseID)
		}
		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", normalizeAuth(token))

		statusCode, apiResp, raw, err := do(req)
		if err != nil {
			return "", err
		}
		if !isSuccess(statusCode) {
			return "", fmt.Errorf("list roomships failed: status=%d message=%s", statusCode, apiResp.Message)
		}

		var data roomshipsResponseData
		if err := json.Unmarshal(apiResp.Data, &data); err != nil {
			return "", fmt.Errorf("unmarshal roomships data failed: %w body=%s", err, string(raw))
		}
		if len(data.Roomships) == 0 {
			return "", nil
		}
		if data.Roomships[0].RoomID != "" {
			return data.Roomships[0].RoomID, nil
		}
		baseID = data.Roomships[len(data.Roomships)-1].ID
	}
	return "", nil
}

func sendMemberRequest(token, roomID, content, password string) (requestID string, status string, message string, err error) {
	body, _ := json.Marshal(map[string]string{
		"room_id":  roomID,
		"password": password,
		"content":  content,
	})
	req, err := http.NewRequest(http.MethodPost, *baseURL+sendMemberRequestPath, bytes.NewBuffer(body))
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", normalizeAuth(token))

	statusCode, apiResp, _, err := do(req)
	if err != nil {
		return "", "", "", err
	}
	message = strings.ToLower(apiResp.Message)
	if isSuccess(statusCode) {
		return "", "ok", message, nil
	}
	if isIdempotentRoomMessage(message) {
		return "", "skipped", message, nil
	}
	return "", "failed", message, fmt.Errorf("send member request failed: status=%d message=%s", statusCode, apiResp.Message)
}

func findLatestMemberRequestID(token, roomID, applicantID string, limit, maxPages int) (string, error) {
	baseID := ""
	for i := 0; i < maxPages; i++ {
		u := fmt.Sprintf("%s%s?limit=%d", *baseURL, listMemberRequestsPath, limit)
		if baseID != "" {
			u += "&base_id=" + url.QueryEscape(baseID)
		}
		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", normalizeAuth(token))

		statusCode, apiResp, raw, err := do(req)
		if err != nil {
			return "", err
		}
		if !isSuccess(statusCode) {
			return "", fmt.Errorf("list member requests failed: status=%d message=%s", statusCode, apiResp.Message)
		}

		var data memberRequestsResponseData
		if err := json.Unmarshal(apiResp.Data, &data); err != nil {
			return "", fmt.Errorf("unmarshal member requests data failed: %w body=%s", err, string(raw))
		}
		if len(data.Requests) == 0 {
			return "", nil
		}
		for _, r := range data.Requests {
			if r.RoomID == roomID && r.ApplicantID == applicantID && strings.EqualFold(r.State, "pending") {
				return r.ID, nil
			}
		}
		baseID = data.Requests[len(data.Requests)-1].ID
	}
	return "", nil
}

func agreeMemberRequest(token, requestID string) (status string, message string, err error) {
	uri := fmt.Sprintf("%s/%s/agree", listMemberRequestsPath, url.PathEscape(requestID))
	req, err := http.NewRequest(http.MethodPut, *baseURL+uri, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", normalizeAuth(token))

	statusCode, apiResp, _, err := do(req)
	if err != nil {
		return "", "", err
	}
	message = strings.ToLower(apiResp.Message)
	if isSuccess(statusCode) {
		return "ok", message, nil
	}
	if isIdempotentRoomMessage(message) {
		return "skipped", message, nil
	}
	return "failed", message, fmt.Errorf("agree member request failed: status=%d message=%s", statusCode, apiResp.Message)
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

func normalizeAuth(token string) string {
	token = strings.TrimSpace(token)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		return token
	}
	return "Bearer " + token
}

func isSuccess(status int) bool {
	return status == http.StatusOK
}

func isIdempotentRoomMessage(message string) bool {
	m := strings.ToLower(message)
	keywords := []string{
		"already",
		"exist",
		"duplicate",
		"friend",
		"member",
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

func runWorkerPool(tasks []roomBuildTask, workers int, fn func(roomBuildTask)) {
	if workers <= 0 {
		workers = 1
	}
	ch := make(chan roomBuildTask, len(tasks))
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

func buildRoomTasks(tokens, userIDs []string, roomCount, roomSize int) []roomBuildTask {
	tasks := make([]roomBuildTask, 0, roomCount)
	for i := 0; i < roomCount; i++ {
		start := i * roomSize
		end := start + roomSize
		members := make([]string, 0, roomSize-1)
		memberTokens := make([]string, 0, roomSize-1)
		for j := start + 1; j < end; j++ {
			members = append(members, userIDs[j])
			memberTokens = append(memberTokens, tokens[j])
		}
		tasks = append(tasks, roomBuildTask{
			Index:      i,
			OwnerID:    userIDs[start],
			OwnerToken: tokens[start],
			MemberIDs:  members,
			MemberToks: memberTokens,
		})
	}
	return tasks
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

func writeLines(path string, lines []string) error {
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
	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func writeJSONL(path string, rows []roomTask) error {
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

func collectRoomIDs(rows []roomTask) []string {
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.RoomID != "" {
			ids = append(ids, row.RoomID)
		}
	}
	return ids
}
