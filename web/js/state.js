/**
 * state.js — a minimal centralized store. No external library.
 * Components subscribe to keys they care about and re-render on change.
 */

const state = {
  runtimeOnline: null,       // null = unknown/checking, true/false after first ping
  containers: [],
  containersLoaded: false,
  images: [],
  network: null,
  loading: {},               // { [key]: boolean }
  errors: {},                // { [key]: string }
  activity: [],              // local UI activity log, newest first
};

const listeners = new Map(); // key -> Set<fn>

function subscribe(key, fn) {
  if (!listeners.has(key)) listeners.set(key, new Set());
  listeners.get(key).add(fn);
  return () => listeners.get(key).delete(fn);
}

function emit(key) {
  const set = listeners.get(key);
  if (set) set.forEach(fn => fn(state));
  const all = listeners.get('*');
  if (all) all.forEach(fn => fn(state));
}

function setState(patch) {
  Object.assign(state, patch);
  Object.keys(patch).forEach(emit);
}

function setLoading(key, value) {
  state.loading = { ...state.loading, [key]: value };
  emit('loading');
}

function setError(key, message) {
  state.errors = { ...state.errors, [key]: message };
  emit('errors');
}

function clearError(key) {
  if (!(key in state.errors)) return;
  const next = { ...state.errors };
  delete next[key];
  state.errors = next;
  emit('errors');
}

function pushActivity(entry) {
  state.activity = [{ id: Date.now() + Math.random(), time: new Date(), ...entry }, ...state.activity].slice(0, 30);
  emit('activity');
}

window.AppState = { state, subscribe, setState, setLoading, setError, clearError, pushActivity };
