package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractRecipientID(t *testing.T) {
	t.Parallel()

	token := buildDryRunAccessToken("019e40df-fdd4-7692-9e4a-37c73ead3e27")
	recipientID, err := extractRecipientID(token)
	if err != nil {
		t.Fatalf("extractRecipientID returned error: %v", err)
	}
	if recipientID != "019e40df-fdd4-7692-9e4a-37c73ead3e27" {
		t.Fatalf("unexpected recipientID: %s", recipientID)
	}
}

func TestWriteOutputsWritesRecipients(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	tokenPath := filepath.Join(root, "tokens.txt")
	detailPath := filepath.Join(root, "token_pool.jsonl")
	recipientPath := filepath.Join(root, "recipients.txt")

	results := []result{
		{
			index:       0,
			idx:         1,
			email:       "bench.20260519153442.1@test.com",
			userNumber:  "2056760889289216002",
			recipientID: "019e40df-fdd4-7692-9e4a-37c73ead3e27",
			accessToken: buildDryRunAccessToken("019e40df-fdd4-7692-9e4a-37c73ead3e27"),
		},
		{
			index: 1,
			idx:   2,
			err:   os.ErrNotExist,
		},
	}

	if err := writeOutputs(tokenPath, detailPath, recipientPath, results, false); err != nil {
		t.Fatalf("writeOutputs returned error: %v", err)
	}

	tokenBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		t.Fatalf("read token file failed: %v", err)
	}
	if got := strings.TrimSpace(string(tokenBytes)); got == "" || !strings.HasPrefix(got, "Bearer ") {
		t.Fatalf("unexpected token file content: %q", got)
	}

	recipientBytes, err := os.ReadFile(recipientPath)
	if err != nil {
		t.Fatalf("read recipient file failed: %v", err)
	}
	if got := strings.TrimSpace(string(recipientBytes)); got != "019e40df-fdd4-7692-9e4a-37c73ead3e27" {
		t.Fatalf("unexpected recipient file content: %q", got)
	}

	detailBytes, err := os.ReadFile(detailPath)
	if err != nil {
		t.Fatalf("read detail file failed: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(detailBytes)), "\n")
	if len(lines) != 1 {
		t.Fatalf("unexpected detail line count: %d", len(lines))
	}
	var detail map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &detail); err != nil {
		t.Fatalf("unmarshal detail line failed: %v", err)
	}
	if detail["user_id"] != "019e40df-fdd4-7692-9e4a-37c73ead3e27" {
		t.Fatalf("unexpected detail user_id: %#v", detail["user_id"])
	}
}
