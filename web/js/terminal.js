/**
 * terminal.js
 *
 * TerminalClient is a transport-agnostic abstraction. The current daemon's
 * /exec hijacks a raw HTTP connection to attach a PTY via nsenter — that
 * cannot be driven from browser fetch(), and the future web server is
 * expected to expose a WebSocket PTY at /api/containers/:name/terminal.
 * Until that lands, connect() reports the transport as unavailable rather
 * than faking a shell. Swap the connect()/send() bodies here once the
 * WebSocket endpoint exists; nothing else in the app needs to change.
 */
class TerminalClient {
  constructor({ onData, onStatus, onClose } = {}) {
    this.onData = onData || (() => {});
    this.onStatus = onStatus || (() => {});
    this.onClose = onClose || (() => {});
    this.socket = null;
    this.connected = false;
  }

  connect(containerName) {
    this.onStatus('connecting');

    if (window.LXRApi && window.LXRApi.DEMO_MODE) {
      this.connected = true;
      this.onStatus('connected');
      this._writeLine(`Connected to ${containerName} (demo session — no real PTY attached).`);
      this._prompt();
      return;
    }

    // Real transport seam: once /api/containers/:name/terminal (WebSocket)
    // exists, this becomes:
    //   const proto = location.protocol === 'https:' ? 'wss' : 'ws';
    //   this.socket = new WebSocket(`${proto}://${location.host}${API_BASE}/containers/${containerName}/terminal`);
    //   this.socket.onopen = () => { this.connected = true; this.onStatus('connected'); };
    //   this.socket.onmessage = (ev) => this.onData(ev.data);
    //   this.socket.onclose = () => { this.connected = false; this.onStatus('closed'); this.onClose(); };
    try {
      window.LXRApi.connectTerminal(containerName);
    } catch (err) {
      this.onStatus('unavailable');
      this.onData(
        `\r\nInteractive terminal transport is not yet available over HTTP/WebSocket.\r\n` +
        `The daemon currently serves /exec by hijacking a raw connection and attaching a PTY via nsenter,\r\n` +
        `which a browser cannot drive directly. This view is wired up and will connect automatically\r\n` +
        `once the web server exposes a WebSocket terminal endpoint.\r\n`
      );
    }
  }

  send(data) {
    if (this.socket && this.connected) {
      this.socket.send(data);
      return;
    }
    if (window.LXRApi && window.LXRApi.DEMO_MODE && this.connected) {
      this._handleDemoInput(data);
    }
  }

  close() {
    if (this.socket) this.socket.close();
    this.connected = false;
  }

  // ---- demo-mode line editing, only active when LXR_DEMO_MODE is true ----
  _writeLine(text) { this.onData(text + '\r\n'); }
  _prompt() { this.onData('\r\n$ '); }
  _handleDemoInput(ch) {
    if (ch === '\r') {
      this._writeLine('');
      this._prompt();
    } else {
      this.onData(ch);
    }
  }
}

window.TerminalClient = TerminalClient;
