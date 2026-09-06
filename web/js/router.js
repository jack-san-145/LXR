/**
 * router.js — a small pushState router. No framework.
 *
 * Routes are registered as { pattern: '/containers/:name', render(params) }.
 * Falls back gracefully: if the app isn't served with a catch-all rewrite
 * for deep links, index.html can be configured server-side to always serve
 * itself; this router only needs the initial load to reach index.html.
 */

const Router = (() => {
  const routes = [];
  let mainEl = null;
  let notFoundRender = () => '<div class="state-block"><div class="state-title">Page not found</div></div>';

  function register(pattern, render) {
    const paramNames = [];
    const regexStr = pattern.replace(/:[^/]+/g, (m) => {
      paramNames.push(m.slice(1));
      return '([^/]+)';
    });
    const regex = new RegExp(`^${regexStr}/?$`);
    routes.push({ pattern, regex, paramNames, render });
  }

  function match(path) {
    for (const r of routes) {
      const m = path.match(r.regex);
      if (m) {
        const params = {};
        r.paramNames.forEach((name, i) => { params[name] = decodeURIComponent(m[i + 1]); });
        return { route: r, params };
      }
    }
    return null;
  }

  async function resolve(path) {
    const found = match(path);
    if (!mainEl) return;
    if (!found) {
      mainEl.innerHTML = notFoundRender();
      return;
    }
    try {
      await found.route.render(found.params, mainEl);
    } catch (err) {
      console.error(err);
      mainEl.innerHTML = `
        <div class="state-block">
          <span class="state-icon">${Icons.warning()}</span>
          <div class="state-title">Something went wrong rendering this page</div>
          <div class="state-desc">${Utils.escapeHtml(err.message || String(err))}</div>
        </div>`;
    }
    document.dispatchEvent(new CustomEvent('lxr:navigated', { detail: { path } }));
  }

  function navigate(path, { replace = false } = {}) {
    if (replace) history.replaceState({}, '', path);
    else history.pushState({}, '', path);
    resolve(path);
    window.scrollTo(0, 0);
  }

  function init(mount) {
    mainEl = mount;
    document.body.addEventListener('click', (e) => {
      const a = e.target.closest('a[data-link]');
      if (!a) return;
      const href = a.getAttribute('href');
      if (!href || href.startsWith('http') || a.target === '_blank') return;
      e.preventDefault();
      navigate(href);
    });
    window.addEventListener('popstate', () => resolve(location.pathname));
    resolve(location.pathname);
  }

  return { register, navigate, init, setNotFound: (fn) => { notFoundRender = fn; } };
})();

window.Router = Router;
