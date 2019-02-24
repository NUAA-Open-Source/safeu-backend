#!/bin/bash

mkdir -p ../data/mariadb

sudo docker-compose -f ../deployments/production/docker-compose.yml up --force-recreate -d
