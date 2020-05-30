#!/bin/sh -l

mkdir -p /data/db
mongod --fork -f /etc/mongod.conf

go test -v -race -coverprofile=coverage.txt -covermode=atomic
