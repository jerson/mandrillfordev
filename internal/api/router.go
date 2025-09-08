package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
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

	// Templates endpoints
	// add
	mux.HandleFunc("/templates/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateAdd(w, r, st)
	})
	mux.HandleFunc("/templates/add.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateAdd(w, r, st)
	})
	mux.HandleFunc("/api/1.0/templates/add.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateAdd(w, r, st)
	})
	// info
	mux.HandleFunc("/templates/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateInfo(w, r, st)
	})
	mux.HandleFunc("/templates/info.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateInfo(w, r, st)
	})
	mux.HandleFunc("/api/1.0/templates/info.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateInfo(w, r, st)
	})
	// update
	mux.HandleFunc("/templates/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateUpdate(w, r, st)
	})
	mux.HandleFunc("/templates/update.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateUpdate(w, r, st)
	})
	mux.HandleFunc("/api/1.0/templates/update.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateUpdate(w, r, st)
	})
	// publish
	mux.HandleFunc("/templates/publish", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplatePublish(w, r, st)
	})
	mux.HandleFunc("/templates/publish.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplatePublish(w, r, st)
	})
	mux.HandleFunc("/api/1.0/templates/publish.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplatePublish(w, r, st)
	})
	// delete
	mux.HandleFunc("/templates/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateDelete(w, r, st)
	})
	mux.HandleFunc("/templates/delete.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateDelete(w, r, st)
	})
	mux.HandleFunc("/api/1.0/templates/delete.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateDelete(w, r, st)
	})
	// list
	mux.HandleFunc("/templates/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateList(w, r, st)
	})
	mux.HandleFunc("/templates/list.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateList(w, r, st)
	})
	mux.HandleFunc("/api/1.0/templates/list.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateList(w, r, st)
	})
	// time-series
	mux.HandleFunc("/templates/time-series", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateTimeSeries(w, r, st)
	})
	mux.HandleFunc("/templates/time-series.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateTimeSeries(w, r, st)
	})
	mux.HandleFunc("/api/1.0/templates/time-series.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateTimeSeries(w, r, st)
	})
	// render
	mux.HandleFunc("/templates/render", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateRender(w, r, st)
	})
	mux.HandleFunc("/templates/render.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateRender(w, r, st)
	})
	mux.HandleFunc("/api/1.0/templates/render.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		handleTemplateRender(w, r, st)
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
	// In debug mode, append the original (sanitized) request to the message body for troubleshooting
	if isDebug() {
		dbg := map[string]any{
			"template_name":    req.TemplateName,
			"template_content": req.TemplateContent,
			"message":          req.Message,
			"async":            req.Async,
			"ip_pool":          req.IPPool,
			"send_at":          req.SendAt,
			"_note":            "Debug info: original send-template request (key omitted)",
		}
		if b, err := json.MarshalIndent(dbg, "", "  "); err == nil {
			// Text part
			if strings.TrimSpace(sr.Message.Text) == "" {
				sr.Message.Text = string(b)
			} else {
				sr.Message.Text += "\n\n---- debug: send-template request ----\n" + string(b)
			}
			// HTML part
			esc := html.EscapeString(string(b))
			if strings.TrimSpace(sr.Message.HTML) == "" {
				sr.Message.HTML = "<pre style=\"white-space:pre-wrap\">" + esc + "</pre>"
			} else {
				sr.Message.HTML += "<hr><h4>Debug: send-template request</h4><pre style=\"white-space:pre-wrap\">" + esc + "</pre>"
			}
		}
	}
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

	// include template name for later stats and discovery
	tags := append([]string{}, sr.Message.Tags...)
	if strings.TrimSpace(req.TemplateName) != "" {
		tags = append(tags, "template:"+req.TemplateName)
	}
	rec := &types.MessageRecord{ID: id, CreatedAt: time.Now(), ScheduledAt: scheduledAt, Status: "queued", Message: sr.Message, From: sr.Message.FromEmail, To: rcpts, Subject: sr.Message.Subject, Tags: tags, TemplateName: req.TemplateName}
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

// Templates Handlers
func handleTemplateAdd(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.TemplateAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := requireKey(req.Key); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing name"})
		return
	}
	now := time.Now()
	t := &types.Template{
		Name:      name,
		FromEmail: req.FromEmail,
		FromName:  req.FromName,
		Subject:   req.Subject,
		Code:      req.Code,
		Text:      req.Text,
		Labels:    append([]string{}, req.Labels...),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if req.Publish {
		t.PublishedCode = t.Code
		t.PublishedText = t.Text
		t.PublishedAt = &now
	}
	st.SaveTemplate(t)
	writeJSON(w, http.StatusOK, t)
}

func handleTemplateInfo(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.TemplateInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := requireKey(req.Key); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	if t, ok := st.GetTemplate(req.Name); ok {
		writeJSON(w, http.StatusOK, t)
		return
	}
	// Not found: return a synthetic template to "fake it's there"
	now := time.Now()
	t := &types.Template{
		Name:      strings.TrimSpace(req.Name),
		CreatedAt: now,
		UpdatedAt: now,
	}
	writeJSON(w, http.StatusOK, t)
}

func handleTemplateUpdate(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.TemplateUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := requireKey(req.Key); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	t, ok := st.GetTemplate(req.Name)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	changed := false
	if req.FromEmail != nil {
		t.FromEmail = *req.FromEmail
		changed = true
	}
	if req.FromName != nil {
		t.FromName = *req.FromName
		changed = true
	}
	if req.Subject != nil {
		t.Subject = *req.Subject
		changed = true
	}
	if req.Code != nil {
		t.Code = *req.Code
		changed = true
	}
	if req.Text != nil {
		t.Text = *req.Text
		changed = true
	}
	if req.Labels != nil {
		t.Labels = append([]string{}, (*req.Labels)...)
		changed = true
	}
	if req.Publish != nil && *req.Publish {
		now := time.Now()
		t.PublishedCode = t.Code
		t.PublishedText = t.Text
		t.PublishedAt = &now
		changed = true
	}
	if changed {
		t.UpdatedAt = time.Now()
		st.SaveTemplate(t)
	}
	writeJSON(w, http.StatusOK, t)
}

func handleTemplatePublish(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.TemplatePublishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := requireKey(req.Key); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	t, ok := st.GetTemplate(req.Name)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	now := time.Now()
	t.PublishedCode = t.Code
	t.PublishedText = t.Text
	t.PublishedAt = &now
	t.UpdatedAt = now
	st.SaveTemplate(t)
	writeJSON(w, http.StatusOK, t)
}

func handleTemplateDelete(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.TemplateDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := requireKey(req.Key); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	if t, ok := st.DeleteTemplate(req.Name); ok {
		writeJSON(w, http.StatusOK, t)
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
}

func handleTemplateList(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.TemplateListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := requireKey(req.Key); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	list := st.ListTemplates(req.Label)
	writeJSON(w, http.StatusOK, list)
}

func handleTemplateTimeSeries(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.TemplateTimeSeriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := requireKey(req.Key); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	// Build hourly buckets for the last 30 days
	type point struct {
		Time string `json:"time"`
		Sent int    `json:"sent"`
	}
	now := time.Now().Truncate(time.Hour)
	start := now.Add(-30 * 24 * time.Hour)
	buckets := map[time.Time]int{}
	for t := start; !t.After(now); t = t.Add(time.Hour) {
		buckets[t] = 0
	}
	for _, m := range st.Messages() {
		if m.SentAt == nil || m.Status != "sent" {
			continue
		}
		if !strings.EqualFold(m.TemplateName, req.Name) {
			continue
		}
		ts := m.SentAt.Truncate(time.Hour)
		if _, ok := buckets[ts]; ok {
			buckets[ts]++
		}
	}
	out := make([]point, 0, len(buckets))
	for t := start; !t.After(now); t = t.Add(time.Hour) {
		out = append(out, point{Time: t.Format(time.RFC3339), Sent: buckets[t]})
	}
	writeJSON(w, http.StatusOK, out)
}

func handleTemplateRender(w http.ResponseWriter, r *http.Request, st *store.Store) {
	var req types.TemplateRenderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := requireKey(req.Key); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	t, _ := st.GetTemplate(req.TemplateName)
	// Prefer published code, fallback to draft code if none
	code := ""
	if t != nil {
		if strings.TrimSpace(t.PublishedCode) != "" {
			code = t.PublishedCode
		} else {
			code = t.Code
		}
	}
	// If no stored template, render the merge-only content; some clients may expect variable injection without stored template
	vars := map[string]string{}
	for _, tc := range req.TemplateContent {
		vars[tc.Name] = tc.Content
	}
	htmlOut := replaceVars(code, vars)
	writeJSON(w, http.StatusOK, map[string]any{"html": htmlOut})
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

// isDebug reports whether debug logging/behaviors are enabled via env
func isDebug() bool {
	v := strings.TrimSpace(os.Getenv("MANDRILL_DEBUG"))
	if v == "" {
		v = strings.TrimSpace(os.Getenv("DEBUG"))
	}
	v = strings.ToLower(v)
	return v == "1" || v == "true" || v == "yes" || v == "on"
}
