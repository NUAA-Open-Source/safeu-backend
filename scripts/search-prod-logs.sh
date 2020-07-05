#!/bin/bash
# Author: TripleZ<me@triplez.cn>

echo -e "
 This script is for searching the safeu-backend application log with
 your specific keyword.
 
 Usage:
   ./search-prod-logs.sh [OPTIONS]             Search logs
  
 Example:
   ./search-prod-logs.sh --tail 1000           Show the last 1000 lines of log
 "

echo -e "\n Input your query keyword: \c"
read kw
echo ""

sudo docker-compose -f ../deployments/prod-safeu/docker-compose.yml logs -t "$@" | grep "$kw" | more

echo -e "\n Bye~"