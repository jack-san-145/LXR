#!/bin/bash

#make the output tracable
# set -ex

#create seperate directory for each container with its containerName
mkdir -p /home/LXR/LXR-data/$CONTAINER_NAME

echo "  [+] extract image rootfs ✔"
#copy image rootfs Recursively to container data
cp -r /home/LXR/LXR-registry/$IMAGE_NAME /home/LXR/LXR-data/$CONTAINER_NAME/


echo "  [+] add dependency script ✔"
mkdir -p /home/LXR/LXR-data/$CONTAINER_NAME/$IMAGE_NAME/rootfs/home/script
cp /home/jack/LXR/LXR-d/script/install-dependencies.sh /home/LXR/LXR-data/$CONTAINER_NAME/$IMAGE_NAME/rootfs/home/script

#enter into it
cd /home/LXR/LXR-data/$CONTAINER_NAME/$IMAGE_NAME

#change the rootfs ownership
chown jack:lxr -R rootfs

#give permissions for /dev to create null further
chmod 775 rootfs/dev

echo "  [+] Rootfs setup completed ✔"


