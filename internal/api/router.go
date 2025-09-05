package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/jerson/mandrillfordev/internal/config"
	"github.com/jerson/mandrillfordev/internal/mailer"
	"github.com/jerson/mandrillfordev/internal/store"
	"github.com/jerson/mandrillfordev/internal/types"
)

func NewMux(cfg config.Config, st *store.Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/messages/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSend(w, r, cfg, st)
	})
	mux.HandleFunc("/messages/send.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSend(w, r, cfg, st)
	})
	mux.HandleFunc("/api/1.0/messages/send.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSend(w, r, cfg, st)
	})

	mux.HandleFunc("/messages/send-template", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSendTemplate(w, r, cfg, st)
	})
	mux.HandleFunc("/messages/send-template.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSendTemplate(w, r, cfg, st)
	})
	mux.HandleFunc("/api/1.0/messages/send-template.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSendTemplate(w, r, cfg, st)
	})

	mux.HandleFunc("/messages/send-raw", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSendRaw(w, r, cfg, st)
	})
	mux.HandleFunc("/messages/send-raw.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSendRaw(w, r, cfg, st)
	})
	mux.HandleFunc("/api/1.0/messages/send-raw.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSendRaw(w, r, cfg, st)
	})

	mux.HandleFunc("/messages/parse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleParse(w, r)
	})
	mux.HandleFunc("/messages/parse.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleParse(w, r)
	})
	mux.HandleFunc("/api/1.0/messages/parse.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleParse(w, r)
	})

	mux.HandleFunc("/messages/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleInfo(w, r, st)
	})
	mux.HandleFunc("/messages/info.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleInfo(w, r, st)
	})
	mux.HandleFunc("/api/1.0/messages/info.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleInfo(w, r, st)
	})

	mux.HandleFunc("/messages/content", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleContent(w, r, st)
	})
	mux.HandleFunc("/messages/content.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleContent(w, r, st)
	})
	mux.HandleFunc("/api/1.0/messages/content.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleContent(w, r, st)
	})

	mux.HandleFunc("/messages/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSearch(w, r, st)
	})
	mux.HandleFunc("/messages/search.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSearch(w, r, st)
	})
	mux.HandleFunc("/api/1.0/messages/search.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSearch(w, r, st)
	})

	mux.HandleFunc("/messages/search-time-series", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSearchTimeSeries(w, r, st)
	})
	mux.HandleFunc("/messages/search-time-series.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSearchTimeSeries(w, r, st)
	})
	mux.HandleFunc("/api/1.0/messages/search-time-series.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleSearchTimeSeries(w, r, st)
	})

	mux.HandleFunc("/messages/list-scheduled", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleListScheduled(w, r, st)
	})
	mux.HandleFunc("/messages/list-scheduled.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleListScheduled(w, r, st)
	})
	mux.HandleFunc("/api/1.0/messages/list-scheduled.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleListScheduled(w, r, st)
	})

	mux.HandleFunc("/messages/cancel-scheduled", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleCancelScheduled(w, r, st)
	})
	mux.HandleFunc("/messages/cancel-scheduled.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleCancelScheduled(w, r, st)
	})
	mux.HandleFunc("/api/1.0/messages/cancel-scheduled.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleCancelScheduled(w, r, st)
	})

	mux.HandleFunc("/messages/reschedule", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleReschedule(w, r, st)
	})
	mux.HandleFunc("/messages/reschedule.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleReschedule(w, r, st)
	})
	mux.HandleFunc("/api/1.0/messages/reschedule.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleReschedule(w, r, st)
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return mux
}

