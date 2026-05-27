# LXR — Linux Container Runtime

LXR is a Linux container runtime and browser-accessible development environment built from scratch using low-level Linux primitives such as namespaces, cgroups, veth networking, PTY execution, and isolated root filesystems.

LXR manually implements core container runtime components including image pulling, rootfs extraction, process isolation, container networking, resource control, and interactive terminal execution.

Each container includes integrated `code-server` support, enabling isolated browser-based development environments directly inside containers.

LXR directly orchestrates namespaces, cgroups, networking, PTY systems, and filesystem isolation using Go to create fully isolated container environments.

---

# Features

* Linux namespace based container isolation
* cgroup based resource control
* Custom bridge networking with veth pairs
* Custom O(1) IP allocator with reusable IP pools
* PTY based interactive shell execution
* Browser-accessible `code-server` containers
* Pull container images directly from Docker Hub
* Extract and manage isolated rootfs environments
* Unix socket based daemon communication
* Persistent container metadata
* CLI driven workflow

---

# Architecture Overview

```text
lxr-cli
    │
unix socket (/var/run/lxr.sock)
    │
LXR daemon
    │
Handlers
    │
Helpers
    │
Namespaces / cgroups / networking / rootfs
```

---

# Internal Architecture

LXR is split into multiple internal layers.

## Handler Layer

The handler layer exposes HTTP APIs over a Unix socket and acts as the control plane entrypoint.

Supported APIs:

* create
* start
* stop
* exec
* kill
* ps
* ps/all
* pull_image

---

## Helper Layer

The helper layer contains the orchestration logic responsible for container lifecycle management.

This layer handles:

* rootfs setup
* image extraction
* namespace lifecycle
* cgroup setup
* networking
* process management
* state restoration

---

## Networking Layer

LXR uses:

* Linux bridge networking
* veth pairs
* custom subnet allocation
* reusable O(1) IP allocation

Sample network configuration:

```env
NETWORK=10.10.0.0
CIDR=17
BRIDGE_IP=10.10.0.1
IP_START_RANGE=10.10.0.2
IP_END_RANGE=10.10.127.254
NETWORK_ADDR=10.10.0.0
BROADCAST_ADDR=10.10.127.255
TOTAL_USABLE_HOST=32766
```

Each container receives:

* isolated network namespace
* dedicated veth interface
* dynamically allocated IP
* bridge connectivity
* default gateway

Released container IPs are automatically returned back to the reusable IP pool when containers are destroyed.

---

## Execution Layer

Interactive terminal access is implemented using:

* PTY (pseudo terminal)
* nsenter
* namespace attachment

This allows containers to behave like real isolated terminal environments while preserving native terminal interaction.

---

# Project Structure

```text
LXR/
├── cmd/
│   └── server/
│       ├── main.go
│       └── routes.go
├── internal/
│   ├── app/
│   ├── handlers/
│   ├── helper/
│   ├── ip/
│   ├── models/
│   └── response/
├── script/
├── .env
└── Makefile
```

---

# Requirements

* Ubuntu 22.04
* Go
* Make
* Root privileges

---

# Local Setup

## Clone Repository

```bash
git clone https://github.com/jack-san-145/LXR.git
cd LXR
```

---

## Setup Environment

Create `.env`

```env
NETWORK=10.10.0.0
CIDR=17
BRIDGE_IP=10.10.0.1
IP_START_RANGE=10.10.0.2
IP_END_RANGE=10.10.127.254
NETWORK_ADDR=10.10.0.0
BROADCAST_ADDR=10.10.127.255
TOTAL_USABLE_HOST=32766
```

---

## Install Dependencies

```bash
sudo apt update

sudo apt install -y \
    bridge-utils \
    curl \
    jq \
    tar \
    iproute2 \
    uidmap
```

---

## Build LXR Daemon

```bash
make build
```

---

## Start LXR Daemon

```bash
sudo make run
```

LXR daemon communicates through:

```text
/var/run/lxr.sock
```

---

# Setup LXR CLI

LXR commands are executed using the separate CLI project.

LXR-cli Repository:

[https://github.com/jack-san-145/LXR-cli](https://github.com/jack-san-145/LXR-cli)

Clone and setup:

```bash
git clone https://github.com/jack-san-145/LXR-cli.git
cd LXR-cli
```

Build CLI:

```bash
make build
```

Copy binary:

```bash
sudo make cpbin
```

Verify installation:

```bash
lxr
```

---

# Commands

## Create Container

```bash
lxr create --name py-con python
```

Container creation workflow:

* pulls container image
* extracts root filesystem
* configures namespaces
* sets up networking
* configures cgroups
* installs container dependencies
* starts isolated environment

![Container Create](./assets/container-create.png)

---

## List Running Containers

```bash
lxr ps
```

---

## List All Containers

```bash
lxr ps --all
```

---

## Execute Shell Inside Container

```bash
lxr exec py-con
```

This attaches an interactive PTY shell into the container namespaces.

---

## Stop Container

```bash
lxr stop py-con
```

Container processes are frozen using cgroups.

---

## Start Container

```bash
lxr start py-con
```

---

## Kill Container

```bash
lxr kill py-con
```

This removes:

* rootfs
* veth pair
* namespaces
* cgroups
* container metadata

Released container IPs are returned back to the reusable IP pool.

---

# Container Environment

Each container automatically includes:

* nano
* git
* iproute2
* ping utilities
* `code-server`

Containers expose browser-accessible development environments through integrated `code-server` support.

---

# Example Workflow

```bash
# create container
lxr create --name go-con golang

# start container
lxr start go-con

# enter container
lxr exec go-con

# stop container
lxr stop go-con

# remove container
lxr kill go-con
```

---

# Storage Layout

LXR maintains:

```text
LXR-registry/
```

Stores downloaded image layers and extracted image rootfs.

```text
LXR-data/
```

Stores per-container isolated filesystems.

---
