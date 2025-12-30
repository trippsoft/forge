#!/usr/bin/env bash

mkdir -p plugins/forge/discover
mkdir -p plugins/forge/core

rm -f plugins/forge/discover/forge-discover_windows_amd64.exe
rm -f plugins/forge/discover/forge-discover_windows_arm64.exe

rm -f plugins/forge/core/forge-core_linux_amd64
rm -f plugins/forge/core/forge-core_linux_arm64
rm -f plugins/forge/core/forge-core_windows_amd64.exe
rm -f plugins/forge/core/forge-core_windows_arm64.exe

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./plugins/forge/discover/forge-discover_windows_amd64.exe ../../../cmd/forge-discover/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ./plugins/forge/discover/forge-discover_windows_arm64.exe ../../../cmd/forge-discover/main.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./plugins/forge/core/forge-core_linux_amd64 ../../../cmd/forge-core/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./plugins/forge/core/forge-core_linux_arm64 ../../../cmd/forge-core/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./plugins/forge/core/forge-core_windows_amd64.exe ../../../cmd/forge-core/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ./plugins/forge/core/forge-core_windows_arm64.exe ../../../cmd/forge-core/main.go

chmod +x ./plugins/forge/discover/forge-discover_windows_amd64.exe
chmod +x ./plugins/forge/discover/forge-discover_windows_arm64.exe

chmod +x ./plugins/forge/core/forge-core_linux_amd64
chmod +x ./plugins/forge/core/forge-core_linux_arm64
chmod +x ./plugins/forge/core/forge-core_windows_amd64.exe
chmod +x ./plugins/forge/core/forge-core_windows_arm64.exe
