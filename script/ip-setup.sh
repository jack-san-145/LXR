#!/bin/bash

set -e



echo "   [+] creating new veth pairs ✔"

#delete if container_veth already exists
ip link delete $CONTAINER_VETH || true

#creationg new veth pairs 
ip link add $CONTAINER_VETH type veth peer $BRIDGE_VETH


#attach one side veth to bridge lxr0
ip link set $BRIDGE_VETH master lxr0
echo "   [+] attach to lxr bridge ✔"

#bring bridge interface up
ip link set $BRIDGE_VETH up

#move another side veth to container from host
ip link set $CONTAINER_VETH netns $CONTAINER_PID
echo "   [+] attach to container ✔"


#allocate ip address for container
nsenter -t $CONTAINER_PID --net \
    ip addr add $CONTAINER_IP dev $CONTAINER_VETH
echo "   [+] allocate IP address ✔"


#bring container interface up
nsenter -t $CONTAINER_PID --net \
    ip link set $CONTAINER_VETH up


#bring conainer's loopback interface up
nsenter -t $CONTAINER_PID --net \
    ip link set lo up
echo "   [+] bring container interface up ✔"


#add default route as bridge to route unknown traffic
nsenter -t $CONTAINER_PID --net \
    ip route add default via $BRIDGE_IP
echo "   [+] add container default ghateway ✔"


echo "  [+] Network setup completed ✔"