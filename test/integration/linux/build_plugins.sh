#!/usr/bin/env bash

mkdir -p plugins/forge/discover
rm -f plugins/forge/discover/forge-discover_linux_amd64
rm -f plugins/forge/discover/forge-discover_linux_arm64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./plugins/forge/discover/forge-discover_linux_amd64 ../../../cmd/forge-discover/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./plugins/forge/discover/forge-discover_linux_arm64 ../../../cmd/forge-discover/main.go
chmod +x ./plugins/forge/discover/forge-discover_linux_amd64
chmod +x ./plugins/forge/discover/forge-discover_linux_arm64
