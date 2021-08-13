#!/bin/bash

docker-compose stop
docker-compose rm -f
docker-compose build
docker-compose up -d mongo mongo-express backend bot
sleep 5
docker-compose up -d loader
docker-compose logs -f backend bot