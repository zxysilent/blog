#!/bin/bash
# 后端开发模式：热加载
# fish 监听文件改动自动重编译并重启，改 .go 不用手动重新 go run
# https://github.com/zxysilent/fish
#
# 前端（mgmt）另开一个终端执行 dev-mgmt.sh（pnpm run dev），
# vite 会把 /api、/static 代理到本进程的 8085 端口。
set -e
cd "$(dirname "$0")" || exit 1

# -all：除 .go 外也监听 .html，改 views/ 下的模板同样会触发重启
exec fish run -all
