#!/bin/bash

set -ex

mkdir -p /home/jack/LXR-data/$CONTAINER_NAME-$CONTAINER_ID

cd /home/jack/LXR-data/$CONTAINER_NAME-$CONTAINER_ID

mkdir rootfs

docker create --name $CONTAINER_NAME $IMAGE_NAME

docker export $CONTAINER_NAME > $CONTAINER_NAME-roofs.tar

sudo tar -xf $CONTAINER_NAME-roofs.tar -C rootfs

#sudo chown $USER:lxr -R rootfs

sudo chmod 775 rootfs/dev

#ls -l rootfs

docker rm $CONTAINER_NAME

rm -rf $CONTAINER_NAME-roofs.tar
