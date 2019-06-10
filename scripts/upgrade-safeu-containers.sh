#!/bin/bash
# Author:   TripleZ<me@triplez.cn>
# Date:     2019-04-30

helpMsg() {
    echo -e " Usage:
  ./upgrade-safeu-containers.sh [ENV] [BRANCH NAME]
  ./upgrade-safeu-containers.sh -h|--help
  
 Environments: 
   dev          Development environment
   staging      Staging deployment environment
   prod         Production deployment environment
 
 Branchs:
   master       Stable version
   dev          Development version
   [custom]     Custom Git branch
   "
}

if [ "$2" != "" ]
then
    # update code
    git stash  # stash the current modification
    git checkout $2
    git pull origin $2
else
    echo -e "\n [ERROR] Unrecognized branch name!\n"
    helpMsg
    exit 1
fi

if [ "$1" == "prod" ]
then
    sudo ./prod-docker-compose.sh build
    sudo ./prod-docker-compose.sh down
    sudo ./prod-docker-compose.sh up

elif [ "$1" == "staging" ]
then
    sudo ./staging-docker-compose.sh build
    sudo ./staging-docker-compose.sh down
    sudo ./staging-docker-compose.sh up

elif [ "$1" == "dev" ]
then
    sudo ./dev-docker-compose.sh build
    sudo ./dev-docker-compose.sh down
    sudo ./dev-docker-compose.sh up

elif [ "$1" == "-h" ] || [ "$1" == "--help" ]
then
    helpMsg
else
    echo -e " [ERROR] Unrecognized environment name!\n"
    helpMsg
    exit 1
fi
