#!/bin/bash

set +x

printf "Startup script - $1 \n"

IP_PREFIX="10.20.1"

# Install Git
sudo apt update -y
sudo apt install -y git

GO_VERSION="1.21.6"
DOCKER_COMPOSE_VERSION="2.24.5"

# Install Go
wget https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

# Create the GOPATH directory
mkdir -p $HOME/go $HOME/go/bin $HOME/go/src

# Install Docker
sudo apt-get remove docker docker-engine docker.io containerd runc
sudo apt-get update -y
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update -y
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# Add the user to the docker group
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/$(DOCKER_COMPOSE_VERSION)/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify installations
echo "Verifying installations..."
git --version
go version
docker --version
docker-compose --version

printf "Installing redis \n"
sudo apt update
sudo apt install redis-server -y
node_num="$1"
if [ "$node_num" == "0" ]; then

elif [ "$node_num" == "1" ]; then

elif [ "$node_num" == "2" ]; then

else

fi

printf "%s: %s\n" "$(date +"%T.%N")" "Redis setup completed!"