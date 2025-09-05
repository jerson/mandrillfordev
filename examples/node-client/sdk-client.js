// Node example using the official Mailchimp Transactional (Mandrill) SDK.
// Note: As of this example, the SDK targets https://mandrillapp.com by default
// and does not publicly expose a base URL override. For local testing against
// this dev server, you have two options:
// 1) Use the HTTP client example (http-client.js) which directly posts to the dev server.
// 2) Temporarily override DNS (e.g., hosts file) to point mandrillapp.com to your
//    dev server and terminate TLS via a local proxy (advanced).
// If your SDK version exposes a base URL override, set MC_BASE to your server.

import transactional from '@mailchimp/mailchimp_transactional';

const API_KEY = process.env.KEY || 'dev';
const MC_BASE = process.env.MC_BASE || 'http://localhost:8080/api/1.0'; // e.g., http://localhost:8080/api/1.0
const FROM = process.env.FROM || 'sender@example.com';
const TO = process.env.TO || 'user@example.com';

async function main() {
  const payload = {
    key: API_KEY,
    message: {
      from_email: FROM,
      subject: 'SDK send test',
      text: 'Hello from SDK client',
      to: [{ email: TO, type: 'to' }]
    }
  };

  // If MC_BASE looks like a local/dev override, call the dev server directly.
  if (MC_BASE && /^https?:\/\//i.test(MC_BASE)) {
    const url = `${MC_BASE.replace(/\/$/, '')}/messages/send.json`;
    console.log('Using local dev server via MC_BASE:', url);
    try {
      const res = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      const text = await res.text();
      console.log('SDK response:', JSON.stringify({ status: res.status, body: safeParse(text) }, null, 2));
      if (!res.ok) process.exit(1);
      return;
    } catch (err) {
      console.error('SDK error (local dev):', err);
      process.exit(1);
    }
  }

  // Fallback to official SDK (targets mandrillapp.com)
  const client = transactional(API_KEY);
  try {
    const res = await client.messages.send(payload);
    console.log('SDK response:', JSON.stringify(res, null, 2));
  } catch (err) {
    console.error('SDK error (sdk host):', err);
    process.exit(1);
  }
}

function safeParse(s) {
  try { return JSON.parse(s); } catch { return s; }
}

main();
