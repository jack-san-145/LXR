/**
 * components.js — shared rendering + action helpers used by multiple pages.
 * Keeps container action semantics (start/stop/kill/terminal/open code) and
 * their modals in one place instead of duplicated across pages.
 */

const LXRActions = {
  /** Which row actions are valid for a given status. */
  actionsFor(status) {
    const s = (status || '').toLowerCase();
    if (s === 'running') return ['terminal', 'open', 'stop', 'kill'];
    if (s === 'frozen') return ['start', 'kill'];
    return ['start', 'kill']; // stopped / inactive / unknown
  },

  codeServerUrl(container) {
    const port = (container.ports && container.ports[0]) || 9000;
    return `http://${window.location.hostname}:${port}/`;
  },

  async doStart(name, { onDone } = {}) {
    const modal = Modal.open({
      title: `Starting ${name}`,
      dismissible: false,
      bodyHtml: `<div class="stream-log cursor-blink" id="stream-log"></div>`,
    });
    const logEl = () => Utils.qs('#stream-log', modal.el);
    const append = (line, cls = '') => {
      const el = Utils.h(`<div class="line ${cls}">${Utils.escapeHtml(line) || '&nbsp;'}</div>`);
      logEl().appendChild(el);
      logEl().scrollTop = logEl().scrollHeight;
    };

    try {
      await LXRApi.startContainer(name, {
        onLine: (line) => append(line),
        onComplete: async () => {
          modal.setDismissible(true);
          modal.setFooter(`<button class="btn btn-primary" data-act="close">Done</button>`);
          modal.el.querySelector('[data-act="close"]').addEventListener('click', modal.close);
          Notify.success(`${name} started`);
          AppState.pushActivity({ type: 'success', text: `Started container “${name}”` });
          onDone && onDone();
        },
        onError: (err) => {
          append(err.message, 'err');
          modal.setDismissible(true);
          modal.setFooter(`<button class="btn btn-secondary" data-act="close">Close</button>`);
          modal.el.querySelector('[data-act="close"]').addEventListener('click', modal.close);
          Notify.error(`Failed to start ${name}`, { detail: err.message });
        },
      });
    } catch (err) {
      modal.close();
      Notify.error(`Failed to start ${name}`, { detail: err.message });
    }
  },

  async doStop(name, { onDone } = {}) {
    AppState.setLoading(`stop:${name}`, true);
    try {
      const res = await LXRApi.stopContainer(name);
      if (!res.container_exists) {
        Notify.error(`${name} no longer exists`);
      } else if (res.container_Stoped) {
        Notify.success(`${name} stopped`, { detail: 'The container was frozen via cgroup.freeze.' });
        AppState.pushActivity({ type: 'success', text: `Stopped (froze) container “${name}”` });
      } else {
        Notify.warning(`${name} was not running`);
      }
      onDone && onDone();
    } catch (err) {
      Notify.error(`Failed to stop ${name}`, { detail: err.message });
    } finally {
      AppState.setLoading(`stop:${name}`, false);
    }
  },

  async doKill(name, { onDone } = {}) {
    const confirmed = await Modal.confirmDanger({
      title: 'Kill container?',
      message: `This permanently removes the container runtime state, filesystem, networking resources, cgroup, and metadata for <strong>${Utils.escapeHtml(name)}</strong>. This cannot be undone.`,
      confirmLabel: 'Kill container',
    });
    if (!confirmed) return;

    AppState.setLoading(`kill:${name}`, true);
    try {
      const res = await LXRApi.killContainer(name);
      if (!res.container_exists) {
        Notify.error(`${name} no longer exists`);
      } else if (res.is_killed) {
        Notify.success(`${name} killed`);
        AppState.pushActivity({ type: 'success', text: `Killed container “${name}”` });
      } else {
        Notify.error(`Failed to kill ${name}`);
      }
      onDone && onDone();
    } catch (err) {
      Notify.error(`Failed to kill ${name}`, { detail: err.message });
    } finally {
      AppState.setLoading(`kill:${name}`, false);
    }
  },

  openCode(container) {
    window.open(this.codeServerUrl(container), '_blank', 'noopener');
  },

  goToTerminal(name) {
    Router.navigate(`/containers/${encodeURIComponent(name)}/terminal`);
  },
};

