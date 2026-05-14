
#!/bin/bash

#make the output tracable
set -ex

#create seperate directory for each container with its containerName+containerID
mkdir -p /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID

#copy image rootfs Recursively to container data
cp -r /home/LXR/LXR-registry/$IMAGE_NAME /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/

#enter into it
cd /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/$IMAGE_NAME

#copy ping and ip binary from host to container to use ping and ip cmds inside container
cp /usr/bin/ping /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/$IMAGE_NAME/rootfs/usr/bin

cp /usr/bin/ip /home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/$IMAGE_NAME/rootfs/usr/bin


#to copy missing libraries from host to container to use ip cmd
#create an indexed array to store ip libraries
declare -a ip_libraries

container_rootfs=/home/LXR/LXR-data/$CONTAINER_NAME-$CONTAINER_ID/$IMAGE_NAME/rootfs

#append all libraries required for ip iperation to arr
$ip_libraries+=("/lib/aarch64-linux-gnu/libbpf.so.0")
$ip_libraries+=("/lib/aarch64-linux-gnu/libelf.so.1")
$ip_libraries+=("/lib/aarch64-linux-gnu/libmnl.so.0")
$ip_libraries+=("/lib/aarch64-linux-gnu/libbsd.so.0")
$ip_libraries+=("/lib/aarch64-linux-gnu/libcap.so.2")
$ip_libraries+=("/lib/aarch64-linux-gnu/libc.so.6")
$ip_libraries+=("/lib/aarch64-linux-gnu/libz.so.1")
$ip_libraries+=("/lib/aarch64-linux-gnu/libmd.so.0")

#loader that loads all the libraries dynamically
dynamic_loader=/lib/ld-linux-aarch64.so.1 


#use for loop to iterate ip_libraries
for library in $ip_libraries;do

    #check if the file exists or not inside container
    if [[ -f $container_rootfs+$library ]];then
        cp $library $container_rootfs/lib/aarch64-linux-gnu/
        echo "$library copied to container"
    else
        echo "file not found: $library"
    fi
done


#check if path exists inside container and copy dynamic_loader to container
if [[ -f $container_rootfs+$dynamic_loader ]];then
    cp $dynamic_loader $container_rootfs/lib/
    echo "$dynamic_loader copied to container"
else
    echo "loader not found: $dynamic_loader"
fi


#change the rootfs ownership
chown jack:lxr -R rootfs

#give permissions for /dev to create null further
chmod 775 rootfs/dev



