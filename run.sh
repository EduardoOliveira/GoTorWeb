#!/bin/bash

GOOS=linux go build main.go
docker build . -t gtw
docker run --rm --name gtw  -v /var/run/docker.sock:/var/run/docker.sock -v `pwd`/config:/config gtw