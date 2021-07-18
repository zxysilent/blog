#!/bin/bash

name="blog"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=prod -ldflags "-s" -o $name main.go


