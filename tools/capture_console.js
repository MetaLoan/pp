const fs = require('fs');
const { chromium } = require('playwright');

(async () => {
  const out = fs.createWriteStream('capture_console.log', { flags: 'a' });
  out.write('=== Capture started ' + new Date().toISOString() + " ===\n");
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext();
  const page = await context.newPage();
  // Inject a fake auth token and user so the SPA routes to /trade
  await context.addInitScript(() => {
    try {
      localStorage.setItem('token', 'dev-token');
      localStorage.setItem('user', JSON.stringify({ id: 1, email: 'dev@local', name: 'dev' }));
    } catch (e) {
      // ignore
    }
  });

  page.on('console', (msg) => {
    const text = `[PAGE console] ${msg.type().toUpperCase()}: ${msg.text()}`;
    out.write(new Date().toISOString() + ' ' + text + '\n');
    console.log(text);
  });

  page.on('pageerror', (err) => {
    const text = `[PAGE error] ${err}`;
    out.write(new Date().toISOString() + ' ' + text + '\n');
    console.error(text);
  });

  page.on('requestfailed', (req) => {
    const text = `[REQUEST FAILED] ${req.method()} ${req.url()} - ${req.failure().errorText}`;
    out.write(new Date().toISOString() + ' ' + text + '\n');
    console.log(text);
  });

  const url = process.argv[2] || 'http://localhost:3000';
  console.log('Opening', url);
  // Ensure we navigate to /trade to reach TradeView
  const tradeUrl = url.endsWith('/') ? url + 'trade' : url + '/trade';
  await page.goto(tradeUrl, { waitUntil: 'networkidle' });

  // Wait and collect logs for a duration
  const durationMs = parseInt(process.argv[3], 10) || 15000; // default 15s
  console.log(`Collecting console for ${durationMs}ms ...`);
  await new Promise((res) => setTimeout(res, durationMs));

  out.write('=== Capture ended ' + new Date().toISOString() + " ===\n\n");
  await browser.close();
  out.end();
})();
