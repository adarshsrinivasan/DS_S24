#!/bin/bash

set +x

# shellcheck disable=SC2059
printf "Startup script - $1 $2 \n"

IP_PREFIX="10.20.1"
GO_VERSION="1.19"
max_attempts=10
delay_seconds=10

# Install Packages
sudo apt update -y
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    git \
    protobuf-compiler

# Setup Go
sudo rm -rf /usr/local/go /tmp/go$GO_VERSION.linux-amd64.tar.gz || true
wget https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz -O /tmp/go$GO_VERSION.linux-amd64.tar.gz
sudo tar -C /usr/local/ -xzf /tmp/go$GO_VERSION.linux-amd64.tar.gz
rm /tmp/go$GO_VERSION.linux-amd64.tar.gz
(echo 'export PATH=$PATH:/usr/local/go/bin'; echo 'export GOPATH=$HOME/go'; echo 'export PATH=$PATH:$GOPATH/bin') >> ~/.bashrc
# shellcheck disable=SC1090
source ~/.bashrc
mkdir -p $HOME/go $HOME/go/bin $HOME/go/src

#Install protoc go libraries
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

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
docker compose --version

node_num="$1"
assignment_num="$2"
if [ "$node_num" == "0" ]; then
  echo "Launching Postgres..."
  if [ "$assignment_num" == "1" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment1
      make run-sql-server
  elif [ "$assignment_num" == "2" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment2
      IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.1 SERVER_PORT=50000  POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-sql-server
  else
    echo "Error: unrecognized assignment number: $assignment_num"
    exit 1
  fi
  echo "Postgres Launch Complete..."
elif [ "$node_num" == "1" ]; then
  echo "Launching MongoDB..."
  if [ "$assignment_num" == "1" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment1
      make run-nosql-server
  elif [ "$assignment_num" == "2" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment2
      IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.2 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace make run-nosql-server
  else
    echo "Error: unrecognized assignment number: $assignment_num"
    exit 1
  fi
  echo "MongoDB Launch Complete..."
elif [ "$node_num" == "2" ]; then
  echo "Launching Server-Seller..."
  if [ "$assignment_num" == "1" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment1
      for ((attempt=1; attempt<=max_attempts; attempt++)); do
          echo "Attempt $attempt"
          IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.3 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace  POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-server-seller

          # Check the exit status
          exit_status=$?

          if [ $exit_status -eq 0 ]; then
              echo "Command succeeded!"
              break
          else
              echo "Command failed. Retrying in $delay_seconds seconds..."
              sleep $delay_seconds
          fi
      done

  elif [ "$assignment_num" == "2" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment2
      for ((attempt=1; attempt<=max_attempts; attempt++)); do
          echo "Attempt $attempt"
          IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.3 SERVER_PORT=50000 NOSQL_RPC_HOST=$IP_PREFIX.2 NOSQL_RPC_PORT=50000  SQL_RPC_HOST=$IP_PREFIX.1 SQL_RPC_PORT=50000 TRANSACTION_HOST=$IP_PREFIX.7 TRANSACTION_PORT=50000 make run-server-seller

          # Check the exit status
          exit_status=$?

          if [ $exit_status -eq 0 ]; then
              echo "Command succeeded!"
              break
          else
              echo "Command failed. Retrying in $delay_seconds seconds..."
              sleep $delay_seconds
          fi
      done
  else
    echo "Error: unrecognized assignment number: $assignment_num"
    exit 1
  fi
  echo "Server-Seller Launch Complete..."
elif [ "$node_num" == "3" ]; then
  echo "Launching Server-Buyer..."
  if [ "$assignment_num" == "1" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment1
      for ((attempt=1; attempt<=max_attempts; attempt++)); do
          echo "Attempt $attempt"
          IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.3 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace  POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-server-buyer

          # Check the exit status
          exit_status=$?

          if [ $exit_status -eq 0 ]; then
              echo "Command succeeded!"
              break
          else
              echo "Command failed. Retrying in $delay_seconds seconds..."
              sleep $delay_seconds
          fi
      done
  elif [ "$assignment_num" == "2" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment2
      for ((attempt=1; attempt<=max_attempts; attempt++)); do
          echo "Attempt $attempt"
          IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 NOSQL_RPC_HOST=$IP_PREFIX.2 NOSQL_RPC_PORT=50000  SQL_RPC_HOST=$IP_PREFIX.1 SQL_RPC_PORT=50000 TRANSACTION_HOST=$IP_PREFIX.7 TRANSACTION_PORT=50000 make run-server-buyer

          # Check the exit status
          exit_status=$?

          if [ $exit_status -eq 0 ]; then
              echo "Command succeeded!"
              break
          else
              echo "Command failed. Retrying in $delay_seconds seconds..."
              sleep $delay_seconds
          fi
      done
  else
    echo "Error: unrecognized assignment number: $assignment_num"
    exit 1
  fi
  echo "Server-Buyer Launch Complete..."
elif [ "$node_num" == "4" ]; then
  echo "Launching Client-Seller..."
  if [ "$assignment_num" == "1" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment1
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.3 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-client-seller
  elif [ "$assignment_num" == "2" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment2
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.3 SERVER_PORT=50000 make run-client-seller
  else
    echo "Error: unrecognized assignment number: $assignment_num"
    exit 1
  fi
  echo "Client-Seller Launch Complete..."
elif [ "$node_num" == "5" ]; then
  echo "Launching Client-Buyer..."
  if [ "$assignment_num" == "1" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment1
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-client-buyer
  elif [ "$assignment_num" == "2" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment2
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-client-buyer
  else
    echo "Error: unrecognized assignment number: $assignment_num"
    exit 1
  fi
  echo "Client-Buyer Launch Complete..."
elif [ "$node_num" == "6" ]; then
    echo "Launching Transaction-Server..."
    if [ "$assignment_num" == "1" ]; then
        # shellcheck disable=SC2164
        cd /local/repository/Assignment1
        #NO-OP :)
    elif [ "$assignment_num" == "2" ]; then
        # shellcheck disable=SC2164
        cd /local/repository/Assignment2
        IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.7 SERVER_PORT=50000 make run-server-transaction
    else
      echo "Error: unrecognized assignment number: $assignment_num"
      exit 1
    fi
    echo "Transaction-Server Launch Complete..."
elif [ "$node_num" == "7" ]; then
  echo "Setting-up Test-Buyer..."
  if [ "$assignment_num" == "1" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment1
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-test-latency 1 0
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-test-latency 10 0
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-test-latency 100 0
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-test-throughput 1 0
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-test-throughput 10 0
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-test-throughput 100 0
      echo ":)"
  elif [ "$assignment_num" == "2" ]; then
      # shellcheck disable=SC2164
      cd /local/repository/Assignment2
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-test-latency 1 0
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-test-latency 10 0
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-test-latency 100 0
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-test-throughput 1 0
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-test-throughput 10 0
      #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-test-throughput 100 0
      echo ":)"
  else
    echo "Error: unrecognized assignment number: $assignment_num"
    exit 1
  fi
  echo "Test-Buyer Setup Complete..."
elif [ "$node_num" == "8" ]; then
  echo "Setting-up Test-Seller.."
  if [ "$assignment_num" == "1" ]; then
        # shellcheck disable=SC2164
        cd /local/repository/Assignment1
        #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-test-latency 1 1
        #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-test-latency 10 1
        #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-test-latency 100 1
        #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-test-throughput 1 1
        #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-test-throughput 10 1
        #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 MONGO_HOST=$IP_PREFIX.2 MONGO_PORT=27017 MONGO_USERNAME=admin MONGO_PASSWORD=admin MONGO_DB=marketplace POSTGRES_HOST=$IP_PREFIX.1 POSTGRES_PORT=5432 POSTGRES_USERNAME=admin POSTGRES_PASSWORD=admin POSTGRES_DB=marketplace POSTGRES_MAX_CONN=500 make run-test-throughput 100 1
        echo ":)"
    elif [ "$assignment_num" == "2" ]; then
        # shellcheck disable=SC2164
        cd /local/repository/Assignment2
        #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-test-latency 1 1
        #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-test-latency 10 1
        #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-test-latency 100 1
        #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-test-throughput 1 1
        #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-test-throughput 10 1
        #IP_PREFIX="10.20.1" SERVER_HOST=$IP_PREFIX.4 SERVER_PORT=50000 make run-test-throughput 100 1
        echo ":)"
    else
      echo "Error: unrecognized assignment number: $assignment_num"
      exit 1
    fi
  echo "Testr-Seller Setup Complete..."
else
  echo "Invalid Node Number $node_num"
fi

printf "%s: %s\n" "$(date +"%T.%N")" "Setup completed!"