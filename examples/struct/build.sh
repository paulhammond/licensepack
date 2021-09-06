#!/bin/sh

set -e

# change to this directory
cd "${0%/*}"

# generate
go generate .

# build
go build -o hello .

echo "done"