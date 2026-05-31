#!/bin/bash
export DEBIAN_FRONTEND=noninteractive  

echo -e "\n[+] Setting up $IMAGE_NAME rootfs..."


#first setup temporary container with host root privilages 
#directory to store the container's data
IMAGE_PATH=/home/LXR/LXR-registry/$IMAGE_NAME

mkdir -p $IMAGE_PATH/modified_rootfs

cp -a $IMAGE_PATH/rootfs/. $IMAGE_PATH/modified_rootfs


#change rootfs to modified rootfs
ROOT_FS=$IMAGE_PATH/modified_rootfs

#after this mount wont affects the host 
mount --make-rprivate /

#make the same directory as mount point
mount --bind $ROOT_FS $ROOT_FS

#create directory named old_root to store the old root(host rootfs)
mkdir -p $ROOT_FS/old_root

#it changes the new rootfs as (/ root directory) and move host rootfs to old_root
pivot_root $ROOT_FS $ROOT_FS/old_root

cd /

#mount proc (virtual filesystem) ,sysfs,tmpfs
mount -t proc proc /proc
mount -t sysfs sysfs /sys

#create fresh /dev inside the new rootfs
mount -t tmpfs tmpfs /dev

#mount the host's /dev/null to new rootfs /dev/null
touch /dev/null
mount --bind /old_root/dev/null /dev/null

#mount pts to /dev/pts and add symlink for ptmx
mkdir -p /dev/pts
mount -t devpts devpts /dev/pts
ln -sf /dev/pts/ptmx /dev/ptmx

#clear cache apt modules
rm -rf /var/lib/apt/lists/*

#clear old rootfs inside container
umount -l /old_root



#install dependencies to the actuall image rootfs
echo -e "\n[+] Installing dependencies..."

#mount temporary fs
mount -t tmpfs tmpfs /tmp 

#change mode for /tmp
chmod 1777 /tmp 

#update apt 
apt-get update > /dev/null 2>&1

#install nano editor
apt-get install -y nano > /dev/null 2>&1
echo "  [+] nano ✔"

#install ip package to use all the iptables oerations
apt-get install -y iproute2 > /dev/null 2>&1
echo "  [+] iproute2 ✔"

#install ping globally orelse between containers
apt-get install -y iputils-ping > /dev/null 2>&1
echo "  [+] iputils-ping ✔"


#install git
apt-get install -y git > /dev/null 2>&1
echo "  [+] git ✔"


echo "  [+] Installing code-server..."
echo "  [+] This may take a few minutes. Please wait while the installation completes."

#install code server to use vscode on web
curl -fsSL https://code-server.dev/install.sh | sh 

echo "  [+] code-server ✔"
echo "  [+] Dependency Installation completed ✔"