#!/bin/bash

set -ex 

#mount temporary fs
mount -t tmpfs tmpfs /tmp 

#change mode for /tmp
chmod 1777 /tmp 

#create both randam and urandom for random number generation
mknod -m 666 /dev/random c 1 8
mknod -m 666 /dev/urandom c 1 9
