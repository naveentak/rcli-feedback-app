// Tabs
const tabs = document.querySelectorAll('.tab');
const submitPanel = document.getElementById('submit-panel');
const browsePanel = document.getElementById('browse-panel');

tabs.forEach(tab => {
  tab.addEventListener('click', () => {
    tabs.forEach(t => t.classList.remove('active'));
    tab.classList.add('active');
    const isSubmit = tab.dataset.tab === 'submit';
    submitPanel.classList.toggle('hidden', !isSubmit);
    browsePanel.classList.toggle('hidden', isSubmit);
    if (!isSubmit) loadTickets();
  });
});

// Submit form
const submitForm = document.getElementById('submit-form');
const submitBtn = document.getElementById('submit-btn');
const submitResult = document.getElementById('submit-result');

submitForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  submitBtn.disabled = true;
  submitBtn.textContent = 'Submitting…';
  submitResult.classList.add('hidden');

  const data = Object.fromEntries(new FormData(submitForm));

  try {
    const res = await fetch('/api/v1/feedback', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    const body = await res.json();

    if (!res.ok) {
      throw new Error(body.error || 'Submission failed');
    }

    submitResult.className = 'submit-result success';
    submitResult.innerHTML = `Submitted <strong>#${body.number}</strong>: ${escapeHtml(body.title)} — <a href="${body.url}" target="_blank">View on GitHub</a>`;
    submitResult.classList.remove('hidden');
    submitForm.reset();
  } catch (err) {
    submitResult.className = 'submit-result error';
    submitResult.textContent = err.message;
    submitResult.classList.remove('hidden');
  } finally {
    submitBtn.disabled = false;
    submitBtn.textContent = 'Submit Feedback';
  }
});

// Browse tickets
const ticketsEl = document.getElementById('tickets');
const detailEl = document.getElementById('detail');
const appFilter = document.getElementById('app-filter');
const statusFilter = document.getElementById('status-filter');
const refreshBtn = document.getElementById('refresh-btn');

async function loadTickets() {
  const app = appFilter.value;
  const status = statusFilter.value;
  const params = new URLSearchParams({ app, status, state: 'open' });

  ticketsEl.innerHTML = '<p class="empty">Loading…</p>';
  detailEl.classList.add('hidden');

  try {
    const res = await fetch(`/api/v1/feedback?${params}`);
    const tickets = await res.json();

    if (!tickets.length) {
      ticketsEl.innerHTML = '<p class="empty">No tickets found.</p>';
      return;
    }

    ticketsEl.innerHTML = tickets.map(t => `
      <div class="ticket" data-number="${t.number}">
        <div class="ticket-header">
          <span class="ticket-number">#${t.number}</span>
          <span class="badge badge-${t.type}">${t.type}</span>
          <span class="badge badge-app">${t.app}</span>
        </div>
        <div class="ticket-title">${escapeHtml(t.title)}</div>
      </div>
    `).join('');

    document.querySelectorAll('.ticket').forEach(el => {
      el.addEventListener('click', () => showDetail(el.dataset.number));
    });
  } catch {
    ticketsEl.innerHTML = '<p class="empty">Failed to load tickets.</p>';
  }
}

async function showDetail(number) {
  const res = await fetch(`/api/v1/feedback/${number}`);
  const t = await res.json();

  detailEl.classList.remove('hidden');
  detailEl.innerHTML = `
    <h2>#${t.number} ${escapeHtml(t.title)}</h2>
    <div class="meta">
      ${t.app} · ${t.type} · ${t.status || t.state}
      · <a href="${t.url}" target="_blank">View on GitHub</a>
    </div>
    <div class="body">${escapeHtml(t.body)}</div>
  `;
  detailEl.scrollIntoView({ behavior: 'smooth' });
}

function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

refreshBtn.addEventListener('click', loadTickets);
appFilter.addEventListener('change', loadTickets);
statusFilter.addEventListener('change', loadTickets);