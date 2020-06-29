#!/bin/sh

# Prerequisites
# 1. Docker installed on remote machine
# 2. Docker Compose installed on remote machine

docker network create -d bridge nginx-proxy-network
docker-compose -f nginx-proxy-compose.yaml up -d || exit
docker-compose -f service-compose.yaml build --build-arg app_env=production || exit
docker-compose -f service-compose.yaml up -d || exit