/** Renders a "Create container" modal wired to the real streaming endpoint. */
function openCreateContainerModal({ onCreated } = {}) {
  const modal = Modal.open({
    title: 'Create container',
    bodyHtml: `
      <form id="create-form">
        <div class="field">
          <label class="field-label" for="cc-name">Container name</label>
          <input class="input" id="cc-name" placeholder="python-dev" autocomplete="off" required />
        </div>
        <div class="field">
          <label class="field-label" for="cc-image">Image</label>
          <input class="input" id="cc-image" placeholder="python" autocomplete="off" required />
          <span class="field-hint">Pulled from Docker Hub if not already cached locally.</span>
        </div>
        <fieldset class="disabled-future">
          <div class="field-hint" style="margin-bottom:6px;">
            <span class="tag-soon">Coming soon</span> Resource limits, port mapping, and environment variables
          </div>
        </fieldset>
      </form>
    `,
    footerHtml: `
      <button class="btn btn-secondary" data-act="cancel">Cancel</button>
      <button class="btn btn-primary" data-act="submit">${Icons.plus()}Create container</button>
    `,
  });

  modal.el.querySelector('[data-act="cancel"]').addEventListener('click', modal.close);
  modal.el.querySelector('[data-act="submit"]').addEventListener('click', () => {
    const name = Utils.qs('#cc-name', modal.el).value.trim();
    const image = Utils.qs('#cc-image', modal.el).value.trim();
    if (!name || !image) {
      Notify.warning('Container name and image are both required');
      return;
    }
    runCreateStream(modal, name, image, onCreated);
  });
}

const CREATE_STEP_MATCHERS = [
  { key: 'image', label: 'Image found', test: (l) => /find image/i.test(l) },
  { key: 'rootfs', label: 'Root filesystem', test: (l) => /rootfs/i.test(l) },
  { key: 'env', label: 'Container environment', test: (l) => /building container environment/i.test(l) },
  { key: 'net', label: 'Networking', test: (l) => /networking/i.test(l) },
  { key: 'limits', label: 'Resource limits', test: (l) => /resource(s)? limit/i.test(l) },
  { key: 'code', label: 'Code-server', test: (l) => /code-server/i.test(l) },
];

function runCreateStream(modal, name, image, onCreated) {
  modal.setDismissible(false);
  const steps = CREATE_STEP_MATCHERS.map(s => ({ ...s, done: false }));
  const renderSteps = () => steps.map(s => `
    <div class="step-item ${s.done ? 'done' : ''}">
      <span class="step-dot">${s.done ? Icons.check() : ''}</span>
      <span>${s.label}</span>
    </div>
  `).join('');

  modal.setBody(`
    <div class="stream-log cursor-blink" id="stream-log"></div>
    <div class="step-list" id="step-list">${renderSteps()}</div>
  `);
  modal.setFooter(`<button class="btn btn-secondary" data-act="close" disabled>Creating…</button>`);

  const logEl = () => Utils.qs('#stream-log', modal.el);
  const append = (line, cls = '') => {
    const el = Utils.h(`<div class="line ${cls}">${Utils.escapeHtml(line) || '&nbsp;'}</div>`);
    logEl().appendChild(el);
    logEl().scrollTop = logEl().scrollHeight;
  };

  LXRApi.createContainer({ container_name: name, image_name: image }, {
    onLine: (line) => {
      append(line);
      let changed = false;
      steps.forEach(s => { if (!s.done && s.test(line)) { s.done = true; changed = true; } });
      if (changed) Utils.qs('#step-list', modal.el).innerHTML = renderSteps();
    },
    onComplete: () => {
      steps.forEach(s => s.done = true);
      Utils.qs('#step-list', modal.el).innerHTML = renderSteps();
      modal.setDismissible(true);
      modal.setFooter(`<button class="btn btn-primary" data-act="close">Done</button>`);
      modal.el.querySelector('[data-act="close"]').addEventListener('click', modal.close);
      Notify.success(`Container “${name}” created`);
      AppState.pushActivity({ type: 'success', text: `Created container “${name}” from image “${image}”` });
      onCreated && onCreated();
    },
    onError: (err) => {
      append(err.message, 'err');
      modal.setDismissible(true);
      modal.setFooter(`<button class="btn btn-secondary" data-act="close">Close</button>`);
      modal.el.querySelector('[data-act="close"]').addEventListener('click', modal.close);
      Notify.error(`Failed to create ${name}`, { detail: err.message });
    },
  });
}