// Handlers
func handleSend(w http.ResponseWriter, r *http.Request, cfg config.Config, st *store.Store) {
	var req types.SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := requireKey(req.Key); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	var scheduledAt *time.Time
	if strings.TrimSpace(req.SendAt) != "" {
		if t, err := parseTime(req.SendAt); err == nil {
			scheduledAt = &t
		}
	}

	id := genID()
	rcpts := recipientsFromMessage(req.Message)
	if len(rcpts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no recipients"})
		return
	}

	rec := &types.MessageRecord{ID: id, CreatedAt: time.Now(), ScheduledAt: scheduledAt, Status: "queued", Message: req.Message, From: req.Message.FromEmail, To: rcpts, Subject: req.Message.Subject, Tags: req.Message.Tags}

	var results []types.SendResult
	for _, rcpt := range rcpts {
		results = append(results, types.SendResult{Email: rcpt, Status: "queued", ID: id})
	}

	if scheduledAt != nil && scheduledAt.After(time.Now()) {
		rec.Status = "scheduled"
		st.AddScheduled(rec)
		for i := range results {
			results[i].Status = "scheduled"
		}
		writeJSON(w, http.StatusOK, results)
		return
	}

	if err := mailer.SendMessage(cfg, req.Message, id, &rec.Raw); err != nil {
		rec.Status = "rejected"
		rec.RejectReason = err.Error()
		st.SaveMessage(rec)
		for i := range results {
			results[i].Status = "rejected"
			results[i].RejectReason = rec.RejectReason
		}
		writeJSON(w, http.StatusOK, results)
		return
	}
	now := time.Now()
	rec.SentAt = &now
	rec.Status = "sent"
	st.SaveMessage(rec)
	for i := range results {
		results[i].Status = "sent"
	}
	writeJSON(w, http.StatusOK, results)
}

func handleSendTemplate(w http.ResponseWriter, r *http.Request, cfg config.Config, st *store.Store) {
	var req types.SendTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := requireKey(req.Key); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	vars := map[string]string{}
	for _, tc := range req.TemplateContent {
		vars[tc.Name] = tc.Content
	}
	req.Message.HTML = replaceVars(req.Message.HTML, vars)
	req.Message.Text = replaceVars(req.Message.Text, vars)

	sr := types.SendRequest{Key: req.Key, Message: req.Message, Async: req.Async, IPPool: req.IPPool, SendAt: req.SendAt}
	var scheduledAt *time.Time
	if strings.TrimSpace(sr.SendAt) != "" {
		if t, err := parseTime(sr.SendAt); err == nil {
			scheduledAt = &t
		}
	}
	id := genID()
	rcpts := recipientsFromMessage(sr.Message)
	if len(rcpts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no recipients"})
		return
	}

	rec := &types.MessageRecord{ID: id, CreatedAt: time.Now(), ScheduledAt: scheduledAt, Status: "queued", Message: sr.Message, From: sr.Message.FromEmail, To: rcpts, Subject: sr.Message.Subject, Tags: sr.Message.Tags}
	var results []types.SendResult
	for _, rcpt := range rcpts {
		results = append(results, types.SendResult{Email: rcpt, Status: "queued", ID: id})
	}
	if scheduledAt != nil && scheduledAt.After(time.Now()) {
		rec.Status = "scheduled"
		st.AddScheduled(rec)
		for i := range results {
			results[i].Status = "scheduled"
		}
		writeJSON(w, http.StatusOK, results)
		return
	}
	if err := mailer.SendMessage(cfg, sr.Message, id, &rec.Raw); err != nil {
		rec.Status = "rejected"
		rec.RejectReason = err.Error()
		st.SaveMessage(rec)
		for i := range results {
			results[i].Status = "rejected"
			results[i].RejectReason = rec.RejectReason
		}
		writeJSON(w, http.StatusOK, results)
		return
	}
	now := time.Now()
	rec.SentAt = &now
	rec.Status = "sent"
	st.SaveMessage(rec)
	for i := range results {
		results[i].Status = "sent"
	}
	writeJSON(w, http.StatusOK, results)
}

