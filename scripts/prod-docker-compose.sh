#!/bin/bash
# Author:   TripleZ<me@triplez.cn>
# Date:     2019-03-11

echo -e "\n Build, up, down, check logs for SafeU production docker clusters.\n"

if [ "$1" == "up" ]
then
    mkdir -p ../data/mariadb
    sudo docker-compose -f ../deployments/production/docker-compose.yml up -d

elif [ "$1" == "down" ]
then
    sudo docker-compose -f ../deployments/production/docker-compose.yml down

elif [ "$1" == "build" ]
then
    sudo docker-compose -f ../deployments/production/docker-compose.yml build --force-rm

elif [ "$1" == "help" ] || [ "$1" == "-h" ] || [ "$1" == "--help" ]
then
    echo -e " Usage:
  ./prod-docker-compose.sh [COMMAND]
  ./prod-docker-compose.sh -h|--help

 Commands:
   build      Build SafeU prod container images.
   down       Down SafeU prod containers.
   help       Show this help message.
   logs       View output from prod containers.
   up         Up SafeU prod containers with force recreate and build.
"

else
    echo -e " Cannot match the command \"$1\", please type \"help\" command for help."
fi




