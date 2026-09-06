/**
 * demo.js
 *
 * Mock data for previewing the console without a running LXR daemon.
 * This module is inert unless window.LXR_DEMO_MODE === true (see index.html).
 * Every value here is clearly demo data and is never shown to a user in a
 * way that could be mistaken for a real runtime response — pages that
 * render demo data render a persistent "Demo data" tag (see app.js).
 */

const LXRDemo = {
  containers: [
    { container_id: 'a1b2c3d4e5f6a1b2c3d4e5f6', container_name: 'python-dev', image: 'python', pid: 41210, status: 'running', ip_address: '10.10.0.2/17', ports: [9000] },
    { container_id: 'b2c3d4e5f6a1b2c3d4e5f6a1', container_name: 'ubuntu-dev', image: 'ubuntu', pid: 41288, status: 'stopped', ip_address: '10.10.0.3/17', ports: [9000] },
    { container_id: 'c3d4e5f6a1b2c3d4e5f6a1b2', container_name: 'node-api', image: 'node', pid: 41310, status: 'running', ip_address: '10.10.0.4/17', ports: [9000] },
    { container_id: 'd4e5f6a1b2c3d4e5f6a1b2c3', container_name: 'alpine-worker', image: 'alpine', pid: 41402, status: 'frozen', ip_address: '10.10.0.5/17', ports: [9000] },
  ],

  images: [
    { name: 'python', tag: 'latest', pulled_at: new Date(Date.now() - 3 * 86400000).toISOString() },
    { name: 'ubuntu', tag: 'latest', pulled_at: new Date(Date.now() - 9 * 86400000).toISOString() },
    { name: 'node', tag: 'latest', pulled_at: new Date(Date.now() - 1 * 86400000).toISOString() },
    { name: 'alpine', tag: 'latest', pulled_at: new Date(Date.now() - 12 * 86400000).toISOString() },
  ],

  runtime: {
    socket: '/var/run/lxr.sock',
    web_api: '/api',
    bridge: 'lxr0',
  },

  network: {
    bridge: 'lxr0',
    bridge_ip: '10.10.0.1',
    subnet: '10.10.0.0/17',
    usable_hosts: 32766,
    used: 4,
  },

  async _fakeStream(lines, { onLine, onComplete }, delay = 260) {
    for (const line of lines) {
      await new Promise(r => setTimeout(r, delay));
      onLine && onLine(line);
    }
    onComplete && onComplete();
  },

  simulateCreate(name, image, handlers) {
    return this._fakeStream([
      'creation started...',
      '',
      '[+] Find Image locally in LXR-registry...',
      '',
      '[+] Setting up container rootfs...',
      '',
      '[+] Building container environment with rootfs..',
      '',
      '[+] Setting up container networking..',
      '',
      '[+] Setting up container resources limit...',
      '',
      '[+] code-server activated at port 9000 ✔',
      '',
      `CONTAINER ID: ${Math.random().toString(16).slice(2, 18)}`,
      `CONTAINER NAME: ${name}`,
      '',
      'Container Created Successfully...',
    ], handlers);
  },

  simulateStart(name, handlers) {
    return this._fakeStream([
      '[+] Building container environment with rootfs..',
      '',
      '[+] Setting up container networking..',
      '',
      '[+] Setting up container resources limit...',
      '',
      'Container Running...',
    ], handlers);
  },

  async simulateStop(name) {
    await new Promise(r => setTimeout(r, 300));
    return { container_exists: true, container_Stoped: true };
  },

  async simulateKill(name) {
    await new Promise(r => setTimeout(r, 300));
    return { container_exists: true, is_killed: true };
  },

  simulatePull(image, handlers) {
    return this._fakeStream([
      '[+]Find Image locally in LXR-registry..',
      '',
      '[+] Initialize Image Pull...',
      '',
      '[+] Pulling image layers...',
      '',
      '      [+] layer 1 pulled ✔',
      '      [+] layer 2 pulled ✔',
      '',
      '[+] Image pull completed',
    ], handlers);
  },
};

window.LXRDemo = LXRDemo;
