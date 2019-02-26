#!/bin/sh

# Build safeu-backend-dev docker image for development
sudo docker build -t safeu-backend-dev:latest -f ../build/package/safeu-backend-dev/Dockerfile ..

# Build safeu-backend docker image for production
sudo docker build -t safeu-backend:latest -f ../build/package/safeu-backend/Dockerfile ..
