#!/bin/bash

set -e

image_path="/home/LXR/LXR-registry/$IMAGE_NAME"
container_path="/home/LXR/LXR-data/$CONTAINER_NAME/$IMAGE_NAME"

# create overlay directories
mkdir -p "$container_path"/{upper,work,merged}


# Create container rootfs using OverlayFS (base + writable layer)
# - lowerdir: base image rootfs (read-only)
# - upperdir: container's writable layer (stores changes)
# - workdir: internal workspace used by kernel
# - merged: final combined rootfs used by container
mount -t overlay overlay \
  -o lowerdir="$image_path/rootfs",\
upperdir="$container_path/upper",\
workdir="$container_path/work" \
  "$container_path/merged"

echo "  [+] overlay filesystem ✔"

# add dependency script into container
mkdir -p "$container_path/merged/home/script"
cp "/home/jack/LXR/LXR-d/script/install-dependencies.sh" \
   "$container_path/merged/home/script"

echo "  [+] add dependency script ✔"

#create directory /dev
mkdir -p "$container_path/merged"/dev

mkdir -p $container_path/upper/old_root

echo "  [+] Rootfs setup completed ✔"



