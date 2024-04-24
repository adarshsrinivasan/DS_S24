#!/bin/sh

docker-entrypoint.sh mongod &

sleep 15

/app/main

