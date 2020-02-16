#!/bin/bash

name="blog"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=prod -o $name main.go


