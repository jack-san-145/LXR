window.Pages = window.Pages || {};

Pages.images = {
  async render(mount) {
    mount.innerHTML = `
      <div class="page-head">
        <div>
          <div class="page-title">Images</div>
          <div class="page-subtitle">Pull and manage container images from Docker Hub</div>
        </div>
        <div class="page-actions">
          <button class="btn btn-primary" id="btn-pull">${Icons.download()}<span class="btn-label-full">Pull image</span></button>
        </div>
      </div>
      <div class="card">
        <div class="card-header"><span class="card-title">Local images</span></div>
        <div class="card-body" id="images-body">
          <div class="skeleton skeleton-row"></div>
        </div>
      </div>
    `;

    mount.querySelector('#btn-pull').addEventListener('click', () => openPullImageModal({ onPulled: () => this._load(mount) }));
    await this._load(mount);
  },

  async _load(mount) {
    const body = mount.querySelector('#images-body');
    try {
      const images = await LXRApi.getImages();
      if (!images.length) {
        body.innerHTML = this._empty();
        this._wireEmpty(mount, body);
        return;
      }
      body.innerHTML = `
        <div class="image-grid">
          ${images.map(img => `
            <div class="image-card">
              <div class="image-card-head">
                <span class="image-icon">${Icons.image()}</span>
                <div>
                  <div class="image-name">${Utils.escapeHtml(img.name)}:${Utils.escapeHtml(img.tag || 'latest')}</div>
                  <div class="image-meta">Pulled ${Utils.timeAgo(img.pulled_at)}</div>
                </div>
              </div>
              <button class="btn btn-secondary btn-sm btn-block" data-create-from="${Utils.escapeHtml(img.name)}">
                ${Icons.plus()}Create container
              </button>
            </div>
          `).join('')}
        </div>
        ${LXRApi.DEMO_MODE ? `<div class="help-text" style="margin-top:var(--sp-4);">Demo data</div>` : ''}
      `;
      Utils.qsa('[data-create-from]', body).forEach(btn => {
        btn.addEventListener('click', () => openCreateContainerModal({}));
      });
    } catch (err) {
      if (err instanceof LXRApi.UnsupportedError) {
        body.innerHTML = `
          <div class="state-block">
            <span class="state-icon">${Icons.info()}</span>
            <div class="state-title">Image listing isn't available yet</div>
            <div class="state-desc">
              The current LXR daemon doesn't expose <code class="mono">GET /images</code>. This page is wired
              to render local images automatically the moment that endpoint ships. You can still pull new
              images below — pulling works today via <code class="mono">POST /pull_image</code>.
            </div>
          </div>`;
      } else {
        body.innerHTML = `
          <div class="state-block">
            <span class="state-icon">${Icons.warning()}</span>
            <div class="state-title">Could not load images</div>
            <div class="state-desc">${Utils.escapeHtml(err.message)}</div>
          </div>`;
      }
    }
  },

  _empty() {
    return `
      <div class="state-block">
        <span class="state-icon">${Icons.image()}</span>
        <div class="state-title">No local images</div>
        <div class="state-desc">Pull an image from Docker Hub.</div>
        <button class="btn btn-primary" id="empty-pull">${Icons.download()}Pull image</button>
      </div>`;
  },

  _wireEmpty(mount, body) {
    const btn = body.querySelector('#empty-pull');
    if (btn) btn.addEventListener('click', () => openPullImageModal({ onPulled: () => this._load(mount) }));
  },
};
