#!/bin/sh
set -e
cd /app
echo "Downloading Go modules..."
go mod download
exec air -c .air.toml
