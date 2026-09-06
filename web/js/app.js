/**
 * app.js — boots the shell (sidebar/header), wires global behaviors, and
 * registers routes. Individual pages live in js/pages/*.js.
 */

const NAV_ITEMS = [
  { path: '/', label: 'Dashboard', icon: 'dashboard' },
  { path: '/containers', label: 'Containers', icon: 'container' },
  { path: '/images', label: 'Images', icon: 'image' },
  { path: '/network', label: 'Network', icon: 'network' },
  { path: '/settings', label: 'Settings', icon: 'settings' },
];

function renderShell() {
  const shell = Utils.h(`
    <div class="app-shell">
      <aside class="sidebar" id="sidebar">
        <div class="brand">
          <span class="brand-mark">${Icons.logo()}</span>
          <div class="brand-text">
            <span class="brand-name">LXR</span>
            <span class="brand-sub">Linux Container Runtime</span>
          </div>
        </div>
        <nav class="nav" id="nav-list" aria-label="Primary"></nav>
        <div class="sidebar-footer">
          <span class="runtime-pill">
            <span class="runtime-dot pending" id="runtime-dot"></span>
            <span id="runtime-label">Checking runtime…</span>
          </span>
        </div>
      </aside>
      <div class="sidebar-scrim" id="sidebar-scrim"></div>
      <header class="topbar">
        <div class="topbar-left">
          <button class="menu-toggle" id="menu-toggle" aria-label="Toggle navigation">${Icons.menu()}</button>
          <div class="breadcrumbs" id="breadcrumbs"></div>
        </div>
        <div class="topbar-right">
          <span class="runtime-pill" id="topbar-status" title="LXR daemon health">
            <span class="runtime-dot pending" id="runtime-dot-top"></span>
            <span id="runtime-label-top">Checking…</span>
          </span>
        </div>
      </header>
      <main class="main" id="main"></main>
    </div>
  `);
  document.body.appendChild(shell);

  const navList = shell.querySelector('#nav-list');
  NAV_ITEMS.forEach(item => {
    const el = Utils.h(`
      <a class="nav-item" href="${item.path}" data-link data-path="${item.path}">
        ${Icons[item.icon]()}
        <span>${item.label}</span>
      </a>
    `);
    navList.appendChild(el);
  });

  // mobile drawer
  const sidebar = shell.querySelector('#sidebar');
  const scrim = shell.querySelector('#sidebar-scrim');
  shell.querySelector('#menu-toggle').addEventListener('click', () => {
    sidebar.classList.toggle('open');
    scrim.classList.toggle('open');
  });
  scrim.addEventListener('click', () => {
    sidebar.classList.remove('open');
    scrim.classList.remove('open');
  });
  document.body.addEventListener('click', (e) => {
    if (e.target.closest('a[data-link]')) {
      sidebar.classList.remove('open');
      scrim.classList.remove('open');
    }
  });

  return shell;
}

function updateActiveNav() {
  const path = location.pathname;
  Utils.qsa('.nav-item').forEach(a => {
    const p = a.dataset.path;
    const active = p === '/' ? path === '/' : path.startsWith(p);
    a.classList.toggle('active', active);
  });
}

function setBreadcrumbs(parts) {
  const el = Utils.qs('#breadcrumbs');
  if (!el) return;
  el.innerHTML = parts.map((p, i) => {
    const isLast = i === parts.length - 1;
    if (isLast) return `<span class="crumb-current">${Utils.escapeHtml(p.label)}</span>`;
    const sep = i > 0 ? `<span class="sep">${Icons.chevronRight()}</span>` : '';
    return `${sep}<a href="${p.path}" data-link>${Utils.escapeHtml(p.label)}</a>`;
  }).join('');
}

