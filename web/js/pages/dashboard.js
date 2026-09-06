window.Pages = window.Pages || {};

Pages.dashboard = {
  async render(mount) {
    mount.innerHTML = `
      <div class="page-head">
        <div>
          <div class="page-title">Your LXR environment</div>
          <div class="page-subtitle">Container infrastructure at a glance</div>
        </div>
        <div class="page-actions">
          <button class="btn btn-secondary" id="btn-refresh">${Icons.refresh()}<span class="btn-label-full">Refresh</span></button>
          <button class="btn btn-primary" id="btn-create">${Icons.plus()}<span class="btn-label-full">Create container</span></button>
        </div>
      </div>

      <div class="stat-grid" id="stat-grid">${this._skeletonStats()}</div>

      <div class="dash-grid">
        <div class="card">
          <div class="card-header">
            <span class="card-title">Running containers</span>
            <a href="/containers" data-link class="btn btn-ghost btn-sm">View all${Icons.chevronRight()}</a>
          </div>
          <div class="card-body no-pad" id="running-table"></div>
        </div>

        <div class="card">
          <div class="card-header"><span class="card-title">Quick actions</span></div>
          <div class="card-body">
            <div class="quick-actions">
              <button class="quick-action" id="qa-create">
                <span class="quick-action-icon">${Icons.plus()}</span>
                <span class="quick-action-title">Create container</span>
                <span class="quick-action-desc">Launch a new isolated container</span>
              </button>
              <button class="quick-action" id="qa-pull">
                <span class="quick-action-icon">${Icons.download()}</span>
                <span class="quick-action-title">Pull image</span>
                <span class="quick-action-desc">Fetch an image from Docker Hub</span>
              </button>
              <a class="quick-action" href="/containers" data-link>
                <span class="quick-action-icon">${Icons.container()}</span>
                <span class="quick-action-title">Open containers</span>
                <span class="quick-action-desc">Manage everything in one place</span>
              </a>
            </div>
          </div>
        </div>
      </div>

      <div class="dash-grid" style="margin-top:0;">
        <div class="card">
          <div class="card-header"><span class="card-title">Recent activity</span></div>
          <div class="card-body" id="activity-panel"></div>
        </div>
        <div class="card">
          <div class="card-header">
            <span class="card-title">Network</span>
            <a href="/network" data-link class="btn btn-ghost btn-sm">Details${Icons.chevronRight()}</a>
          </div>
          <div class="card-body" id="network-panel"></div>
        </div>
      </div>
    `;

    mount.querySelector('#btn-refresh').addEventListener('click', () => this._load(mount, true));
    mount.querySelector('#btn-create').addEventListener('click', () => openCreateContainerModal({ onCreated: () => this._load(mount, true) }));
    mount.querySelector('#qa-create').addEventListener('click', () => openCreateContainerModal({ onCreated: () => this._load(mount, true) }));
    mount.querySelector('#qa-pull').addEventListener('click', () => openPullImageModal({}));

    this._unsubActivity = AppState.subscribe('activity', () => this._renderActivity(mount));
    this._renderActivity(mount);

    await this._load(mount);
    this._loadNetwork(mount);
  },

  _skeletonStats() {
    return Array.from({ length: 4 }).map(() => `
      <div class="stat-card">
        <div class="skeleton skeleton-line" style="width:70px;"></div>
        <div class="skeleton skeleton-line" style="width:40px;height:26px;"></div>
      </div>
    `).join('');
  },

  async _load(mount, isManualRefresh = false) {
    try {
      const [online, containers] = await Promise.all([
        LXRApi.ping().catch(() => false),
        LXRApi.getContainers(),
      ]);
      AppState.setState({ runtimeOnline: online, containers, containersLoaded: true });
      this._renderStats(mount, online, containers);
      this._renderRunningTable(mount, containers.filter(c => (c.status || '').toLowerCase() === 'running'));
      if (isManualRefresh) Notify.info('Refreshed');
    } catch (err) {
      mount.querySelector('#stat-grid').innerHTML = `
        <div class="state-block" style="grid-column:1/-1;">
          <span class="state-icon">${Icons.warning()}</span>
          <div class="state-title">LXR daemon unavailable</div>
          <div class="state-desc">The web console cannot reach the LXR runtime.</div>
          <button class="btn btn-secondary" id="retry-btn">Retry</button>
        </div>`;
      mount.querySelector('#retry-btn').addEventListener('click', () => this._load(mount, true));
    }
  },

  _renderStats(mount, online, containers) {
    const running = containers.filter(c => (c.status || '').toLowerCase() === 'running').length;
    const stopped = containers.length - running;
    mount.querySelector('#stat-grid').innerHTML = `
      <div class="stat-card">
        <span class="stat-label">${Icons.pulse()} LXR runtime</span>
        <span class="stat-value ${online ? 'stat-online' : 'stat-offline'}">${online ? 'Online' : 'Offline'}</span>
        <span class="stat-meta">via GET /ping</span>
      </div>
      <div class="stat-card">
        <span class="stat-label">${Icons.container()} Total containers</span>
        <span class="stat-value">${containers.length}</span>
        <span class="stat-meta">via GET /ps/all</span>
      </div>
      <div class="stat-card">
        <span class="stat-label">${Icons.play()} Running</span>
        <span class="stat-value">${running}</span>
        <span class="stat-meta">active right now</span>
      </div>
      <div class="stat-card">
        <span class="stat-label">${Icons.stop()} Stopped</span>
        <span class="stat-value">${stopped}</span>
        <span class="stat-meta">frozen or inactive</span>
      </div>
    `;
  },

  _renderRunningTable(mount, running) {
    const el = mount.querySelector('#running-table');
    if (!running.length) {
      el.innerHTML = `
        <div class="state-block">
          <span class="state-icon">${Icons.container()}</span>
          <div class="state-title">No containers running</div>
          <div class="state-desc">Create your first LXR container to get started.</div>
          <button class="btn btn-primary" id="empty-create">${Icons.plus()}Create container</button>
        </div>`;
      el.querySelector('#empty-create').addEventListener('click', () => openCreateContainerModal({ onCreated: () => this._load(mount, true) }));
      return;
    }
    el.innerHTML = `
      <div class="table-wrap">
        <table class="data-table">
          <thead><tr><th>Name</th><th>Image</th><th>Status</th><th>IP</th><th>Ports</th><th></th></tr></thead>
          <tbody>
            ${running.map(c => `
              <tr>
                <td><a href="/containers/${encodeURIComponent(c.container_name)}" data-link class="cell-primary mono">${Utils.escapeHtml(c.container_name)}</a></td>
                <td>${Utils.escapeHtml(c.image || c.image_name || '—')}</td>
                <td><span class="badge ${Utils.statusBadgeClass(c.status)}">${Utils.statusLabel(c.status)}</span></td>
                <td class="mono">${Utils.escapeHtml(Utils.ipOnly(c.ip_address) || '—')}</td>
                <td class="mono">${(c.ports || []).join(', ') || '—'}</td>
                <td class="row-actions">
                  <button class="btn btn-ghost btn-icon btn-sm" title="Terminal" data-term="${Utils.escapeHtml(c.container_name)}">${Icons.terminal()}</button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
    Utils.qsa('[data-term]', el).forEach(btn => btn.addEventListener('click', () => LXRActions.goToTerminal(btn.dataset.term)));
  },

  _renderActivity(mount) {
    const el = mount.querySelector('#activity-panel');
    if (!el) return;
    const activity = AppState.state.activity;
    if (!activity.length) {
      el.innerHTML = `<div class="unavailable">${Icons.clock()} No local activity yet this session.</div>`;
      return;
    }
    el.innerHTML = `
      <div class="activity-list">
        ${activity.slice(0, 8).map(a => `
          <div class="activity-item">
            <span class="activity-dot ${a.type}"></span>
            <div>
              <div class="activity-text">${Utils.escapeHtml(a.text)} <span class="local-tag">Local UI</span></div>
              <div class="activity-time">${Utils.timeAgo(a.time)}</div>
            </div>
          </div>
        `).join('')}
      </div>
    `;
  },

  async _loadNetwork(mount) {
    const el = mount.querySelector('#network-panel');
    try {
      const net = await LXRApi.getNetwork();
      el.innerHTML = `
        <dl class="kv-grid">
          <dt>Bridge</dt><dd class="mono">${Utils.escapeHtml(net.bridge)}</dd>
          <dt>Bridge IP</dt><dd class="mono">${Utils.escapeHtml(net.bridge_ip)}</dd>
          <dt>Subnet</dt><dd class="mono">${Utils.escapeHtml(net.subnet)}</dd>
        </dl>
        ${LXRApi.DEMO_MODE ? `<div class="help-text" style="margin-top:10px;">Demo data</div>` : ''}
      `;
    } catch (err) {
      el.innerHTML = `<div class="unavailable">${Icons.info()} Network overview isn't exposed by the current LXR API yet.</div>`;
    }
  },
};
