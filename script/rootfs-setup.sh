#!/bin/bash

#make the output tracable
set -ex

#create seperate directory for each container with its containerName+containerID
mkdir -p /home/jack/LXR-data/$CONTAINER_NAME-$CONTAINER_ID

#enter into it
cd /home/jack/LXR-data/$CONTAINER_NAME-$CONTAINER_ID

#make rootfs directory
mkdir rootfs

#Create temporary container
docker create --name $CONTAINER_NAME $IMAGE_NAME

#extract container's rootfs
docker export $CONTAINER_NAME > $CONTAINER_NAME-roofs.tar

#extract the tar to get actual rootfs
tar -xf $CONTAINER_NAME-roofs.tar -C rootfs

#change the rootfs ownership
chown jack:lxr -R rootfs

#give permissions for /dev to create null further
chmod 775 rootfs/dev

#remove created temporary container
docker rm $CONTAINER_NAME

#remove rootfs tar
rm -rf $CONTAINER_NAME-roofs.tar
