#!/bin/bash

docker-compose up -d

sleep 5

docker-compose exec mongo_primary /scripts/rs-init.sh

sleep 10

docker-compose down
