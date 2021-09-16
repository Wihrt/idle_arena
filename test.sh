#!/bin/bash

docker-compose stop 
docker-compose rm -f
docker-compose build
docker-compose up -d mongo mongo-express backend bot
if [ $? -gt 0 ]
then
    exit 1
fi

docker system prune -f
sleep 5
docker-compose up -d loader
if [ $? -gt 0 ]
then
    exit 1
fi
docker-compose logs -f backend bot