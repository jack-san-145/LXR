#!/bin/bash

set -e

BRIDGE="lxr0"
SUBNET="10.10.0.0/17"
HOST_IFACE="enp0s1"

echo "[LXR] Initializing..."

# Bridge
if ! ip link show "$BRIDGE" >/dev/null 2>&1; then
    echo "[LXR] Creating bridge $BRIDGE"
    ip link add "$BRIDGE" type bridge
fi

ip addr add 10.10.0.1/17 dev "$BRIDGE" 2>/dev/null || true
ip link set "$BRIDGE" up

# IPv4 forwarding
echo "[LXR] Enabling IPv4 forwarding"
sysctl -w net.ipv4.ip_forward=1 >/dev/null

# NAT
iptables -t nat -C POSTROUTING \
    -s "$SUBNET" ! -o "$BRIDGE" -j MASQUERADE 2>/dev/null || \
iptables -t nat -A POSTROUTING \
    -s "$SUBNET" ! -o "$BRIDGE" -j MASQUERADE

# Forward: container -> host
iptables -C FORWARD \
    -i "$BRIDGE" -o "$HOST_IFACE" -j ACCEPT 2>/dev/null || \
iptables -A FORWARD \
    -i "$BRIDGE" -o "$HOST_IFACE" -j ACCEPT

# Forward: host -> container
iptables -C FORWARD \
    -i "$HOST_IFACE" -o "$BRIDGE" \
    -m conntrack --ctstate ESTABLISHED,RELATED \
    -j ACCEPT 2>/dev/null || \
iptables -A FORWARD \
    -i "$HOST_IFACE" -o "$BRIDGE" \
    -m conntrack --ctstate ESTABLISHED,RELATED \
    -j ACCEPT

# cgroups
echo "[LXR] Setting up cgroups"

mkdir -p /sys/fs/cgroup/lxr

echo "+cpu +memory +pids" \
    > /sys/fs/cgroup/cgroup.subtree_control

echo "+cpu +memory +pids" \
    > /sys/fs/cgroup/lxr/cgroup.subtree_control

echo "[LXR] Initialization complete"