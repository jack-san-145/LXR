#!/bin/bash

image_path=/home/LXR/LXR-registry/$IMAGE_NAME
container_path=/home/LXR/LXR-data/$CONTAINER_NAME/$IMAGE_NAME

#get user with id=1000
user=$(id -nu 1000)

#create seperate directory for container with overlayfs directories
mkdir -p $container_path/{upper,work,merged}


# Create container rootfs using OverlayFS (base + writable layer)
# - lowerdir: base image rootfs (read-only)
# - upperdir: container's writable layer (stores changes)
# - workdir: internal workspace used by kernel
# - merged: final combined rootfs used by container
mount -t overlay overlay \
  -o lowerdir=$image_path/rootfs,\
upperdir=$container_path/upper,\
workdir=$container_path/work \
  $container_path/merged

echo "  [+] overlay filesystem ✔"


mkdir -p $container_path/merged/home/script
cp /home/jack/LXR/LXR-d/script/install-dependencies.sh $container_path/merged/home/script

echo "  [+] add dependency script ✔"


#change the rootfs ownership
chown $user:lxr $container_path/upper

#give permissions for /dev to create null further
mkdir -p $container_path/merged/dev
chmod 775 $container_path/merged/dev

echo "  [+] Rootfs setup completed ✔"


