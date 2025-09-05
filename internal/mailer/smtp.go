package mailer

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"sort"
	"strings"
	"time"

	"github.com/jerson/mandrillfordev/internal/config"
	"github.com/jerson/mandrillfordev/internal/types"
)

func SendMessage(cfg config.Config, mm types.MandrillMessage, id string, outRaw *[]byte) error {
	from, toHdr, ccHdr, rcpts := extractRecipients(mm)
	raw := buildRFC822(mm, id, from, toHdr, ccHdr)
	if outRaw != nil {
		*outRaw = raw
	}
	return smtpSend(cfg, from, rcpts, raw)
}

func SendRaw(cfg config.Config, from string, to []string, raw []byte) error {
	return smtpSend(cfg, from, to, raw)
}

func smtpSend(cfg config.Config, from string, rcpts []string, raw []byte) error {
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)

	var c *smtp.Client
	var err error
	switch cfg.SMTPMode {
	case config.TLSTLS:
		tlsConn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: cfg.SMTPHost, InsecureSkipVerify: cfg.InsecureTLS})
		if err != nil {
			return err
		}
		c, err = smtp.NewClient(tlsConn, cfg.SMTPHost)
		if err != nil {
			return err
		}
	default:
		c, err = smtp.Dial(addr)
		if err != nil {
			return err
		}
		// STARTTLS if requested
		if cfg.SMTPMode == config.TLSStartTLS {
			if ok, _ := c.Extension("STARTTLS"); ok {
				if err := c.StartTLS(&tls.Config{ServerName: cfg.SMTPHost, InsecureSkipVerify: cfg.InsecureTLS}); err != nil {
					_ = c.Close()
					return err
				}
			}
		}
	}
	defer c.Close()

	if cfg.SMTPUsername != "" {
		if ok, _ := c.Extension("AUTH"); ok {
			auth := smtp.PlainAuth("", cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPHost)
			if err := c.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	for _, r := range rcpts {
		if err := c.Rcpt(r); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(raw); err != nil {
		_ = w.Close()
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

func extractRecipients(mm types.MandrillMessage) (from string, toHdr []string, ccHdr []string, rcpts []string) {
	from = mm.FromEmail
	if from == "" {
		from = "no-reply@example.local"
	}
	var to, cc, bcc []string
	for _, r := range mm.To {
		t := strings.ToLower(strings.TrimSpace(r.Type))
		if t == "" {
			t = "to"
		}
		switch t {
		case "to":
			to = append(to, formatAddress(r.Email, r.Name))
		case "cc":
			cc = append(cc, formatAddress(r.Email, r.Name))
		case "bcc":
			bcc = append(bcc, r.Email)
		}
		rcpts = append(rcpts, r.Email)
	}
	if mm.BccAddress != "" {
		rcpts = append(rcpts, mm.BccAddress)
	}
	return from, to, cc, rcpts
}

func formatAddress(email, name string) string {
	if name == "" {
		return (&mail.Address{Address: email}).String()
	}
	return (&mail.Address{Name: name, Address: email}).String()
}

func buildRFC822(mm types.MandrillMessage, id, from string, toHdr, ccHdr []string) []byte {
	now := time.Now()
	var buf bytes.Buffer

	// Standard headers
	fromName := mm.FromName
	if fromName == "" {
		fromName = "Mandrill Dev"
	}
	fmt.Fprintf(&buf, "Date: %s\r\n", now.Format(time.RFC1123Z))
	fmt.Fprintf(&buf, "From: %s\r\n", formatAddress(from, fromName))
	if len(toHdr) > 0 {
		fmt.Fprintf(&buf, "To: %s\r\n", strings.Join(toHdr, ", "))
	}
	if len(ccHdr) > 0 {
		fmt.Fprintf(&buf, "Cc: %s\r\n", strings.Join(ccHdr, ", "))
	}
	if mm.Subject != "" {
		fmt.Fprintf(&buf, "Subject: %s\r\n", sanitizeHeader(mm.Subject))
	}
	fmt.Fprintf(&buf, "Message-ID: <%s@mandrill-dev.local>\r\n", id)
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")

	// Additional headers
	keys := make([]string, 0, len(mm.Headers))
	for k := range mm.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if strings.EqualFold(k, "Bcc") {
			continue
		} // don't add bcc header
		v := mm.Headers[k]
		if v != "" {
			fmt.Fprintf(&buf, "%s: %s\r\n", k, sanitizeHeader(v))
		}
	}

	// Build body with MIME structure
	hasHTML := strings.TrimSpace(mm.HTML) != ""
	hasText := strings.TrimSpace(mm.Text) != ""
	hasAttachments := len(mm.Attachments) > 0 || len(mm.Images) > 0

	if hasAttachments {
		mixedBoundary := boundary("mixed")
		fmt.Fprintf(&buf, "Content-Type: multipart/mixed; boundary=%q\r\n\r\n", mixedBoundary)
		// First part: alternative or single
		if hasHTML && hasText {
			// multipart/alternative
			altBoundary := boundary("alt")
			fmt.Fprintf(&buf, "--%s\r\n", mixedBoundary)
			fmt.Fprintf(&buf, "Content-Type: multipart/alternative; boundary=%q\r\n\r\n", altBoundary)
			writeTextPart(&buf, altBoundary, mm.Text)
			writeHTMLPart(&buf, altBoundary, mm.HTML)
			fmt.Fprintf(&buf, "--%s--\r\n", altBoundary)
		} else if hasHTML {
			fmt.Fprintf(&buf, "--%s\r\n", mixedBoundary)
			writeHTMLPart(&buf, "", mm.HTML)
		} else {
			fmt.Fprintf(&buf, "--%s\r\n", mixedBoundary)
			writeTextPart(&buf, "", mm.Text)
		}
		// Attachments
		for _, a := range mm.Attachments {
			writeAttachment(&buf, mixedBoundary, a, false)
		}
		for _, a := range mm.Images {
			writeAttachment(&buf, mixedBoundary, a, true)
		}
		fmt.Fprintf(&buf, "--%s--\r\n", mixedBoundary)
	} else if hasHTML && hasText {
		altBoundary := boundary("alt")
		fmt.Fprintf(&buf, "Content-Type: multipart/alternative; boundary=%q\r\n\r\n", altBoundary)
		writeTextPart(&buf, altBoundary, mm.Text)
		writeHTMLPart(&buf, altBoundary, mm.HTML)
		fmt.Fprintf(&buf, "--%s--\r\n", altBoundary)
	} else if hasHTML {
		fmt.Fprintf(&buf, "Content-Type: text/html; charset=utf-8\r\n\r\n")
		buf.WriteString(mm.HTML)
		if !strings.HasSuffix(mm.HTML, "\r\n") {
			buf.WriteString("\r\n")
		}
	} else {
		fmt.Fprintf(&buf, "Content-Type: text/plain; charset=utf-8\r\n\r\n")
		buf.WriteString(mm.Text)
		if !strings.HasSuffix(mm.Text, "\r\n") {
			buf.WriteString("\r\n")
		}
	}
	return buf.Bytes()
}

func sanitizeHeader(v string) string {
	v = strings.ReplaceAll(v, "\r", " ")
	v = strings.ReplaceAll(v, "\n", " ")
	return v
}

func boundary(prefix string) string {
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), time.Now().Unix()%1000)
}

func writeTextPart(buf *bytes.Buffer, boundary, body string) {
	if boundary != "" {
		fmt.Fprintf(buf, "--%s\r\n", boundary)
	}
	fmt.Fprintf(buf, "Content-Type: text/plain; charset=utf-8\r\n\r\n")
	buf.WriteString(body)
	if !strings.HasSuffix(body, "\r\n") {
		buf.WriteString("\r\n")
	}
}

func writeHTMLPart(buf *bytes.Buffer, boundary, body string) {
	if boundary != "" {
		fmt.Fprintf(buf, "--%s\r\n", boundary)
	}
	fmt.Fprintf(buf, "Content-Type: text/html; charset=utf-8\r\n\r\n")
	buf.WriteString(body)
	if !strings.HasSuffix(body, "\r\n") {
		buf.WriteString("\r\n")
	}
}

func writeAttachment(buf *bytes.Buffer, mixedBoundary string, a types.MandrillAttachment, inline bool) {
	fmt.Fprintf(buf, "--%s\r\n", mixedBoundary)
	name := a.Name
	if name == "" {
		name = "attachment"
	}
	dispo := "attachment"
	if inline {
		dispo = "inline"
	}
	ctype := a.Type
	if ctype == "" {
		ctype = "application/octet-stream"
	}
	fmt.Fprintf(buf, "Content-Type: %s; name=%q\r\n", ctype, name)
	fmt.Fprintf(buf, "Content-Disposition: %s; filename=%q\r\n", dispo, name)
	if inline {
		// allow referencing via cid:name in HTML
		fmt.Fprintf(buf, "Content-ID: <%s>\r\n", name)
	}
	fmt.Fprintf(buf, "Content-Transfer-Encoding: base64\r\n\r\n")
	// content is base64-encoded string already; ensure canonical wrap at 76 chars
	// decode and re-encode to enforce wrapping
	raw, err := base64.StdEncoding.DecodeString(a.Content)
	if err != nil {
		// fallback: write as-is
		writeBase64Wrapped(buf, []byte(a.Content))
	} else {
		enc := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
		base64.StdEncoding.Encode(enc, raw)
		writeBase64Wrapped(buf, enc)
	}
	if !strings.HasSuffix(buf.String(), "\r\n") {
		buf.WriteString("\r\n")
	}
}

func writeBase64Wrapped(buf *bytes.Buffer, b []byte) {
	for len(b) > 76 {
		buf.Write(b[:76])
		buf.WriteString("\r\n")
		b = b[76:]
	}
	if len(b) > 0 {
		buf.Write(b)
		buf.WriteString("\r\n")
	}
}
