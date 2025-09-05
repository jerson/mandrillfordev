package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jerson/mandrillfordev/internal/api"
	"github.com/jerson/mandrillfordev/internal/config"
	"github.com/jerson/mandrillfordev/internal/scheduler"
	"github.com/jerson/mandrillfordev/internal/store"
)

func main() {
	// Special mode: healthcheck helper for Docker
	if len(os.Args) > 1 && os.Args[1] == "-healthcheck" {
		if err := doHealthcheck(); err != nil {
			log.Printf("healthcheck failed: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	cfg := config.Load()
	st := store.NewStore()
	sched := scheduler.NewScheduler(cfg, st)
	sched.Start()

	mux := api.NewMux(cfg, st)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("Mandrill-dev server listening on %s", addr)
	if err := http.ListenAndServe(addr, loggingMiddleware(mux)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if !isDebug() {
			// Basic logging
			next.ServeHTTP(w, r)
			dur := time.Since(start)
			log.Printf("%s %s %s", r.Method, r.URL.Path, dur)
			return
		}

		// Read and restore request body
		var reqBody []byte
		if r.Body != nil {
			b, _ := io.ReadAll(r.Body)
			reqBody = b
			r.Body = io.NopCloser(bytes.NewBuffer(b))
		}

		// Wrap ResponseWriter to capture response
		lrw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lrw, r)
		dur := time.Since(start)

		// Truncate bodies for log safety
		const maxLog = 100 * 1024 // 100KB
		rb := reqBody
		if len(rb) > maxLog {
			rb = append(rb[:maxLog], []byte("... (truncated)")...)
		}
		sb := lrw.buf.Bytes()
		if len(sb) > maxLog {
			sb = append(sb[:maxLog], []byte("... (truncated)")...)
		}

		log.Printf("[DEBUG] %s %s %s\nRequestBody: %s\nResponseStatus: %d\nResponseBody: %s",
			r.Method, r.URL.Path, dur, string(rb), lrw.status, string(sb))
	})
}

// loggingResponseWriter captures status code and body for logging
type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	buf    bytes.Buffer
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	// Capture copy of body
	_, _ = w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func isDebug() bool {
	v := strings.TrimSpace(os.Getenv("MANDRILL_DEBUG"))
	if v == "" {
		v = strings.TrimSpace(os.Getenv("DEBUG"))
	}
	v = strings.ToLower(v)
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func doHealthcheck() error {
	port := os.Getenv("PORT")
	if strings.TrimSpace(port) == "" {
		port = "8080"
	}
	url := "http://127.0.0.1:" + port + "/healthz"
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}
