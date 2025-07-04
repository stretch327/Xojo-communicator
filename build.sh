#!/bin/bash

# script to build xojo-cli for mac, windows & linux

echo "Making build directories"
rm -rf build
mkdir -p build/linux/intel
mkdir -p build/linux/arm
mkdir -p build/mac/intel
mkdir -p build/mac/arm
mkdir -p build/windows/intel
mkdir -p build/windows/arm

cd src

echo "Building Linux"
GOOS=linux GOARCH=amd64 go build -o ../build/linux/intel/xojo-cli main.go

echo "Building Linux ARM"
GOOS=linux GOARCH=arm64 go build -o ../build/linux/arm/xojo-cli main.go

echo "Building Mac Intel 64"
GOOS=darwin GOARCH=amd64 go build -o ../build/mac/intel/xojo-cli main.go

echo "Building Mac ARM 64"
GOOS=darwin GOARCH=arm64 go build -o ../build/mac/arm/xojo-cli main.go

echo "Building Windows Intel 64"
GOOS=windows GOARCH=amd64 go build -o ../build/windows/intel/xojo-cli.exe main.go

echo "Building Windows ARM 64"
GOOS=windows GOARCH=arm64 go build -o ../build/windows/arm/xojo-cli.exe main.go