function updateRuntimeIndicator(online) {
  const dotIds = ['#runtime-dot', '#runtime-dot-top'];
  const labelIds = ['#runtime-label', '#runtime-label-top'];
  dotIds.forEach(sel => {
    const dot = Utils.qs(sel);
    if (dot) dot.className = `runtime-dot ${online === null ? 'pending' : online ? 'online' : 'offline'}`;
  });
  labelIds.forEach(sel => {
    const label = Utils.qs(sel);
    if (label) label.textContent = online === null ? 'Checking runtime…' : online ? 'LXR Online' : 'LXR Offline';
  });
}

async function pollRuntimeHealth() {
  try {
    const online = await LXRApi.ping();
    AppState.setState({ runtimeOnline: online });
  } catch (err) {
    AppState.setState({ runtimeOnline: false });
  }
}

function startHealthPolling() {
  AppState.subscribe('runtimeOnline', (s) => updateRuntimeIndicator(s.runtimeOnline));
  pollRuntimeHealth();
  let interval = setInterval(() => { if (!document.hidden) pollRuntimeHealth(); }, 8000);
  document.addEventListener('visibilitychange', () => {
    if (!document.hidden) pollRuntimeHealth();
  });
}

/** Periodic container-state refresh, paused while the tab is hidden. */
function startContainerPolling() {
  async function tick() {
    if (document.hidden) return;
    try {
      const containers = await LXRApi.getContainers();
      AppState.setState({ containers, containersLoaded: true });
      AppState.clearError('containers');
    } catch (err) {
      AppState.setError('containers', err.message);
    }
  }
  tick();
  setInterval(tick, 7000);
}

function registerRoutes() {
  Router.register('/', async (params, mount) => {
    setBreadcrumbs([{ label: 'Dashboard', path: '/' }]);
    await Pages.dashboard.render(mount);
  });
  Router.register('/containers', async (params, mount) => {
    setBreadcrumbs([{ label: 'Containers', path: '/containers' }]);
    await Pages.containers.render(mount);
  });
  Router.register('/containers/:name/terminal', async (params, mount) => {
    setBreadcrumbs([
      { label: 'Containers', path: '/containers' },
      { label: params.name, path: `/containers/${params.name}` },
      { label: 'Terminal', path: `/containers/${params.name}/terminal` },
    ]);
    await Pages.terminal.render(mount, params.name);
  });
  Router.register('/containers/:name', async (params, mount) => {
    setBreadcrumbs([
      { label: 'Containers', path: '/containers' },
      { label: params.name, path: `/containers/${params.name}` },
    ]);
    await Pages.container.render(mount, params.name);
  });
  Router.register('/images', async (params, mount) => {
    setBreadcrumbs([{ label: 'Images', path: '/images' }]);
    await Pages.images.render(mount);
  });
  Router.register('/network', async (params, mount) => {
    setBreadcrumbs([{ label: 'Network', path: '/network' }]);
    await Pages.network.render(mount);
  });
  Router.register('/settings', async (params, mount) => {
    setBreadcrumbs([{ label: 'Settings', path: '/settings' }]);
    await Pages.settings.render(mount);
  });

  Router.setNotFound(() => `
    <div class="state-block">
      <span class="state-icon">${Icons.warning()}</span>
      <div class="state-title">Page not found</div>
      <div class="state-desc">That view doesn't exist in the LXR console.</div>
      <a href="/" data-link class="btn btn-primary" style="margin-top:8px;">Back to dashboard</a>
    </div>
  `);
}

document.addEventListener('DOMContentLoaded', () => {
  renderShell();
  registerRoutes();
  Router.init(Utils.qs('#main'));
  document.addEventListener('lxr:navigated', updateActiveNav);
  updateActiveNav();
  startHealthPolling();
  startContainerPolling();

  if (LXRApi.DEMO_MODE) {
    const banner = Utils.h(`
      <div style="background:var(--warning-wash);color:var(--warning);border-bottom:1px solid var(--border-strong);
        font-size:12px;text-align:center;padding:6px;font-weight:600;letter-spacing:.02em;">
        Demo mode — showing mock data, not a live LXR daemon
      </div>
    `);
    document.body.insertBefore(banner, document.body.firstChild);
  }
});