/** Pull-image modal wired to the real streaming endpoint. */
function openPullImageModal({ onPulled } = {}) {
  const modal = Modal.open({
    title: 'Pull image',
    bodyHtml: `
      <div class="field">
        <label class="field-label" for="pi-name">Image name</label>
        <input class="input" id="pi-name" placeholder="ubuntu" autocomplete="off" required />
        <span class="field-hint">Pulled directly from Docker Hub into the LXR registry.</span>
      </div>
    `,
    footerHtml: `
      <button class="btn btn-secondary" data-act="cancel">Cancel</button>
      <button class="btn btn-primary" data-act="submit">${Icons.download()}Pull</button>
    `,
  });
  modal.el.querySelector('[data-act="cancel"]').addEventListener('click', modal.close);
  modal.el.querySelector('[data-act="submit"]').addEventListener('click', () => {
    const image = Utils.qs('#pi-name', modal.el).value.trim();
    if (!image) { Notify.warning('Enter an image name'); return; }
    runPullStream(modal, image, onPulled);
  });
}

function runPullStream(modal, image, onPulled) {
  modal.setDismissible(false);
  modal.setBody(`<div class="stream-log cursor-blink" id="stream-log"></div>`);
  modal.setFooter(`<button class="btn btn-secondary" disabled>Pulling…</button>`);
  const logEl = () => Utils.qs('#stream-log', modal.el);
  const append = (line, cls = '') => {
    const el = Utils.h(`<div class="line ${cls}">${Utils.escapeHtml(line) || '&nbsp;'}</div>`);
    logEl().appendChild(el);
    logEl().scrollTop = logEl().scrollHeight;
  };

  LXRApi.pullImage(image, {
    onLine: (line) => append(line),
    onComplete: () => {
      modal.setDismissible(true);
      modal.setFooter(`<button class="btn btn-primary" data-act="close">Done</button>`);
      modal.el.querySelector('[data-act="close"]').addEventListener('click', modal.close);
      Notify.success(`Image “${image}” pulled`);
      AppState.pushActivity({ type: 'success', text: `Pulled image “${image}”` });
      onPulled && onPulled();
    },
    onError: (err) => {
      append(err.message, 'err');
      modal.setDismissible(true);
      modal.setFooter(`<button class="btn btn-secondary" data-act="close">Close</button>`);
      modal.el.querySelector('[data-act="close"]').addEventListener('click', modal.close);
      Notify.error(`Failed to pull ${image}`, { detail: err.message });
    },
  });
}

/** Small copy-to-clipboard button generator for IDs/IPs/etc. */
function copyable(value, { mono = true } = {}) {
  const safe = Utils.escapeHtml(value ?? '—');
  if (value === null || value === undefined || value === '') return safe;
  return `
    <span class="copyable">
      <span class="${mono ? 'mono' : ''}">${safe}</span>
      <button class="copy-btn" data-copy="${Utils.escapeHtml(value)}" aria-label="Copy">${Icons.copy()}</button>
    </span>
  `;
}

document.addEventListener('click', (e) => {
  const btn = e.target.closest('.copy-btn');
  if (!btn) return;
  const val = btn.dataset.copy;
  Utils.copyToClipboard(val).then(() => Notify.success('Copied to clipboard'));
});

window.LXRActions = LXRActions;
window.openCreateContainerModal = openCreateContainerModal;
window.openPullImageModal = openPullImageModal;
window.copyable = copyable;
