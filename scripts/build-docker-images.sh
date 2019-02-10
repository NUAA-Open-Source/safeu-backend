#!/bin/sh

# Build safeu-backend-dev docker image
sudo docker build -t safeu-backend-dev:latest -f ../build/package/safeu-backend-dev/Dockerfile ..
