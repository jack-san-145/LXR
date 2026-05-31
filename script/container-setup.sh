#!/bin/bash

#directory to store the container's data
ROOT_FS=/home/LXR/LXR-data/$CONTAINER_NAME/$IMAGE_NAME/merged

#after this mount wont affects the host 
mount --make-rprivate /

#make the same directory as mount point
mount --bind $ROOT_FS $ROOT_FS

#create directory named old_root to store the old root(host rootfs)
mkdir -p $ROOT_FS/old_root

#it changes the new rootfs as (/ root directory) and move host rootfs to old_root
pivot_root $ROOT_FS $ROOT_FS/old_root

cd /

#mount proc (virtual filesystem) ,sysfs,tmpfs
mount -t proc proc /proc
mount -t tmpfs tmpfs /tmp

#create fresh /dev inside the new rootfs
mount -t tmpfs tmpfs /dev

#mount the host's /dev/null to new rootfs /dev/null
touch /dev/null /dev/random /dev/urandom

mount --bind /old_root/dev/null /dev/null
mount --bind /old_root/dev/random /dev/random
mount --bind /old_root/dev/urandom /dev/urandom

#mount pts to /dev/pts and add symlink for ptmx
mkdir -p /dev/pts
mount -t devpts devpts /dev/pts
ln -sf /dev/pts/ptmx /dev/ptmx

#clear cache apt modules
rm -rf /var/lib/apt/lists/*

#make container name as container's hostname
hostname $CONTAINER_NAME

#clear old rootfs inside container
umount -l /old_root




