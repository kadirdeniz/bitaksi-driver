#!/bin/bash

docker run --name driver -p 27017:27017 -e MONGO_INITDB_DATABASE=bitaksi -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin -d mongo