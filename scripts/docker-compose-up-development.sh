#!/bin/bash

mkdir -p ../db-data

echo "on daemon ?(y/n) "
read isD

if [ "$isD" = "y" ]||[ "$isD" = "Y" ]
then
    sudo docker-compose -f ../deployments/development/docker-compose.yml up --force-recreate --build -d
else
    sudo docker-compose -f ../deployments/development/docker-compose.yml up --force-recreate --build
fi
