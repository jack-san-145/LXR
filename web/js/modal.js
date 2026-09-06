/**
 * modal.js — reusable modal infrastructure. No alert()/confirm() anywhere.
 */

const Modal = (() => {
  let current = null;

  function close() {
    if (!current) return;
    const { backdrop, onClose } = current;
    backdrop.remove();
    document.removeEventListener('keydown', onEsc);
    current = null;
    onClose && onClose();
  }

  function onEsc(e) {
    if (e.key === 'Escape') close();
  }

  /**
   * open({ title, bodyHtml, footerHtml, size, onMount, dismissible })
   * Returns { el, close } so callers can wire up buttons and update content.
   */
  function open({ title, bodyHtml = '', footerHtml = '', size = '', dismissible = true, onClose } = {}) {
    if (current) close();

    const backdrop = Utils.h(`
      <div class="modal-backdrop">
        <div class="modal ${size === 'lg' ? 'modal-lg' : ''}" role="dialog" aria-modal="true" aria-label="${Utils.escapeHtml(title || 'Dialog')}">
          <div class="modal-header">
            <div class="modal-title">${Utils.escapeHtml(title || '')}</div>
            ${dismissible ? `<button class="modal-close" aria-label="Close">${Icons.x()}</button>` : ''}
          </div>
          <div class="modal-body">${bodyHtml}</div>
          ${footerHtml ? `<div class="modal-footer">${footerHtml}</div>` : ''}
        </div>
      </div>
    `);

    document.body.appendChild(backdrop);

    if (dismissible) {
      backdrop.addEventListener('mousedown', (e) => { if (e.target === backdrop) close(); });
      backdrop.querySelector('.modal-close').addEventListener('click', close);
      document.addEventListener('keydown', onEsc);
    }

    current = { backdrop, onClose };
    const modalEl = backdrop.querySelector('.modal');
    const focusable = modalEl.querySelector('input, button, select, textarea');
    if (focusable) focusable.focus();

    return {
      el: modalEl,
      backdrop,
      close,
      setBody(html) { modalEl.querySelector('.modal-body').innerHTML = html; },
      setFooter(html) {
        let footer = modalEl.querySelector('.modal-footer');
        if (!footer) {
          footer = Utils.h('<div class="modal-footer"></div>');
          modalEl.appendChild(footer);
        }
        footer.innerHTML = html;
      },
      setDismissible(val) {
        const closeBtn = modalEl.querySelector('.modal-close');
        if (closeBtn) closeBtn.style.display = val ? '' : 'none';
      },
    };
  }

  /** Danger confirmation modal (used for kill). Returns a Promise<boolean>. */
  function confirmDanger({ title, message, confirmLabel = 'Confirm', cancelLabel = 'Cancel' }) {
    return new Promise((resolve) => {
      let resolved = false;
      const m = open({
        title,
        dismissible: true,
        bodyHtml: `
          <div class="modal-danger-notice">
            ${Icons.warning()}
            <div>${message}</div>
          </div>
        `,
        footerHtml: `
          <button class="btn btn-secondary" data-act="cancel">${Utils.escapeHtml(cancelLabel)}</button>
          <button class="btn btn-danger" data-act="confirm">${Utils.escapeHtml(confirmLabel)}</button>
        `,
        onClose() { if (!resolved) resolve(false); },
      });
      m.el.querySelector('[data-act="cancel"]').addEventListener('click', () => { resolved = true; m.close(); resolve(false); });
      m.el.querySelector('[data-act="confirm"]').addEventListener('click', () => { resolved = true; m.close(); resolve(true); });
    });
  }

  /** Simple info/error-details modal. */
  function showError(title, message, technical = null) {
    open({
      title,
      bodyHtml: `
        <div class="modal-danger-notice">
          ${Icons.warning()}
          <div>${Utils.escapeHtml(message)}</div>
        </div>
        ${technical ? `
          <details>
            <summary class="help-text" style="cursor:pointer;">Technical details</summary>
            <div class="stream-log" style="margin-top:8px;">${Utils.escapeHtml(technical)}</div>
          </details>
        ` : ''}
      `,
      footerHtml: `<button class="btn btn-secondary" data-act="ok">Close</button>`,
    });
    Utils.qs('.modal-backdrop [data-act="ok"]').addEventListener('click', close);
  }

  return { open, close, confirmDanger, showError };
})();

window.Modal = Modal;
