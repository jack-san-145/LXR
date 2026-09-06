window.Pages = window.Pages || {};

Pages.settings = {
  async render(mount) {
    const online = AppState.state.runtimeOnline;
    mount.innerHTML = `
      <div class="page-head">
        <div>
          <div class="page-title">Settings</div>
          <div class="page-subtitle">Runtime and web console configuration</div>
        </div>
      </div>

      <div class="detail-grid">
        <div class="card">
          <div class="card-header"><span class="card-title">LXR runtime</span></div>
          <div class="card-body">
            <dl class="kv-grid">
              <dt>Daemon status</dt>
              <dd>
                <span class="badge ${online ? 'badge-running' : 'badge-danger'}">${online ? 'Online' : 'Offline'}</span>
              </dd>
              <dt>Unix socket</dt>
              <dd class="mono">/var/run/lxr.sock <span class="tag-soon" style="margin-left:6px;">daemon-internal</span></dd>
              <dt>Web API base</dt>
              <dd class="mono">${Utils.escapeHtml(window.LXR_API_BASE || '/api')}</dd>
              <dt>Bridge</dt>
              <dd class="mono">lxr0</dd>
            </dl>
          </div>
        </div>

        <div class="card">
          <div class="card-header"><span class="card-title">API surface</span></div>
          <div class="card-body no-pad">
            <div class="table-wrap">
              <table class="data-table">
                <thead><tr><th>Endpoint</th><th>Status</th></tr></thead>
                <tbody>
                  ${[
                    ['GET /ping', true], ['GET /ps', true], ['GET /ps/all', true],
                    ['POST /create', true], ['POST /start', true], ['GET /stop', true],
                    ['DELETE /kill', true], ['POST /pull_image', true],
                    ['GET /exec (PTY)', true],
                    ['GET /images', false], ['GET /network', false],
                    ['GET /container/:name/resources', false],
                    ['GET /events', false], ['WS terminal', false],
                  ].map(([ep, avail]) => `
                    <tr>
                      <td class="mono">${ep}</td>
                      <td>
                        ${avail
                          ? `<span class="badge badge-running">Available from daemon</span>`
                          : `<span class="badge badge-neutral">Not exposed by current API</span>`}
                      </td>
                    </tr>
                  `).join('')}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <div class="card" style="margin-top:var(--sp-5);">
        <div class="card-header"><span class="card-title">About this console</span></div>
        <div class="card-body">
          <p class="help-text">
            This is a self-hosted management console for LXR. There is currently no authentication layer —
            treat this deployment as trusted and local, the same way you would treat direct access to the
            Unix socket.
          </p>
        </div>
      </div>
    `;
  },
};
