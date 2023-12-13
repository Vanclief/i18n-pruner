#!/bin/bash

# We assume your Go file is main.go.
# If its not, replace `main.go` with your actual Go filename.

APP_NAME="i18n-pruner"
GO_FILE="main.go"

# Mac ARM64
echo "Compiling for MacOS (ARM)..."
GOOS=darwin GOARCH=arm64 go build -o ./bin/$APP_NAME-mac-arm64 $GO_FILE

# Mac x86_64
echo "Compiling for MacOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o ./bin/$APP_NAME-mac-intel $GO_FILE

# Linux
echo "Compiling for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o ./bin/$APP_NAME-linux $GO_FILE

# Windows
echo "Compiling for Windows (64 bit)..."
GOOS=windows GOARCH=amd64 go build -o ./bin/$APP_NAME-windows.exe $GO_FILE

echo "---- All Done! ----"
