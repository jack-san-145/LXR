/**
 * notifications.js — reusable toast notifications.
 * Usage: Notify.success('Container started'); Notify.error('Failed to start container');
 */

const Notify = (() => {
  let stack;

  function ensureStack() {
    if (!stack) {
      stack = Utils.h('<div class="toast-stack" role="status" aria-live="polite"></div>');
      document.body.appendChild(stack);
    }
    return stack;
  }

  const ICONS = {
    success: Icons.check(),
    error: Icons.warning(),
    warning: Icons.warning(),
    info: Icons.info(),
  };

  function show(message, { type = 'info', duration = 4200, detail = null } = {}) {
    const root = ensureStack();
    const el = Utils.h(`
      <div class="toast toast-${type}" role="alert">
        <span class="toast-icon">${ICONS[type] || ICONS.info}</span>
        <div class="toast-body">
          <div>${Utils.escapeHtml(message)}</div>
          ${detail ? `<div class="help-text" style="margin-top:4px;">${Utils.escapeHtml(detail)}</div>` : ''}
        </div>
        <button class="toast-close" aria-label="Dismiss">${Icons.x()}</button>
      </div>
    `);
    root.appendChild(el);

    const dismiss = () => {
      el.classList.add('leaving');
      setTimeout(() => el.remove(), 200);
    };
    el.querySelector('.toast-close').addEventListener('click', dismiss);
    if (duration) setTimeout(dismiss, duration);
    return dismiss;
  }

  return {
    show,
    success: (msg, opts) => show(msg, { ...opts, type: 'success' }),
    error: (msg, opts) => show(msg, { ...opts, type: 'error', duration: 6000 }),
    warning: (msg, opts) => show(msg, { ...opts, type: 'warning' }),
    info: (msg, opts) => show(msg, { ...opts, type: 'info' }),
  };
})();

window.Notify = Notify;
