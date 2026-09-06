window.Pages = window.Pages || {};

Pages.containers = {
  _all: [],
  _filters: { q: '', status: 'all', image: 'all' },

  async render(mount) {
    mount.innerHTML = `
      <div class="page-head">
        <div>
          <div class="page-title">Containers</div>
          <div class="page-subtitle">Create, start, stop, and manage LXR containers</div>
        </div>
        <div class="page-actions">
          <button class="btn btn-secondary" id="btn-refresh">${Icons.refresh()}<span class="btn-label-full">Refresh</span></button>
          <button class="btn btn-primary" id="btn-create">${Icons.plus()}<span class="btn-label-full">Create container</span></button>
        </div>
      </div>

      <div class="filters-bar">
        <div class="search-input-wrap">
          ${Icons.search()}
          <input class="input" id="f-search" placeholder="Search by name or image…" autocomplete="off" />
        </div>
        <select class="select" id="f-status">
          <option value="all">All statuses</option>
          <option value="running">Running</option>
          <option value="frozen">Frozen</option>
          <option value="stopped">Stopped</option>
        </select>
        <select class="select" id="f-image">
          <option value="all">All images</option>
        </select>
      </div>

      <div class="card">
        <div class="card-body no-pad" id="containers-body"></div>
      </div>
    `;

    mount.querySelector('#btn-refresh').addEventListener('click', () => this._load(mount, true));
    mount.querySelector('#btn-create').addEventListener('click', () => openCreateContainerModal({ onCreated: () => this._load(mount, true) }));
    mount.querySelector('#f-search').addEventListener('input', Utils.debounce((e) => {
      this._filters.q = e.target.value.trim().toLowerCase();
      this._renderList(mount);
    }, 200));
    mount.querySelector('#f-status').addEventListener('change', (e) => {
      this._filters.status = e.target.value;
      this._renderList(mount);
    });
    mount.querySelector('#f-image').addEventListener('change', (e) => {
      this._filters.image = e.target.value;
      this._renderList(mount);
    });

    this._renderSkeleton(mount);
    await this._load(mount);
  },

  _renderSkeleton(mount) {
    mount.querySelector('#containers-body').innerHTML = Array.from({ length: 4 })
      .map(() => `<div class="skeleton skeleton-row"></div>`).join('');
  },

  async _load(mount, manual = false) {
    try {
      this._all = await LXRApi.getContainers();
      AppState.setState({ containers: this._all, containersLoaded: true });
      this._populateImageFilter(mount);
      this._renderList(mount);
      if (manual) Notify.info('Refreshed');
    } catch (err) {
      mount.querySelector('#containers-body').innerHTML = `
        <div class="state-block">
          <span class="state-icon">${Icons.warning()}</span>
          <div class="state-title">LXR daemon unavailable</div>
          <div class="state-desc">The web console cannot reach the LXR runtime.</div>
          <button class="btn btn-secondary" id="retry-btn">Retry</button>
        </div>`;
      mount.querySelector('#retry-btn').addEventListener('click', () => this._load(mount, true));
    }
  },

  _populateImageFilter(mount) {
    const sel = mount.querySelector('#f-image');
    const current = sel.value;
    const images = Array.from(new Set(this._all.map(c => c.image || c.image_name).filter(Boolean))).sort();
    sel.innerHTML = `<option value="all">All images</option>` + images.map(i => `<option value="${Utils.escapeHtml(i)}">${Utils.escapeHtml(i)}</option>`).join('');
    if (images.includes(current)) sel.value = current;
  },

  _filtered() {
    const { q, status, image } = this._filters;
    return this._all.filter(c => {
      if (status !== 'all' && (c.status || '').toLowerCase() !== status) return false;
      if (image !== 'all' && (c.image || c.image_name) !== image) return false;
      if (q) {
        const hay = `${c.container_name} ${c.image || c.image_name || ''}`.toLowerCase();
        if (!hay.includes(q)) return false;
      }
      return true;
    });
  },

  _renderList(mount) {
    const list = this._filtered();
    const body = mount.querySelector('#containers-body');

    if (!this._all.length) {
      body.innerHTML = `
        <div class="state-block">
          <span class="state-icon">${Icons.container()}</span>
          <div class="state-title">No containers yet</div>
          <div class="state-desc">Create your first LXR container to get started.</div>
          <button class="btn btn-primary" id="empty-create">${Icons.plus()}Create container</button>
        </div>`;
      body.querySelector('#empty-create').addEventListener('click', () => openCreateContainerModal({ onCreated: () => this._load(mount, true) }));
      return;
    }

    if (!list.length) {
      body.innerHTML = `
        <div class="state-block">
          <span class="state-icon">${Icons.search()}</span>
          <div class="state-title">No matching containers</div>
          <div class="state-desc">Try adjusting your search or filters.</div>
        </div>`;
      return;
    }

    body.innerHTML = `
      <div class="table-wrap">
        <table class="data-table">
          <thead>
            <tr><th>Container</th><th>Image</th><th>Status</th><th>IP</th><th>PID</th><th>Ports</th><th></th></tr>
          </thead>
          <tbody>
            ${list.map(c => this._row(c)).join('')}
          </tbody>
        </table>
      </div>
      <div class="record-cards" style="padding: var(--sp-4);">
        ${list.map(c => this._recordCard(c)).join('')}
      </div>
    `;
    this._wireActions(mount, body);
  },

  _actionButtons(c) {
    const actions = LXRActions.actionsFor(c.status);
    const map = {
      start: `<button class="btn btn-ghost btn-icon btn-sm" data-act="start" data-name="${Utils.escapeHtml(c.container_name)}" title="Start">${Icons.play()}</button>`,
      stop: `<button class="btn btn-ghost btn-icon btn-sm" data-act="stop" data-name="${Utils.escapeHtml(c.container_name)}" title="Stop / freeze">${Icons.stop()}</button>`,
      terminal: `<button class="btn btn-ghost btn-icon btn-sm" data-act="terminal" data-name="${Utils.escapeHtml(c.container_name)}" title="Terminal">${Icons.terminal()}</button>`,
      open: `<button class="btn btn-ghost btn-icon btn-sm" data-act="open" data-name="${Utils.escapeHtml(c.container_name)}" title="Open development environment">${Icons.code()}</button>`,
      kill: `<button class="btn btn-ghost btn-icon btn-sm" data-act="kill" data-name="${Utils.escapeHtml(c.container_name)}" title="Kill">${Icons.kill()}</button>`,
    };
    return actions.map(a => map[a]).join('');
  },

  _row(c) {
    return `
      <tr>
        <td>
          <a href="/containers/${encodeURIComponent(c.container_name)}" data-link class="cell-primary mono">${Utils.escapeHtml(c.container_name)}</a>
        </td>
        <td>${Utils.escapeHtml(c.image || c.image_name || '—')}</td>
        <td><span class="badge ${Utils.statusBadgeClass(c.status)}">${Utils.statusLabel(c.status)}</span></td>
        <td class="mono">${Utils.escapeHtml(Utils.ipOnly(c.ip_address) || '—')}</td>
        <td class="mono">${c.pid || '—'}</td>
        <td class="mono">${(c.ports || []).join(', ') || '—'}</td>
        <td class="row-actions">${this._actionButtons(c)}</td>
      </tr>
    `;
  },

  _recordCard(c) {
    return `
      <div class="record-card">
        <div class="record-card-head">
          <div>
            <a href="/containers/${encodeURIComponent(c.container_name)}" data-link class="record-card-title mono">${Utils.escapeHtml(c.container_name)}</a>
            <div class="record-card-sub">${Utils.escapeHtml(c.image || c.image_name || '—')}</div>
          </div>
          <span class="badge ${Utils.statusBadgeClass(c.status)}">${Utils.statusLabel(c.status)}</span>
        </div>
        <div class="record-field"><span class="record-field-label">IP</span><span class="mono">${Utils.escapeHtml(Utils.ipOnly(c.ip_address) || '—')}</span></div>
        <div class="record-field"><span class="record-field-label">PID</span><span class="mono">${c.pid || '—'}</span></div>
        <div class="record-field"><span class="record-field-label">Ports</span><span class="mono">${(c.ports || []).join(', ') || '—'}</span></div>
        <div class="record-card-actions">${this._actionButtons(c)}</div>
      </div>
    `;
  },

  _wireActions(mount, body) {
    Utils.qsa('[data-act]', body).forEach(btn => {
      const name = btn.dataset.name;
      const act = btn.dataset.act;
      btn.addEventListener('click', async () => {
        if (act === 'start') return LXRActions.doStart(name, { onDone: () => this._load(mount, true) });
        if (act === 'stop') return LXRActions.doStop(name, { onDone: () => this._load(mount, true) });
        if (act === 'kill') return LXRActions.doKill(name, { onDone: () => this._load(mount, true) });
        if (act === 'terminal') return LXRActions.goToTerminal(name);
        if (act === 'open') {
          const c = this._all.find(x => x.container_name === name);
          if (c) LXRActions.openCode(c);
        }
      });
    });
  },
};
