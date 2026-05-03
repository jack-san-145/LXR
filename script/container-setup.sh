#!/bin/bash

set -ex

#directory to store the container's data
ROOT_FS=/home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/$IMAGE/rootfs

#after this mount wont affects the host 
mount --make-rprivate /

#make the same directory as mount point
mount --bind $ROOT_FS $ROOT_FS
 
#create /dev inside the new rootfs
mkdir -p $ROOT_FS/dev

#create null inside new rootfs
touch $ROOT_FS/dev/null

#mount the host's /dev/null to new rootfs /dev/null
mount --bind /dev/null $ROOT_FS/dev/null

#create directory named old_root to store the old root(host rootfs)
mkdir -p $ROOT_FS/old_root

#it changes the new rootfs as (/ root directory) and move host rootfs to old_root
pivot_root $ROOT_FS $ROOT_FS/old_root

cd /

#mount proc (virtual filesystem)
mount -t proc proc /proc

#clear old rootfs inside container
umount -l /old_root

#mount temporary fs to /tmp
mount -t tmpfs tmpfs /tmp

#root process of the container
exec sleep infinity
