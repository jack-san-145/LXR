#!/bin/bash

set -ex

#container data directory
con_dir="/home/LXR/LXR-data/$CONTAINER_NAME"

#remove bridge veth and another veth automatically removed
ip link delete $BRIDGE_VETH 

#kill container's parent process(unshare)
kill -9 $CONTAINER_PID

sleep 2

#umount container's overlay setup
umount -l $con_dir/$IMAGE_NAME/merged

#remove entire container directory
rm -rf $con_dir



