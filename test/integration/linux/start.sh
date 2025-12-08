#!/usr/bin/env bash

docker image pull ghcr.io/trippsoft/docker-forge-debian13:latest
docker image pull ghcr.io/trippsoft/docker-forge-debian12:latest
docker image pull ghcr.io/trippsoft/docker-forge-fedora42:latest
docker image pull ghcr.io/trippsoft/docker-forge-fedora41:latest
docker image pull ghcr.io/trippsoft/docker-forge-rocky10:latest
docker image pull ghcr.io/trippsoft/docker-forge-rocky9:latest
docker image pull ghcr.io/trippsoft/docker-forge-rocky8:latest
docker image pull ghcr.io/trippsoft/docker-forge-ubuntu2404:latest
docker image pull ghcr.io/trippsoft/docker-forge-ubuntu2204:latest

docker run -d --rm --privileged --name forge-debian13 -p 2201:22 ghcr.io/trippsoft/docker-forge-debian13:latest
docker run -d --rm --privileged --name forge-debian12 -p 2202:22 ghcr.io/trippsoft/docker-forge-debian12:latest
docker run -d --rm --privileged --name forge-fedora42 -p 2211:22 ghcr.io/trippsoft/docker-forge-fedora42:latest
docker run -d --rm --privileged --name forge-fedora41 -p 2212:22 ghcr.io/trippsoft/docker-forge-fedora41:latest
docker run -d --rm --privileged --name forge-rocky10 -p 2221:22 ghcr.io/trippsoft/docker-forge-rocky10:latest
docker run -d --rm --privileged --name forge-rocky9 -p 2222:22 ghcr.io/trippsoft/docker-forge-rocky9:latest
docker run -d --rm --privileged --name forge-rocky8 -p 2223:22 ghcr.io/trippsoft/docker-forge-rocky8:latest
docker run -d --rm --privileged --name forge-ubuntu2404 -p 2231:22 ghcr.io/trippsoft/docker-forge-ubuntu2404:latest
docker run -d --rm --privileged --name forge-ubuntu2204 -p 2232:22 ghcr.io/trippsoft/docker-forge-ubuntu2204:latest