func handleSendRaw(w http.ResponseWriter, r *http.Request, cfg config.Config, st *store.Store) {
	var req types.SendRawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := requireKey(req.Key); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	from := strings.TrimSpace(req.FromEmail)
	if from == "" {
		from = "no-reply@example.local"
	}
	to := make([]string, 0, len(req.To))
	for _, a := range req.To {
		if strings.TrimSpace(a) != "" {
			to = append(to, strings.TrimSpace(a))
		}
	}
	if len(to) == 0 {
		if hdrTo := extractHeader(req.RawMessage, "To"); hdrTo != "" {
			parts := strings.Split(hdrTo, ",")
			for _, p := range parts {
				if addr, err := mail.ParseAddress(strings.TrimSpace(p)); err == nil {
					to = append(to, addr.Address)
				}
			}
		}
	}
	if len(to) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no recipients"})
		return
	}

	id := genID()
	rec := &types.MessageRecord{ID: id, CreatedAt: time.Now(), Status: "queued", From: from, To: to, Subject: extractHeader(req.RawMessage, "Subject")}
	var scheduledAt *time.Time
	if strings.TrimSpace(req.SendAt) != "" {
		if t, err := parseTime(req.SendAt); err == nil {
			scheduledAt = &t
		}
	}
	if scheduledAt != nil && scheduledAt.After(time.Now()) {
		rec.ScheduledAt = scheduledAt
		rec.Raw = []byte(req.RawMessage)
		rec.Message = types.MandrillMessage{FromEmail: from, FromName: req.FromName, To: toRecipients(to)}
		rec.Status = "scheduled"
		st.AddScheduled(rec)
		var results []types.SendResult
		for _, rcpt := range to {
			results = append(results, types.SendResult{Email: rcpt, Status: "scheduled", ID: id})
		}
		writeJSON(w, http.StatusOK, results)
		return
	}
	if err := mailer.SendRaw(cfg, from, to, []byte(req.RawMessage)); err != nil {
		rec.Status = "rejected"
		rec.RejectReason = err.Error()
		st.SaveMessage(rec)
		var results []types.SendResult
		for _, rcpt := range to {
			results = append(results, types.SendResult{Email: rcpt, Status: "rejected", ID: id, RejectReason: rec.RejectReason})
		}
		writeJSON(w, http.StatusOK, results)
		return
	}
	now := time.Now()
	rec.SentAt = &now
	rec.Status = "sent"
	rec.Raw = []byte(req.RawMessage)
	st.SaveMessage(rec)
	var results []types.SendResult
	for _, rcpt := range to {
		results = append(results, types.SendResult{Email: rcpt, Status: "sent", ID: id})
	}
	writeJSON(w, http.StatusOK, results)
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	var req types.ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	msg, err := mail.ReadMessage(strings.NewReader(req.RawMessage))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid mime"})
		return
	}
	subj := msg.Header.Get("Subject")
	from := msg.Header.Get("From")
	to := msg.Header.Get("To")
	b, _ := io.ReadAll(msg.Body)
	writeJSON(w, http.StatusOK, map[string]any{"subject": subj, "from": from, "to": to, "raw": string(b), "headers": msg.Header})
}

func handleInfo(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.InfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if req.Id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing id"})
		return
	}
	if m, ok := st.GetMessage(req.Id); ok {
		writeJSON(w, http.StatusOK, m)
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
}

