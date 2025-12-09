#!/usr/bin/env bash

mkdir -p plugins/forge/discover
rm -f plugins/forge/discover/forge-discover_windows_amd64.exe
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./plugins/forge/discover/forge-discover_windows_amd64.exe ../../../cmd/forge-discover/main.go
chmod +x ./plugins/forge/discover/forge-discover_windows_amd64.exe
