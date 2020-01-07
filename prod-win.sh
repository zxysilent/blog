#!/bin/bash

swag init
name="blog"
go build -tags=prod -o $name.exe main.go


