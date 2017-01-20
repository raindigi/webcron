#!/bin/sh

tarfile="webcron-$1.tar.gz"

echo "Start compressing $tarfile..."

export GOARCH=amd64
export GOOS=linux

bee pack

mv webcron.tar.gz $tarfile
