#!/bin/bash
CGO_ENABLED=0 go build -o giaginfra main.go
rm -f linux_build.tgz
tar -czvf linux_build.tgz build conf diaginfra