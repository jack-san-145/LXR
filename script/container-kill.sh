#!/bin/bash

set -ex

#container data directory
con_dir="/home/LXR/LXR-data/$CONTAINER_NAME"

#remove bridge veth and another veth automatically removed
ip link delete $BRIDGE_VETH 2>/dev/null || true

#kill container's parent process(unshare)
kill -9 $CONTAINER_PID || true

# wait until fully dead
while kill -0 $CONTAINER_PID 2>/dev/null; do
  sleep 0.5
done

#umount container's overlay setup with recursive lazy mount
umount -lR $con_dir/$IMAGE_NAME/merged || true

#remove entire container directory
rm -rf $con_dir