func handleContent(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.ContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if m, ok := st.GetMessage(req.Id); ok {
		writeJSON(w, http.StatusOK, map[string]any{"raw": string(m.Raw), "subject": m.Subject, "from": m.From, "to": m.To})
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
}

func handleSearch(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	var fromT, toT *time.Time
	if strings.TrimSpace(req.DateFrom) != "" {
		if t, err := parseTime(req.DateFrom); err == nil {
			fromT = &t
		}
	}
	if strings.TrimSpace(req.DateTo) != "" {
		if t, err := parseTime(req.DateTo); err == nil {
			toT = &t
		}
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	res := st.Search(req.Query, fromT, toT, req.Tags, req.Senders, limit)
	writeJSON(w, http.StatusOK, res)
}

func handleSearchTimeSeries(w http.ResponseWriter, r *http.Request, st *store.Store) {
	type point struct {
		Time string `json:"time"`
		Sent int    `json:"sent"`
	}
	now := time.Now().Truncate(time.Hour)
	start := now.Add(-7 * 24 * time.Hour)
	buckets := map[time.Time]int{}
	for t := start; !t.After(now); t = t.Add(time.Hour) {
		buckets[t] = 0
	}
	res := st.Search("", &start, &now, nil, nil, 0)
	for _, m := range res {
		if m.SentAt != nil {
			ts := m.SentAt.Truncate(time.Hour)
			if _, ok := buckets[ts]; ok && m.Status == "sent" {
				buckets[ts]++
			}
		}
	}
	out := make([]point, 0, len(buckets))
	for t := start; !t.After(now); t = t.Add(time.Hour) {
		out = append(out, point{Time: t.Format(time.RFC3339), Sent: buckets[t]})
	}
	writeJSON(w, http.StatusOK, out)
}

func handleListScheduled(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.ListScheduledRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	items := st.ListScheduled(req.To)
	writeJSON(w, http.StatusOK, items)
}

func handleCancelScheduled(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.CancelScheduledRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if m, ok := st.RemoveScheduled(req.Id); ok {
		m.Status = "canceled"
		st.SaveMessage(m)
		writeJSON(w, http.StatusOK, map[string]any{"status": "canceled", "id": req.Id})
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
}

func handleReschedule(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.RescheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	t, err := parseTime(req.SendAt)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid send_at"})
		return
	}
	if m, ok := st.GetScheduled(req.Id); ok {
		m.ScheduledAt = &t
		m.Status = "scheduled"
		st.SaveMessage(m)
		writeJSON(w, http.StatusOK, m)
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
}

// Helpers
func parseTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	fmts := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02T15:04:05", "2006-01-02"}
	var last error
	for _, f := range fmts {
		if t, err := time.ParseInLocation(f, s, time.Local); err == nil {
			return t, nil
		} else {
			last = err
		}
	}
	return time.Time{}, last
}

func extractHeader(raw string, key string) string {
	lines := strings.Split(raw, "\n")
	keyLower := strings.ToLower(key) + ":"
	var val strings.Builder
	found := false
	for _, l := range lines {
		if !found {
			if strings.HasPrefix(strings.ToLower(strings.TrimSpace(l)), keyLower) {
				found = true
				v := l[strings.Index(l, ":")+1:]
				val.WriteString(strings.TrimSpace(v))
			}
		} else {
			if len(l) > 0 && (l[0] == ' ' || l[0] == '\t') {
				val.WriteString(" ")
				val.WriteString(strings.TrimSpace(l))
			} else {
				break
			}
		}
	}
	return val.String()
}

func toRecipients(to []string) []types.MandrillRecipient {
	out := make([]types.MandrillRecipient, 0, len(to))
	for _, a := range to {
		out = append(out, types.MandrillRecipient{Email: a, Type: "to"})
	}
	return out
}

func replaceVars(s string, vars map[string]string) string {
	if s == "" {
		return s
	}
	out := s
	for k, v := range vars {
		token := fmt.Sprintf("*|%s|*", k)
		out = strings.ReplaceAll(out, token, v)
	}
	return out
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func genID() string { b := make([]byte, 12); _, _ = rand.Read(b); return hex.EncodeToString(b) }

func requireKey(key string) error {
	allowed := strings.TrimSpace(os.Getenv("MANDRILL_KEYS"))
	if allowed == "" {
		return nil
	}
	keys := strings.Split(allowed, ",")
	for _, k := range keys {
		if strings.TrimSpace(k) == key {
			return nil
		}
	}
	return fmt.Errorf("invalid mandrill api key")
}

func recipientsFromMessage(m types.MandrillMessage) []string {
	rcpts := make([]string, 0, len(m.To)+1)
	for _, r := range m.To {
		e := strings.TrimSpace(r.Email)
		if e != "" {
			rcpts = append(rcpts, e)
		}
	}
	if strings.TrimSpace(m.BccAddress) != "" {
		rcpts = append(rcpts, strings.TrimSpace(m.BccAddress))
	}
	return rcpts
}
