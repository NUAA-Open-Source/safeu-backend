#!/bin/bash

echo -e "\n Are you sure to remove all the database files? (y/n): \c"
read isSure

if [ "$isSure" == "y" ] || [ "$isSure" == "Y" ]
then
    echo -e "\n Removing the release database data..."
    sudo rm -rf ../data
    echo -e "\n Removing the debug database data..."
    sudo rm -rf ../db-data
else
    echo -e "\n Be careful!"
fi

echo -e "\n Bye~"
