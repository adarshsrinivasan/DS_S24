#!/bin/bash

set +x

# shellcheck disable=SC2059
printf "Startup script - $1 \n"

IP_PREFIX="10.20.1"
GO_VERSION="1.21.6"

# Install Packages
sudo apt update -y
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    git \
    golang-go

# Setup Go
#wget https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz -O /tmp/go$GO_VERSION.linux-amd64.tar.gz
#sudo tar -C /usr/local -xzf /tmp/go$GO_VERSION.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
# shellcheck disable=SC1090
source ~/.bashrc
mkdir -p $HOME/go $HOME/go/bin $HOME/go/src

# Install Docker
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove $pkg; done
sudo apt-get update -y
sudo apt-get -y install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update -y

sudo apt-get -y install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Add the user to the docker group
sudo groupadd docker || true
sudo usermod -aG docker $USER
newgrp docker

# Verify installations
echo "Verifying installations..."
git --version
go version
docker --version
docker-compose --version

node_num="$1"
if [ "$node_num" == "0" ]; then
  echo "Launching Postgres..."
  docker compose -f /local/repository/Assignment1/deployment/docker/docker-compose.yaml up -d postgres
  docker compose -f /local/repository/Assignment1/deployment/docker/docker-compose.yaml ps
  echo "Postgres Launch Complete..."
elif [ "$node_num" == "1" ]; then
  echo "Launching MongoDB..."
  docker compose -f /local/repository/Assignment1/deployment/docker/docker-compose.yaml up -d mongodb
  docker compose -f /local/repository/Assignment1/deployment/docker/docker-compose.yaml ps
  echo "MongoDB Launch Complete..."
elif [ "$node_num" == "2" ]; then
  echo "Launching Server-Seller..."
  # shellcheck disable=SC2164
  cd /local/repository/Assignment1/cmd/server/
  rm server-seller || true
  go build -o server-seller .
  SERVER_HOST=$IP_PREFIX.3 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace ./server-seller
  echo "Server-Seller Launch Complete..."
elif [ "$node_num" == "3" ]; then
  echo "Launching Server-Buyer..."
  # shellcheck disable=SC2164
  cd /local/repository/Assignment1/cmd/server/
  rm server-buyer || true
  go build -o server-buyer .
  SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace ./server-buyer
  echo "Server-Buyer Launch Complete..."
elif [ "$node_num" == "4" ]; then
  echo "Launching Client-Seller..."
  # shellcheck disable=SC2164
  cd /local/repository/Assignment1/cmd/seller/
  rm client-seller || true
  go build -o client-seller .
  #SERVER_HOST=$IP_PREFIX.3 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace ./client-seller
  echo "Client-Seller Launch Complete..."
elif [ "$node_num" == "5" ]; then
  echo "Launching Client-Buyer..."
  # shellcheck disable=SC2164
  cd /local/repository/Assignment1/cmd/buyer/
  rm client-buyer || true
  go build -o client-buyer .
  #SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace ./client-buyer
  echo "Client-Buyer Launch Complete..."
elif [ "$node_num" == "6" ]; then
  echo "Launching Test-Client-Buyer-Latency..."
  # shellcheck disable=SC2164
  cd /local/repository/Assignment1/cmd/test_buyer_response/
  rm test-buyer-latency || true
  go build -o test-buyer-latency .
  #SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace ./test-buyer-latency
  echo "Test-Client-Buyer-Latency Launch Complete..."
elif [ "$node_num" == "7" ]; then
  echo "Launching Test-Client-Buyer-Throughput..."
  # shellcheck disable=SC2164
  cd /local/repository/Assignment1/cmd/test_buyer_throughput/
  rm test-buyer-throughput || true
  go build -o test-buyer-throughput .
  #SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace ./test-buyer-throughput
  echo "Test-Client-Buyer-Throughput Launch Complete..."
else
  echo "Invalid Node Number $node_num"
fi

printf "%s: %s\n" "$(date +"%T.%N")" "Setup completed!"