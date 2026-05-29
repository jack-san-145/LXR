#!/bin/bash
export DEBIAN_FRONTEND=noninteractive  

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

#start code-server in background 
code-server --bind-addr 0.0.0.0:9000 &

#immediately store code-server's background pid to kill
cs_pid=$!

#wait for 2 seconds to ensure server run
sleep 2

#kill code-server pid 
kill -2 $cs_pid



#when code-server starts that created directory with config.yaml 
#so after that code-server pid killed and now change code-server's password to config file

echo -e "\n  [+] configuring code-server.."

cat << EOF > ~/.config/code-server/config.yaml
bind-addr: 127.0.0.1:9000
auth: password
password: $PASSWORD
cert: false
EOF


#run again code-server in background and redirects its output,errors to null
nohup code-server --bind-addr 0.0.0.0:9000 > /dev/null 2>&1 &
echo "  [+] code-server ✔"

sleep 2 

#create both randam and urandom for random number generation
(mknod -m 666 /dev/random c 1 8 || true) > /dev/null 2>&1
(mknod -m 666 /dev/urandom c 1 9 || true) > /dev/null 2>&1

echo "  [+] Dependency Installation completed ✔"