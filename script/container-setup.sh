#!/bin/bash

set -ex

ROOT_FS=/home/jack/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/rootfs

mount --make-rprivate /

mount --bind $ROOT_FS $ROOT_FS
 
mkdir -p $ROOT_FS/dev

touch $ROOT_FS/dev/null

mount --bind /dev/null $ROOT_FS/dev/null

#ls -l $ROOT_FS

mkdir -p $ROOT_FS/old_root

pivot_root $ROOT_FS $ROOT_FS/old_root

cd /

mount -t proc proc /proc

umount -l /old_root

mount -t tmpfs tmpfs /tmp

hostname $HOSTNAME

# run container app
exec /bin/bash
