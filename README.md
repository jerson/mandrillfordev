Mandrill-compatible SMTP Redirector (Go)

Overview

- Implements a subset of Mandrill Transactional API endpoints for local development.
- Accepts Mandrill-style JSON requests and relays messages to a regular SMTP server.
- Stores messages in-memory for search, info, and content endpoints.
- Supports simple scheduling of messages.

Endpoints

- POST `/messages/send` and `/messages/send.json`
- POST `/api/1.0/messages/send.json` (alias for SDK compatibility)
- POST `/messages/send-template` and `/messages/send-template.json`
- POST `/api/1.0/messages/send-template.json`
- POST `/messages/send-raw` and `/messages/send-raw.json`
- POST `/api/1.0/messages/send-raw.json`
- POST `/messages/parse` and `/messages/parse.json`
- POST `/api/1.0/messages/parse.json`
- POST `/messages/info` and `/messages/info.json`
- POST `/api/1.0/messages/info.json`
- POST `/messages/content` and `/messages/content.json`
- POST `/api/1.0/messages/content.json`
- POST `/messages/search` and `/messages/search.json`
- POST `/api/1.0/messages/search.json`
- POST `/messages/search-time-series` and `/messages/search-time-series.json`
- POST `/api/1.0/messages/search-time-series.json`
- POST `/messages/list-scheduled` and `/messages/list-scheduled.json`
- POST `/api/1.0/messages/list-scheduled.json`
- POST `/messages/cancel-scheduled` and `/messages/cancel-scheduled.json`
- POST `/api/1.0/messages/cancel-scheduled.json`
- POST `/messages/reschedule` and `/messages/reschedule.json`
- POST `/api/1.0/messages/reschedule.json`
- GET `/healthz`

Configuration (env)

- `SMTP_HOST` (default: `localhost`)
- `SMTP_PORT` (default: `1025`)
- `SMTP_USERNAME` (optional)
- `SMTP_PASSWORD` (optional)
- `SMTP_TLS` one of `none` (default), `starttls`, `tls`
- `SMTP_INSECURE_TLS` `true|false` (default: `false`)
- `MANDRILL_KEYS` comma-separated list. If set, incoming `key` must match one of these.
- `DEFAULT_FROM_NAME` default sender name when missing (default: `Mandrill Dev`).
- `PORT` HTTP port (default: `8080`).

Run locally

1) Run an SMTP server (e.g., MailHog): `docker run -p 1025:1025 -p 8025:8025 mailhog/mailhog`
2) Start server:

```
SMTP_HOST=localhost SMTP_PORT=1025 PORT=8080 go run ./cmd/mandrill-dev
```

Docker

```
docker build -t mandrill-dev .
docker run --rm -p 8080:8080 \
  -e SMTP_HOST=host.docker.internal \
  -e SMTP_PORT=1025 \
  mandrill-dev
```

Docker Compose (with smtp4dev)

This repo includes a `docker-compose.yml` that brings up smtp4dev and this Mandrill-compatible server together.

```
docker compose -f compose.yml up --build
```

- Mandrill API: `http://localhost:8080`
- smtp4dev UI: `http://localhost:3000`

The Mandrill server relays email to the smtp4dev container (`SMTP_HOST=smtp4dev`, `SMTP_PORT=25`).

Additional example services:
- `bun-http-client` posts to `/api/1.0/messages/send.json`.
- `bun-sdk-client` uses the SDK example configured to talk to the local server via `MC_BASE`.
- `bun-send-template-client` exercises `/api/1.0/messages/send-template.json` with merge vars.

Health checks

- The server exposes `GET /healthz` which returns `200 OK` and `ok` body.
- The Docker image includes a `HEALTHCHECK` that runs `/mandrill-dev -healthcheck` against `http://127.0.0.1:$PORT/healthz`.

Client example

- Go client: `examples/go-client/main.go` (Compose service: `client`).
- Node client (HTTP): `examples/node-client/http-client.js` (Compose service: `node-http-client`).
- Node client (Official SDK): `examples/node-client/sdk-client.js` — by default the SDK targets mandrillapp.com. When `MC_BASE` is set (e.g. `http://localhost:8080/api/1.0`), this example will post directly to the local dev server instead of the public API.
- Node client (send-template): `examples/node-client/send-template.js` — uses the official SDK by default. If `MC_BASE` is set (e.g., `http://localhost:8080/api/1.0`), it posts directly to the local dev server’s `/messages/send-template.json`.

Run the client with Compose:

```
docker compose -f compose.yml up --build client
```

Or run locally against a running server:

```
API_BASE=http://localhost:8080 \
KEY=dev \
FROM=sender@example.com \
TO=user@example.com \
go run ./examples/go-client
```

Node HTTP client (no SDK):

```
API_BASE=http://localhost:8080 \
KEY=dev \
FROM=sender@example.com \
TO=user@example.com \
node examples/node-client/http-client.js
```

Node SDK client (requires package install):

```
cd examples/node-client
npm install
# Point to local server via MC_BASE
MC_BASE=http://localhost:8080/api/1.0 \
KEY=dev \
FROM=sender@example.com \
TO=user@example.com \
node sdk-client.js
```

Examples

- Send:

```
curl -s localhost:8080/messages/send -H 'Content-Type: application/json' -d '{
  "key":"dev",
  "message":{
    "from_email":"sender@example.com",
    "from_name":"Sender",
    "subject":"Hello",
    "text":"Plain text body",
    "html":"<b>HTML body</b>",
    "to":[{"email":"user@example.com","type":"to"}],
    "headers":{"Reply-To":"reply@example.com"}
  }
}' | jq .
```

- Send raw:

```
curl -s localhost:8080/messages/send-raw -H 'Content-Type: application/json' -d '{
  "key":"dev",
  "from_email":"sender@example.com",
  "to":["user@example.com"],
  "raw_message":"From: Sender <sender@example.com>\r\nTo: user@example.com\r\nSubject: Raw test\r\n\r\nRaw body\r\n"
}' | jq .
```

Notes

- This is for local development; there is no persistence across restarts.
- Attachments and inline images are supported via base64 in `attachments` and `images` arrays.
- Template sending performs simple `*|NAME|*` token replacement using `template_content` items.
- Scheduler is a best-effort background loop checking once per second.
Node send-template client (local server):

```
MC_BASE=http://localhost:8080/api/1.0 \
KEY=dev \
FROM=sender@example.com \
TO=user@example.com \
node examples/node-client/send-template.js
```
