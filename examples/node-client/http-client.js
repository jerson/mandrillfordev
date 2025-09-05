import fetch from 'node-fetch';

const API_BASE = process.env.API_BASE || 'http://localhost:8080';
const KEY = process.env.KEY || 'dev';
const FROM = process.env.FROM || 'sender@example.com';
const TO = process.env.TO || 'user@example.com';
const REPLY_TO = process.env.REPLY_TO || 'reply@example.com';

async function main() {
  await waitFor(`${API_BASE}/healthz`, 30000);
  const url = `${API_BASE}/api/1.0/messages/send.json`;
  const payload = {
    key: KEY,
    message: {
      from_email: FROM,
      from_name: 'SDK-ish Client',
      subject: 'Mandrill SDK-compatible path test',
      text: 'Hello via /api/1.0/messages/send.json',
      html: '<p>Hello via <b>/api/1.0/messages/send.json</b></p>',
      to: [{ email: TO, type: 'to' }],
      headers: { 'Reply-To': REPLY_TO }
    }
  };
  const res = await fetch(url, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) });
  const body = await res.text();
  console.log('Status:', res.status, res.statusText);
  console.log('Response:', body);
}

async function waitFor(url, timeoutMs) {
  const end = Date.now() + timeoutMs;
  while (Date.now() < end) {
    try {
      const r = await fetch(url);
      if (r.ok) { return; }
    } catch {}
    await new Promise(r => setTimeout(r, 500));
  }
  throw new Error(`timeout waiting for ${url}`);
}

main().catch(err => { console.error(err); process.exit(1); });

