// Node example for send-template using the official SDK when possible.
// If MC_BASE is set (e.g., http://localhost:8080/api/1.0), it calls the local
// dev server directly; otherwise it uses the SDK which targets mandrillapp.com.

import transactional from '@mailchimp/mailchimp_transactional';

const API_KEY = process.env.KEY || 'dev';
const MC_BASE = process.env.MC_BASE || 'http://localhost:8080/api/1.0';
const FROM = process.env.FROM || 'sender@example.com';
const TO = process.env.TO || 'user@example.com';
const REPLY_TO = process.env.REPLY_TO || 'reply@example.com';

async function main() {
  const payload = {
    key: API_KEY,
    template_name: 'welcome',
    template_content: [
      { name: 'NAME', content: 'Friend' },
      { name: 'FEATURE', content: 'local Mandrill dev' }
    ],
    message: {
      from_email: FROM,
      from_name: 'Template Client',
      subject: 'Template test for *|NAME|*',
      text: 'Hello *|NAME|*, welcome to *|FEATURE|*!',
      html: '<p>Hello <b>*|NAME|*</b>, welcome to <i>*|FEATURE|*</i>!</p>',
      to: [{ email: TO, type: 'to' }],
      headers: { 'Reply-To': REPLY_TO },
      tags: ['example', 'send-template']
    }
  };

  // Prefer SDK unless MC_BASE explicitly points to a local dev server
  if (MC_BASE && /^https?:\/\//i.test(MC_BASE)) {
    const healthz = MC_BASE.replace(/\/$/, '').replace(/\/api\/1\.0$/, '') + '/healthz';
    await waitFor(healthz, 30000);
    const url = `${MC_BASE.replace(/\/$/, '')}/messages/send-template.json`;
    console.log('Using local dev server via MC_BASE:', url);
    const res = await fetch(url, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) });
    const body = await res.text();
    console.log('Status:', res.status, res.statusText);
    console.log('Response:', safeParse(body));
    if (!res.ok) process.exit(1);
    return;
  }

  const client = transactional(API_KEY);
  try {
    const res = await client.messages.sendTemplate(payload);
    console.log('SDK response:', JSON.stringify(res, null, 2));
  } catch (err) {
    console.error('SDK error (sdk host):', err);
    process.exit(1);
  }
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

function safeParse(s) {
  try { return JSON.parse(s); } catch { return s; }
}

main().catch(err => { console.error(err); process.exit(1); });
