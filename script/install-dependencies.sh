#!/bin/bash

set -ex 

#mount temporary fs
mount -t tmpfs tmpfs /tmp 

#change mode for /tmp
chmod 1777 /tmp 

#create both randam and urandom for random number generation
mknod -m 666 /dev/random c 1 8
mknod -m 666 /dev/urandom c 1 9

#update apt 
apt update 

#install nano editor
apt install nano

#install ip package to use all the iptables oerations
apt install iproute2 -y 

#install ping globally orelse between containers
apt install iputils-ping 

#install code server to use vscode on web
curl -fsSL https://code-server.dev/install.sh | sh 

#start code-server in background 
code-server --bind-addr 0.0.0.0:8080 &

#immediately store code-server's background pid to kill
cs_pid=$!

#wait for 2 seconds to ensure server run
sleep 2

#kill code-server pid 
kill -2 $cs_pid

#when code-server starts that created directory with config.yaml 
#so after that code-server pid killed and now change code-server's password to config file

>  ~/.config/code-server/config.yaml cat << EOF
bind-addr: 127.0.0.1:8080
auth: password
password: $PASSWORD
cert: false
EOF

