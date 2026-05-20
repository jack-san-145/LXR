#!/bin/bash

set -ex

#container data directory
con_dir="/home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID"

#remove bridge veth and another veth automatically removed
ip link delete $BRIDGE_VETH

#remove entire container directory
rm -rf $con_dir

#kill container's parent process(unshare)
kill -9 $CONTAINER_PID
