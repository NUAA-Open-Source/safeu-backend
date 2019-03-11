#!/bin/bash
# Author:   TripleZ<me@triplez.cn>
# Date:     2019-03-11

echo -e "\n Build, up, down, check logs for SafeU development docker clusters.\n"

if [ "$1" == "up" ]
then
    mkdir -p ../db-data
    echo -e " Running dockers on daemon mode? (y/n, default: n): \c"
    read isD

    if [ "$isD" = "y" ]||[ "$isD" = "Y" ]
    then
        sudo docker-compose -f ../deployments/development/docker-compose.yml up -d
    else
        sudo docker-compose -f ../deployments/development/docker-compose.yml up
    fi

elif [ "$1" == "down" ]
then
    sudo docker-compose -f ../deployments/development/docker-compose.yml down

elif [ "$1" == "build" ]
then
    sudo docker-compose -f ../deployments/development/docker-compose.yml build --force-rm

elif [ "$1" == "logs" ]
then
    echo -e " Follow log output? (y/n, default: n): \c"
    read isF
    echo ""
    if [ "$isF" == "y" ] || [ "$isF" == "Y" ]
    then 
        sudo docker-compose -f ../deployments/development/docker-compose.yml logs -f
    else
        sudo docker-compose -f ../deployments/development/docker-compose.yml logs
    fi

elif [ "$1" == "help" ] || [ "$1" == "-h" ] || [ "$1" == "--help" ]
then
    echo -e " Usage:
  ./dev-docker-compose.sh [COMMAND]
  ./dev-docker-compose.sh -h|--help

 Commands:
   build      Build SafeU dev container images.
   down       Down SafeU dev containers.
   help       Show this help message.
   logs       View output from dev containers.
   up         Up SafeU dev containers.
"

else
    echo -e " Cannot match the command \"$1\", please type \"help\" command for help."
fi
