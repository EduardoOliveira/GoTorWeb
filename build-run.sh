#!/bin/bash

GOOS=linux go build main.go
docker build . -t gtw
docker run --rm --name gtw  -v /var/run/docker.sock:/var/run/docker.sock   -e "GTW-NETWORK=gwt" -e "GTW-USERNAME=user" -e "GTW-PASSWORD=changeme" -v `pwd`/config:/config gtw