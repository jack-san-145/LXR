window.Pages = window.Pages || {};

Pages.terminal = {
  async render(mount, name) {
    mount.innerHTML = `
      <div class="page-head">
        <div>
          <div class="page-title mono">${Utils.escapeHtml(name)}</div>
          <div class="page-subtitle">Interactive shell via nsenter into the container's namespaces</div>
        </div>
        <div class="page-actions">
          <a href="/containers/${encodeURIComponent(name)}" data-link class="btn btn-secondary">Back to container</a>
        </div>
      </div>

      <div class="term-shell">
        <div class="term-titlebar">
          <div class="term-titlebar-left">
            <span class="term-dots"><span></span><span></span><span></span></span>
            <span class="term-name">${Utils.escapeHtml(name)}</span>
          </div>
          <span class="term-status" id="term-status">
            <span class="spinner"></span> Connecting…
          </span>
        </div>
        <div class="term-body" id="term-body" tabindex="0"></div>
        <div class="term-hint-bar">
          <span><kbd>Ctrl</kbd>+<kbd>C</kbd> interrupt</span>
          <span><kbd>↑</kbd>/<kbd>↓</kbd> history</span>
          <span>Click the terminal to focus it</span>
        </div>
      </div>
    `;

    this._history = [];
    this._historyPos = -1;
    this._lineBuffer = '';
    const body = mount.querySelector('#term-body');
    const statusEl = mount.querySelector('#term-status');

    const client = new TerminalClient({
      onStatus: (status) => this._setStatus(statusEl, status),
      onData: (chunk) => this._write(body, chunk),
      onClose: () => this._setStatus(statusEl, 'closed'),
    });

    this._client = client;
    this._currentLineEl = null;

    body.addEventListener('keydown', (e) => this._handleKey(e, client, body));
    body.addEventListener('click', () => body.focus());

    // Kick off connection
    this._write(body, '');
    client.connect(name);

    this._cleanup = () => client.close();
    document.addEventListener('lxr:navigated', this._cleanup, { once: true });
  },

  _setStatus(el, status) {
    const map = {
      connecting: { text: 'Connecting…', cls: '' },
      connected: { text: 'Connected', cls: 'online' },
      unavailable: { text: 'Transport unavailable', cls: 'offline' },
      closed: { text: 'Disconnected', cls: 'offline' },
    };
    const s = map[status] || map.connecting;
    const spinner = status === 'connecting' ? '<span class="spinner"></span>' : `<span class="runtime-dot ${s.cls || 'pending'}"></span>`;
    el.innerHTML = `${spinner} ${s.text}`;
  },

  _write(body, text) {
    if (!text) return;
    // Render raw text as a single flowing block; CR is treated as newline
    // for this minimal terminal (no full ANSI/cursor-addressing support).
    const normalized = text.replace(/\r\n/g, '\n').replace(/\r/g, '\n');
    body.appendChild(document.createTextNode(normalized));
    body.scrollTop = body.scrollHeight;
  },

  _handleKey(e, client, body) {
    e.preventDefault();
    if (e.ctrlKey && e.key.toLowerCase() === 'c') {
      client.send('\x03');
      this._write(body, '^C');
      return;
    }
    if (e.key === 'Enter') {
      this._history.push(this._lineBuffer);
      this._historyPos = this._history.length;
      client.send('\r');
      this._lineBuffer = '';
      return;
    }
    if (e.key === 'Backspace') {
      if (this._lineBuffer.length) {
        this._lineBuffer = this._lineBuffer.slice(0, -1);
        client.send('\b \b');
      }
      return;
    }
    if (e.key === 'ArrowUp') {
      if (this._historyPos > 0) {
        this._historyPos -= 1;
        this._lineBuffer = this._history[this._historyPos] || '';
        client.send(this._lineBuffer);
      }
      return;
    }
    if (e.key === 'ArrowDown') {
      if (this._historyPos < this._history.length - 1) {
        this._historyPos += 1;
        this._lineBuffer = this._history[this._historyPos] || '';
        client.send(this._lineBuffer);
      }
      return;
    }
    if (e.key.length === 1) {
      this._lineBuffer += e.key;
      client.send(e.key);
    }
  },
};
