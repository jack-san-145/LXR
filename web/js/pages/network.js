window.Pages = window.Pages || {};

Pages.network = {
  async render(mount) {
    mount.innerHTML = `
      <div class="page-head">
        <div>
          <div class="page-title">Network</div>
          <div class="page-subtitle">lxr0 bridge, IP allocation, and connected containers</div>
        </div>
      </div>

      <div class="detail-grid">
        <div class="card">
          <div class="card-header"><span class="card-title">LXR network</span></div>
          <div class="card-body" id="network-overview">
            <div class="skeleton skeleton-line"></div>
            <div class="skeleton skeleton-line"></div>
            <div class="skeleton skeleton-line"></div>
          </div>
        </div>
        <div class="card">
          <div class="card-header"><span class="card-title">Bridge topology</span></div>
          <div class="card-body">
            <div class="network-topo">${this._topologySvg()}</div>
          </div>
        </div>
      </div>

      <div class="card" style="margin-top:var(--sp-5);">
        <div class="card-header"><span class="card-title">Container network addresses</span></div>
        <div class="card-body no-pad" id="network-containers"></div>
      </div>
    `;

    await this._loadOverview(mount);
    await this._loadContainers(mount);
  },

  _topologySvg() {
    return `
      <svg width="100%" height="140" viewBox="0 0 420 140" fill="none" aria-hidden="true">
        <rect x="8" y="55" width="70" height="30" rx="5" stroke="var(--border-strong)" stroke-width="1.5"/>
        <text x="43" y="74" text-anchor="middle" font-size="10" fill="var(--muted)" font-family="var(--font-mono)">container</text>
        <line x1="78" y1="70" x2="150" y2="70" stroke="var(--border-strong)" stroke-width="1.5"/>
        <text x="114" y="62" text-anchor="middle" font-size="9" fill="var(--faint)">veth</text>

        <rect x="150" y="40" width="120" height="60" rx="6" stroke="var(--primary)" stroke-width="1.5"/>
        <text x="210" y="74" text-anchor="middle" font-size="12" fill="var(--primary-strong)" font-weight="700" font-family="var(--font-mono)">lxr0</text>

        <line x1="270" y1="70" x2="342" y2="70" stroke="var(--border-strong)" stroke-width="1.5"/>
        <rect x="342" y="55" width="70" height="30" rx="5" stroke="var(--border-strong)" stroke-width="1.5"/>
        <text x="377" y="74" text-anchor="middle" font-size="10" fill="var(--muted)" font-family="var(--font-mono)">host</text>
      </svg>
    `;
  },

  async _loadOverview(mount) {
    const el = mount.querySelector('#network-overview');
    try {
      const net = await LXRApi.getNetwork();
      const used = net.used ?? null;
      const total = net.usable_hosts ?? null;
      const pct = (used !== null && total) ? Utils.clamp((used / total) * 100, 0.5, 100) : null;
      el.innerHTML = `
        <dl class="kv-grid">
          <dt>Bridge</dt><dd class="mono">${Utils.escapeHtml(net.bridge)}</dd>
          <dt>Bridge IP</dt><dd class="mono">${Utils.escapeHtml(net.bridge_ip)}</dd>
          <dt>Subnet</dt><dd class="mono">${Utils.escapeHtml(net.subnet)}</dd>
          <dt>Usable hosts</dt><dd class="mono">${total !== null ? total.toLocaleString() : '—'}</dd>
        </dl>
        ${pct !== null ? `
          <div class="ip-bar"><div class="ip-bar-used" style="width:${pct}%"></div></div>
          <div class="network-legend">
            <span><span class="dot" style="background:var(--primary)"></span>Used: ${used}</span>
            <span><span class="dot" style="background:var(--surface-active)"></span>Available: ${total - used}</span>
          </div>
        ` : ''}
        ${LXRApi.DEMO_MODE ? `<div class="help-text" style="margin-top:10px;">Demo data</div>` : ''}
      `;
    } catch (err) {
      el.innerHTML = `
        <div class="unavailable" style="margin-bottom:8px;">${Icons.info()} A dedicated network endpoint isn't exposed by the current LXR API yet.</div>
        <p class="help-text">
          Subnet, gateway, and IP-pool totals will appear here automatically once the daemon exposes
          <code class="mono">GET /api/network</code>. Container IP addresses below are derived from
          <code class="mono">GET /ps/all</code> in the meantime.
        </p>
      `;
    }
  },

  async _loadContainers(mount) {
    const el = mount.querySelector('#network-containers');
    try {
      const containers = await LXRApi.getContainers();
      const withIps = containers.filter(c => c.ip_address);
      if (!withIps.length) {
        el.innerHTML = `<div class="state-block"><div class="state-desc">No containers with an assigned IP address.</div></div>`;
        return;
      }
      el.innerHTML = `
        <div class="table-wrap">
          <table class="data-table">
            <thead><tr><th>Container</th><th>Status</th><th>IP address</th></tr></thead>
            <tbody>
              ${withIps.map(c => `
                <tr>
                  <td><a href="/containers/${encodeURIComponent(c.container_name)}" data-link class="cell-primary mono">${Utils.escapeHtml(c.container_name)}</a></td>
                  <td><span class="badge ${Utils.statusBadgeClass(c.status)}">${Utils.statusLabel(c.status)}</span></td>
                  <td>${copyable(c.ip_address)}</td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      `;
    } catch (err) {
      el.innerHTML = `<div class="state-block"><div class="state-desc">${Utils.escapeHtml(err.message)}</div></div>`;
    }
  },
};
