/**
 * icons.js
 * A small local icon utility. Every icon is a plain function returning an
 * SVG string, stroke-based, 1.75px, currentColor — no icon framework.
 */

const ICON_DEFAULTS = 'width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"';

function icon(inner, attrs = '') {
  return `<svg ${ICON_DEFAULTS} ${attrs} aria-hidden="true">${inner}</svg>`;
}

const Icons = {
  logo: (size = 26) => `
    <svg width="${size}" height="${size}" viewBox="0 0 32 32" fill="none" aria-hidden="true">
      <rect x="4" y="6" width="24" height="4.5" rx="1.5" fill="var(--primary)"/>
      <rect x="4" y="13.75" width="24" height="4.5" rx="1.5" fill="var(--text)"/>
      <rect x="4" y="21.5" width="24" height="4.5" rx="1.5" fill="var(--border-strong)"/>
    </svg>`,

  dashboard: () => icon('<rect x="3" y="3" width="7" height="9" rx="1.5"/><rect x="14" y="3" width="7" height="5" rx="1.5"/><rect x="14" y="12" width="7" height="9" rx="1.5"/><rect x="3" y="16" width="7" height="5" rx="1.5"/>'),
  container: () => icon('<path d="M3 8l9-5 9 5-9 5-9-5z"/><path d="M3 8v8l9 5 9-5V8"/><path d="M12 13v8"/>'),
  image: () => icon('<rect x="3" y="4" width="18" height="16" rx="2"/><circle cx="8.5" cy="9.5" r="1.75"/><path d="M21 16l-5.5-5.5L4 21"/>'),
  network: () => icon('<circle cx="12" cy="5" r="2.5"/><circle cx="5" cy="19" r="2.5"/><circle cx="19" cy="19" r="2.5"/><path d="M12 7.5v4M12 11.5L6.5 17M12 11.5L17.5 17"/>'),
  terminal: () => icon('<rect x="3" y="4" width="18" height="16" rx="2"/><path d="M7 9l3 3-3 3"/><path d="M12 15h5"/>'),
  settings: () => icon('<circle cx="12" cy="12" r="3"/><path d="M19.4 13a1.7 1.7 0 00.34 1.87l.06.06a2 2 0 11-2.83 2.83l-.06-.06a1.7 1.7 0 00-1.87-.34 1.7 1.7 0 00-1.04 1.56V19a2 2 0 11-4 0v-.09A1.7 1.7 0 008.44 17a1.7 1.7 0 00-1.87.34l-.06.06a2 2 0 11-2.83-2.83l.06-.06A1.7 1.7 0 004.08 13a1.7 1.7 0 00-1.56-1.04H2.5a2 2 0 110-4h.09A1.7 1.7 0 004.08 6.9a1.7 1.7 0 00-.34-1.87l-.06-.06a2 2 0 112.83-2.83l.06.06A1.7 1.7 0 008.44 2.5a1.7 1.7 0 001.04-1.56V.94a2 2 0 114 0V1a1.7 1.7 0 001.04 1.56 1.7 1.7 0 001.87-.34l.06-.06a2 2 0 112.83 2.83l-.06.06A1.7 1.7 0 0018.9 8.44c.13.62.62 1.13 1.24 1.24H20a2 2 0 010 4h-.09a1.7 1.7 0 00-1.51 1.32z"/>'),
  play: () => icon('<path d="M6 4l14 8-14 8V4z"/>'),
  stop: () => icon('<rect x="6" y="6" width="12" height="12" rx="1.5"/>'),
  kill: () => icon('<path d="M18 6L6 18M6 6l12 12"/>'),
  refresh: () => icon('<path d="M21 4v6h-6"/><path d="M3 20v-6h6"/><path d="M4.6 9A9 9 0 0119.4 6.4M19.4 15a9 9 0 01-14.8 2.6"/>'),
  search: () => icon('<circle cx="11" cy="11" r="7"/><path d="M21 21l-4.3-4.3"/>'),
  plus: () => icon('<path d="M12 5v14M5 12h14"/>'),
  external: () => icon('<path d="M14 4h6v6"/><path d="M20 4L10 14"/><path d="M18 13v6a1 1 0 01-1 1H5a1 1 0 01-1-1V6a1 1 0 011-1h6"/>'),
  copy: () => icon('<rect x="9" y="9" width="12" height="12" rx="2"/><path d="M5 15H4a1 1 0 01-1-1V4a1 1 0 011-1h10a1 1 0 011 1v1"/>'),
  warning: () => icon('<path d="M10.3 3.9L1.8 18a2 2 0 001.7 3h17a2 2 0 001.7-3L13.7 3.9a2 2 0 00-3.4 0z"/><path d="M12 9v4"/><path d="M12 17h.01"/>'),
  check: () => icon('<path d="M20 6L9 17l-5-5"/>'),
  x: () => icon('<path d="M18 6L6 18M6 6l12 12"/>'),
  menu: () => icon('<path d="M3 6h18M3 12h18M3 18h18"/>'),
  code: () => icon('<path d="M8 5l-6 7 6 7"/><path d="M16 5l6 7-6 7"/>'),
  cpu: () => icon('<rect x="7" y="7" width="10" height="10" rx="1.5"/><rect x="2" y="10" width="3" height="4"/><rect x="19" y="10" width="3" height="4"/><rect x="10" y="2" width="4" height="3"/><rect x="10" y="19" width="4" height="3"/>'),
  memory: () => icon('<rect x="3" y="6" width="18" height="12" rx="1.5"/><path d="M7 6v12M11 6v12M15 6v12"/>'),
  pulse: () => icon('<path d="M22 12h-4l-3 8-6-16-3 8H2"/>'),
  chevronRight: () => icon('<path d="M9 6l6 6-6 6"/>'),
  info: () => icon('<circle cx="12" cy="12" r="9"/><path d="M12 16v-5"/><path d="M12 8h.01"/>'),
  clock: () => icon('<circle cx="12" cy="12" r="9"/><path d="M12 7v5l3.5 2"/>'),
  bridge: () => icon('<path d="M3 20V10a2 2 0 012-2h14a2 2 0 012 2v10"/><path d="M3 20h18"/><path d="M7 20v-6M12 20v-6M17 20v-6"/>'),
  download: () => icon('<path d="M12 3v12"/><path d="M7 10l5 5 5-5"/><path d="M5 21h14"/>'),
  boxOpen: () => icon('<path d="M3 8l9-5 9 5"/><path d="M3 8l9 5 9-5"/><path d="M12 13v8"/><path d="M3 8v8l4 2.2M21 8v8l-4 2.2"/>'),
  server: () => icon('<rect x="3" y="4" width="18" height="7" rx="1.5"/><rect x="3" y="13" width="18" height="7" rx="1.5"/><path d="M7 7.5h.01M7 16.5h.01"/>'),
};

window.Icons = Icons;
