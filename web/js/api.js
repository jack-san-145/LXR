/**
 * api.js
 *
 * Central API client for the LXR web console.
 *
 * IMPORTANT: this file only claims to talk to endpoints that actually exist
 * on the current LXR daemon (see section 14-24 of the product spec):
 *
 *   GET    /ping
 *   GET    /ps
 *   GET    /ps/all
 *   POST   /create        (streams plain text)
 *   POST   /start         (streams plain text)
 *   GET    /stop
 *   DELETE /kill
 *   GET    /exec           (hijacked PTY connection — not fetch-able)
 *   POST   /pull_image     (streams plain text)
 *
 * Everything else the product spec describes (images list, network info,
 * per-container resources/processes, events, websocket terminal, etc.) is
 * NOT implemented by the daemon yet. Those calls are represented below as
 * clearly-labeled "unsupported" adapters so the rest of the app can call
 * them today and simply switch to real data the day the backend adds them,
 * with zero call-site changes.
 */

const API_BASE = window.LXR_API_BASE || '/api';

/** Demo/mock mode — OFF by default. See js/demo.js for fixtures. */
const DEMO_MODE = window.LXR_DEMO_MODE === true;

class ApiError extends Error {
  constructor(message, { status = null, cause = null } = {}) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.cause = cause;
  }
}

/** Marks a response that is architecturally reserved for a future API. */
class UnsupportedError extends Error {
  constructor(featureName) {
    super(`${featureName} is not exposed by the current LXR API yet.`);
    this.name = 'UnsupportedError';
    this.feature = featureName;
  }
}

async function request(path, options = {}) {
  let res;
  try {
    res = await fetch(`${API_BASE}${path}`, {
      headers: { 'Content-Type': 'application/json', ...(options.headers || {}) },
      ...options,
    });
  } catch (err) {
    throw new ApiError('Could not reach the LXR web server.', { cause: err });
  }

  if (!res.ok) {
    let detail = '';
    try { detail = await res.text(); } catch (_) { /* ignore */ }
    throw new ApiError(detail || `Request failed (${res.status})`, { status: res.status });
  }
  return res;
}

async function requestJson(path, options = {}) {
  const res = await request(path, options);
  const text = await res.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch (err) {
    throw new ApiError('The server returned a response that was not valid JSON.', { cause: err });
  }
}

const LXRApi = {
  ApiError,
  UnsupportedError,
  DEMO_MODE,

  // ---------------------------------------------------------------- health
  async ping() {
    if (DEMO_MODE) return true;
    const res = await request('/ping');
    const text = (await res.text()).trim();
    return /pong/i.test(text) || res.ok;
  },

  // ------------------------------------------------------------ containers
  async getRunningContainers() {
    if (DEMO_MODE) return window.LXRDemo.containers.filter(c => c.status === 'running');
    const data = await requestJson('/ps');
    return (data && data.containers) || [];
  },

  async getContainers() {
    if (DEMO_MODE) return window.LXRDemo.containers;
    const data = await requestJson('/ps/all');
    return (data && data.containers) || [];
  },

  async getContainer(name) {
    // Not exposed as a single-resource endpoint yet; derive from /ps/all.
    const all = await this.getContainers();
    return all.find(c => c.container_name === name) || null;
  },

  async createContainer({ container_name, image_name }, streamHandlers) {
    if (DEMO_MODE) return window.LXRDemo.simulateCreate(container_name, image_name, streamHandlers);
    return streamRequest(`${API_BASE}/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ container_name, image_name }),
    }, streamHandlers);
  },

  async startContainer(name, streamHandlers) {
    if (DEMO_MODE) return window.LXRDemo.simulateStart(name, streamHandlers);
    return streamRequest(`${API_BASE}/start`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ container_name: name }),
    }, streamHandlers);
  },

  async stopContainer(name) {
    if (DEMO_MODE) return window.LXRDemo.simulateStop(name);
    return requestJson(`/stop?container_name=${encodeURIComponent(name)}`);
  },

  async killContainer(name) {
    if (DEMO_MODE) return window.LXRDemo.simulateKill(name);
    return requestJson(`/kill?container_name=${encodeURIComponent(name)}`, { method: 'DELETE' });
  },

  /**
   * Interactive terminal. The current daemon hijacks a plain HTTP connection
   * for /exec and attaches a PTY over nsenter — that transport cannot be
   * driven from a browser fetch(). The future web server is expected to
   * expose this over WebSocket at /api/containers/:name/terminal. This
   * function is the seam: swap the implementation here once that lands.
   */
  connectTerminal(name) {
    throw new UnsupportedError('Interactive terminal transport');
  },

  // ----------------------------------------------------------------images
  async pullImage(imageName, streamHandlers) {
    if (DEMO_MODE) return window.LXRDemo.simulatePull(imageName, streamHandlers);
    return streamRequest(`${API_BASE}/pull_image`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ img_name: imageName }),
    }, streamHandlers);
  },

  async getImages() {
    if (DEMO_MODE) return window.LXRDemo.images;
    throw new UnsupportedError('Image listing');
  },

  // --------------------------------------------------------------- future
  async getRuntime() {
    if (DEMO_MODE) return window.LXRDemo.runtime;
    throw new UnsupportedError('Runtime configuration endpoint');
  },
  async getContainerProcesses(_name) { throw new UnsupportedError('Process listing'); },
  async getContainerResources(_name) { throw new UnsupportedError('Resource metrics'); },
  async getContainerNetwork(_name) { throw new UnsupportedError('Per-container network detail'); },
  async getNetwork() {
    if (DEMO_MODE) return window.LXRDemo.network;
    throw new UnsupportedError('Network overview endpoint');
  },
  async getEvents() { throw new UnsupportedError('Runtime event stream'); },
};

window.LXRApi = LXRApi;
window.ApiError = ApiError;
