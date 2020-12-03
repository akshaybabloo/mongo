#!/bin/bash

mkdir -p /data/db
mongod --fork -f /etc/mongod.conf.orig -dbpath /data/db

go test -v -race -coverprofile=coverage.txt -covermode=atomic

bash <(curl -s https://codecov.io/bash)