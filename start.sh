#!/bin/bash

set +x

IP_PREFIX="10.20.1"
PORT="6379"
PASSWORD="password"

printf "Startup script to setup Redis - $1 $2 \n"

printf "Installing redis \n"
sudo apt update
sudo apt install redis-server -y
node_num="$2"
if [ "$node_num" == "0" ]; then

elif [ "$node_num" == "1" ]; then

elif [ "$node_num" == "2" ]; then

else

fi

printf "%s: %s\n" "$(date +"%T.%N")" "Redis setup completed!"