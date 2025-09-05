package config

import (
	"os"
	"strconv"
)

type SMTPMode string

const (
	TLSNone     SMTPMode = "none"
	TLSStartTLS SMTPMode = "starttls"
	TLSTLS      SMTPMode = "tls"
)

type Config struct {
	SMTPHost        string
	SMTPPort        int
	SMTPUsername    string
	SMTPPassword    string
	SMTPMode        SMTPMode
	InsecureTLS     bool
	DefaultFromName string
}

func envOr(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func Load() Config {
	port, _ := strconv.Atoi(envOr("SMTP_PORT", "1025"))
	mode := envOr("SMTP_TLS", "none")
	insecure := envOr("SMTP_INSECURE_TLS", "false") == "true"
	return Config{
		SMTPHost:        envOr("SMTP_HOST", "localhost"),
		SMTPPort:        port,
		SMTPUsername:    envOr("SMTP_USERNAME", ""),
		SMTPPassword:    envOr("SMTP_PASSWORD", ""),
		SMTPMode:        SMTPMode(mode),
		InsecureTLS:     insecure,
		DefaultFromName: envOr("DEFAULT_FROM_NAME", "Mandrill Dev"),
	}
}
