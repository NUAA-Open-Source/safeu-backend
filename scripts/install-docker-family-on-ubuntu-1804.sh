#!/bin/sh

# remove older versions
sudo apt-get remove docker docker-engine docker.io containerd runc -y

# Install docker via USTC mirror
# sudo apt-get update -y
# sudo apt-get install -y \
#     apt-transport-https \
#     ca-certificates \
#     curl \
#     gnupg-agent \
#     software-properties-common

# curl -fsSL https://mirrors.ustc.edu.cn/docker-ce/linux/ubuntu/gpg | sudo apt-key add -y -

# sudo add-apt-repository -y \
#    "deb [arch=amd64] https://mirrors.ustc.edu.cn/docker-ce/linux/ubuntu \
#    $(lsb_release -cs) \
#    stable"

# sudo apt-get update -y
# sudo apt-get install docker-ce docker-ce-cli containerd.io -y

# Install docker via DaoCloud mirror
curl -sSL https://get.daocloud.io/docker | sh

if [ ! -f "/etc/docker/daemon.json" ]; then
echo "{
    \"registry-mirrors\": [\"https://docker.mirrors.ustc.edu.cn/\"]
}   
" >> /etc/docker/daemon.json
fi

sudo systemctl restart docker

# Install docker compose via DaoCloud
curl -L https://get.daocloud.io/docker/compose/releases/download/1.23.2/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

echo ""
echo "Installed docker & docker-compose successfully~"
