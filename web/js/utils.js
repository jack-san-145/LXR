/**
 * utils.js — small reusable helpers. No dependencies.
 */

const Utils = {
  /** Create a DOM element from an HTML string (single root element). */
  h(html) {
    const tpl = document.createElement('template');
    tpl.innerHTML = html.trim();
    return tpl.content.firstElementChild;
  },

  qs(sel, root = document) { return root.querySelector(sel); },
  qsa(sel, root = document) { return Array.from(root.querySelectorAll(sel)); },

  escapeHtml(str) {
    if (str === null || str === undefined) return '';
    return String(str)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
  },

  /** Debounce helper for search inputs etc. */
  debounce(fn, wait = 250) {
    let t;
    return (...args) => {
      clearTimeout(t);
      t = setTimeout(() => fn(...args), wait);
    };
  },

  timeAgo(date) {
    const d = date instanceof Date ? date : new Date(date);
    const diff = Math.floor((Date.now() - d.getTime()) / 1000);
    if (diff < 5) return 'just now';
    if (diff < 60) return `${diff}s ago`;
    if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
    if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
    return `${Math.floor(diff / 86400)}d ago`;
  },

  shortId(id, len = 12) {
    if (!id) return '—';
    return id.length > len ? id.slice(0, len) : id;
  },

  /** Strip CIDR suffix like /17 for display, keep full value available on hover. */
  ipOnly(ipCidr) {
    if (!ipCidr) return null;
    return String(ipCidr).split('/')[0];
  },

  copyToClipboard(text) {
    if (navigator.clipboard && window.isSecureContext) {
      return navigator.clipboard.writeText(text);
    }
    // Fallback for non-secure contexts
    const ta = document.createElement('textarea');
    ta.value = text;
    ta.style.position = 'fixed';
    ta.style.opacity = '0';
    document.body.appendChild(ta);
    ta.select();
    try { document.execCommand('copy'); } finally { document.body.removeChild(ta); }
    return Promise.resolve();
  },

  clamp(n, min, max) { return Math.max(min, Math.min(max, n)); },

  /** Human label for a container's simplified status. */
  statusLabel(status) {
    switch ((status || '').toLowerCase()) {
      case 'running': return 'Running';
      case 'frozen': return 'Frozen';
      case 'stopped': return 'Stopped';
      case 'inactive': return 'Stopped';
      default: return status ? status : 'Unknown';
    }
  },

  statusBadgeClass(status) {
    switch ((status || '').toLowerCase()) {
      case 'running': return 'badge-running';
      case 'frozen': return 'badge-frozen';
      case 'stopped':
      case 'inactive': return 'badge-stopped';
      default: return 'badge-neutral';
    }
  },
};

window.Utils = Utils;
