package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Minimal Mandrill-like types for the client example
type Recipient struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
}

type Message struct {
	HTML      string            `json:"html,omitempty"`
	Text      string            `json:"text,omitempty"`
	Subject   string            `json:"subject,omitempty"`
	FromEmail string            `json:"from_email"`
	FromName  string            `json:"from_name,omitempty"`
	To        []Recipient       `json:"to"`
	Headers   map[string]string `json:"headers,omitempty"`
}

type SendRequest struct {
	Key     string  `json:"key"`
	Message Message `json:"message"`
	SendAt  string  `json:"send_at,omitempty"`
}

func main() {
	base := getenv("API_BASE", "http://localhost:8080")
	key := getenv("KEY", "dev")

	// Wait for health
	if err := waitFor(base+"/healthz", 30*time.Second); err != nil {
		fmt.Println("dev server not ready:", err)
		os.Exit(1)
	}

	// Build a send request
	req := SendRequest{
		Key: key,
		Message: Message{
			FromEmail: getenv("FROM", "sender@example.com"),
			FromName:  getenv("FROM_NAME", "Sender"),
			Subject:   "Mandrill GO Dev Test",
			Text:      "Hello from client example (text).",
			HTML:      "<p><b>Hello</b> from client example (HTML).</p>",
			To:        []Recipient{{Email: getenv("TO", "user@example.com"), Type: "to"}},
			Headers:   map[string]string{"Reply-To": getenv("REPLY_TO", "reply@example.com")},
		},
		// SendAt: time.Now().Add(10*time.Second).Format(time.RFC3339), // uncomment to test scheduling
	}

	b, _ := json.Marshal(req)
	url := strings.TrimRight(base, "/") + "/messages/send"
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		fmt.Println("request failed:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %s\nResponse: %s\n", resp.Status, string(body))
}

func getenv(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}

func waitFor(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for %s", url)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
