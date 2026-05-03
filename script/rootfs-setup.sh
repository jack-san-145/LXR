#!/bin/bash

#make the output tracable
set -ex

#create seperate directory for each container with its containerName+containerID
mkdir -p /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID

#copy image rootfs Recursively to container data
cp -r /home/LXR/LXR-registry/$IMAGE_NAME /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/

#enter into it
cd /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/$IMAGE_NAME

#change the rootfs ownership
chown jack:lxr -R rootfs

#give permissions for /dev to create null further
chmod 775 rootfs/dev
