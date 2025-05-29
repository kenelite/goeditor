#!/usr/bin/env bash

# install dependencies on linux
sudo apt-get install -y golang gcc libgl1-mesa-dev xorg-dev libxkbcommon-dev

# install dependencies on macOS
xcode-select --install

#build for macOS
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o goeditor-macos main.go


# build for macOS arm64
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o goeditor-macos-arm64 main.go


# build for windows amd64
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o goeditor-windows main.go


# build for linux amd64
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o goeditor-linux main.go
