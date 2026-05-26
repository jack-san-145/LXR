#!/bin/bash

# set -ex

echo -e "[+]creating new veth pairs ✔"
#creationg new veth pairs 
ip link add $CONTAINER_VETH type veth peer $BRIDGE_VETH

echo -e "[+]attach to lxr bridge ✔"
#attach one side veth to bridge lxr0
ip link set $BRIDGE_VETH master lxr0

#bring bridge interface up
ip link set $BRIDGE_VETH up

#move another side veth to container from host
ip link set $CONTAINER_VETH netns $CONTAINER_PID

#allocate ip address for container
nsenter -t $CONTAINER_PID --net \
    ip addr add $CONTAINER_IP dev $CONTAINER_VETH

#bring container interface up
nsenter -t $CONTAINER_PID --net \
    ip link set $CONTAINER_VETH up

#bring conainer's loopback interface up
nsenter -t $CONTAINER_PID --net \
    ip link set lo up

echo -e "[+]add container default ghateway ✔"
#add default route as bridge to route unknown traffic
nsenter -t $CONTAINER_PID --net \
    ip route add default via $BRIDGE_IP














