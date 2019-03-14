#!/bin/bash

echo -e "\n Are you sure to remove all the database files? (y/n, default: n): \c"
read isSure

if [ "$isSure" == "y" ] || [ "$isSure" == "Y" ]
then
    echo -e "\n Removing the development database data..."
    sudo rm -rf ../data-dev

    echo -e "\n Are you SURE to remove the production database data???
 ALL DATA WILL BE LOST!!!\n (y/n, default: n): \c"
    read isProdSure
    if [ "$isProdSure" == "y" ] || [ "$isProdSure" == "Y" ]
    then
        echo -e "\n Removing the production database data..."
        sudo rm -rf ../data
    fi
else
    echo -e "\n Be careful!"
fi

echo -e "\n Bye~\n"
