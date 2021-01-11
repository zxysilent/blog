#!/usr/bin/env bash

rm -rf logs
mkdir logs

cd vue
npm install
cd ..

go get -u github.com/swaggo/swag/cmd/swag
swag init

go run main.go

#导入sql数据

