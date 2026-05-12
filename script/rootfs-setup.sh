#!/bin/bash

#make the output tracable
set -ex

#create seperate directory for each container with its containerName+containerID
mkdir -p /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID

#copy image rootfs Recursively to container data
cp -r /home/LXR/LXR-registry/$IMAGE_NAME /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/

#enter into it
cd /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/$IMAGE_NAME

#copy ping and ip binary from host to container to use ping and ip cmds inside container
cp /usr/bin/ping /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/$IMAGE_NAME/rootfs/usr/bin

cp /usr/bin/ip /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/$IMAGE_NAME/rootfs/usr/bin

#copy missed libraries from host to container to use ip cmd
cp /lib/aarch64-linux-gnu/libbpf.so.0 rootfs/lib/aarch64-linux-gnu/
cp /lib/aarch64-linux-gnu/libelf.so.1 rootfs/lib/aarch64-linux-gnu/
cp /lib/aarch64-linux-gnu/libmnl.so.0 rootfs/lib/aarch64-linux-gnu/

#change the rootfs ownership
chown jack:lxr -R rootfs

#give permissions for /dev to create null further
chmod 775 rootfs/dev



