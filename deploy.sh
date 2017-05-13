#!/usr/bin/env bash

cd ~/go/socketio-server && go build -ldflags "-X main.Env=docker" -o ~/go/bin/socketio-server
docker stop $(docker ps -a -q)
# remove container
docker rm $(docker ps -a -q --filter ancestor=zhanat87/golang-socketio-server) -f
docker rmi $(docker images --filter=reference='zhanat87/golang-socketio-server') -f

# create new docker image, push to docker hub and pull
docker build -t zhanat87/golang-socketio-server .
docker push zhanat87/golang-socketio-server
docker pull zhanat87/golang-socketio-server
# list of all docker images on host machine
docker images

curl http://zhanat.site:9000/hooks/golang-socketio-server

echo "deploy success"
