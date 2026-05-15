#!/bin/bash

set -ex

#creationg new veth pairs 
ip link add $CONTAINER_VETH type veth peer $BRIDGE_VETH

#attach one side veth to bridge lxr0
ip link set $BRIDGE_VETH master lxr0

#bring bridge interface up
ip link set $BRIDGE_VETH up

#move another side veth to container from host
ip link set $CONTAINER_VETH netns $CONTAINER_PID
















