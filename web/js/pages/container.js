window.Pages = window.Pages || {};

Pages.container = {
  async render(mount, name) {
    mount.innerHTML = `<div class="skeleton" style="height:200px;border-radius:var(--radius-lg);"></div>`;

    let container;
    try {
      container = await LXRApi.getContainer(name);
    } catch (err) {
      mount.innerHTML = `
        <div class="state-block">
          <span class="state-icon">${Icons.warning()}</span>
          <div class="state-title">Could not load this container</div>
          <div class="state-desc">${Utils.escapeHtml(err.message)}</div>
        </div>`;
      return;
    }

    if (!container) {
      mount.innerHTML = `
        <div class="state-block">
          <span class="state-icon">${Icons.container()}</span>
          <div class="state-title">Container not found</div>
          <div class="state-desc">“${Utils.escapeHtml(name)}” doesn't exist, or the daemon no longer knows about it.</div>
          <a href="/containers" data-link class="btn btn-primary">Back to containers</a>
        </div>`;
      return;
    }

    this._draw(mount, container);
  },

  _draw(mount, c) {
    const actions = LXRActions.actionsFor(c.status);
    mount.innerHTML = `
      <div class="detail-head">
        <div>
          <div class="detail-title-row">
            <span class="detail-title">${Utils.escapeHtml(c.container_name)}</span>
            <span class="badge ${Utils.statusBadgeClass(c.status)}">${Utils.statusLabel(c.status)}</span>
          </div>
          <div class="detail-sub">Image: ${Utils.escapeHtml(c.image || c.image_name || '—')}</div>
        </div>
        <div class="detail-actions">
          ${actions.includes('terminal') ? `<button class="btn btn-secondary" data-act="terminal">${Icons.terminal()}Terminal</button>` : ''}
          ${actions.includes('open') ? `<button class="btn btn-secondary" data-act="open">${Icons.code()}Open code</button>` : ''}
          ${actions.includes('start') ? `<button class="btn btn-primary" data-act="start">${Icons.play()}Start</button>` : ''}
          ${actions.includes('stop') ? `<button class="btn btn-secondary" data-act="stop">${Icons.stop()}Stop</button>` : ''}
          ${actions.includes('kill') ? `<button class="btn btn-danger-outline" data-act="kill">${Icons.kill()}Kill</button>` : ''}
        </div>
      </div>

      <div class="detail-grid">
        <div class="card">
          <div class="card-header"><span class="card-title">Overview</span></div>
          <div class="card-body">
            <dl class="kv-grid">
              <dt>Container ID</dt><dd>${copyable(c.container_id)}</dd>
              <dt>Image</dt><dd>${Utils.escapeHtml(c.image || c.image_name || '—')}</dd>
              <dt>PID</dt><dd>${copyable(c.pid)}</dd>
              <dt>IP address</dt><dd>${copyable(Utils.ipOnly(c.ip_address))}</dd>
              <dt>Bridge</dt><dd class="mono">${Utils.escapeHtml(c.bridge || 'lxr0')}</dd>
              <dt>Ports</dt><dd class="mono">${(c.ports || []).join(', ') || '—'}</dd>
              <dt>Status</dt><dd><span class="badge ${Utils.statusBadgeClass(c.status)}">${Utils.statusLabel(c.status)}</span></dd>
            </dl>
          </div>
        </div>

        <div class="card">
          <div class="card-header"><span class="card-title">Network</span></div>
          <div class="card-body" id="network-card">
            <dl class="kv-grid">
              <dt>IP</dt><dd>${copyable(c.ip_address)}</dd>
              <dt>Bridge</dt><dd class="mono">${Utils.escapeHtml(c.bridge || 'lxr0')}</dd>
              <dt>Container veth</dt><dd id="con-veth"></dd>
              <dt>Bridge veth</dt><dd id="br-veth"></dd>
            </dl>
          </div>
        </div>
      </div>

      <div class="card" style="margin-top: var(--sp-5);">
        <div class="card-header"><span class="card-title">Resources</span></div>
        <div class="card-body">
          <div class="resource-placeholder-grid">
            <div class="resource-block">
              <div class="resource-block-label">${Icons.cpu()} CPU</div>
              <div class="unavailable">Not available</div>
            </div>
            <div class="resource-block">
              <div class="resource-block-label">${Icons.memory()} Memory</div>
              <div class="unavailable">Not available</div>
            </div>
            <div class="resource-block">
              <div class="resource-block-label">${Icons.pulse()} Processes</div>
              <div class="unavailable">Not available</div>
            </div>
          </div>
          <p class="help-text" style="margin-top:var(--sp-4);">
            Live resource metrics aren't exposed by the current LXR API. The daemon configures cgroup limits
            (256 MB memory, 50% CPU, 100 PIDs) at creation time, but does not yet report live usage.
          </p>
        </div>
      </div>

      <div class="card" style="margin-top: var(--sp-5);">
        <div class="card-header"><span class="card-title">Development environment</span></div>
        <div class="card-body">
          <div class="dev-env-box">
            <div class="dev-env-info">
              <span class="dev-env-icon">${Icons.code()}</span>
              <div>
                <div style="font-weight:700;">code-server</div>
                <div class="help-text">Port ${(c.ports && c.ports[0]) || 9000}</div>
              </div>
            </div>
            ${actions.includes('open')
              ? `<button class="btn btn-primary" data-act="open">${Icons.external()}Open development environment</button>`
              : `<span class="unavailable">${Icons.info()} Start the container to open its environment</span>`}
          </div>
        </div>
      </div>
    `;

    Utils.qsa('[data-act]', mount).forEach(btn => {
      btn.addEventListener('click', async () => {
        const act = btn.dataset.act;
        if (act === 'terminal') return LXRActions.goToTerminal(c.container_name);
        if (act === 'open') return LXRActions.openCode(c);
        if (act === 'start') return LXRActions.doStart(c.container_name, { onDone: () => this.render(mount, c.container_name) });
        if (act === 'stop') return LXRActions.doStop(c.container_name, { onDone: () => this.render(mount, c.container_name) });
        if (act === 'kill') {
          await LXRActions.doKill(c.container_name, { onDone: () => Router.navigate('/containers') });
        }
      });
    });

    // veth fields are part of the model but the current /ps and /ps/all
    // responses documented in the spec don't include them — show clearly.
    const conVeth = mount.querySelector('#con-veth');
    const brVeth = mount.querySelector('#br-veth');
    if (c.con_veth) conVeth.innerHTML = `<span class="mono">${Utils.escapeHtml(c.con_veth)}</span>`;
    else conVeth.innerHTML = `<span class="unavailable">${Icons.info()} Not available</span>`;
    if (c.br_veth) brVeth.innerHTML = `<span class="mono">${Utils.escapeHtml(c.br_veth)}</span>`;
    else brVeth.innerHTML = `<span class="unavailable">${Icons.info()} Not available</span>`;
  },
};
