#!/bin/bash
# 管理端（mgmt）开发模式：vite 热更新
# 依赖后端已在 8085 运行（根目录执行 dev-srv.sh，即 fish run）
#
# 统一用 pnpm：mgmt/.npmrc 里的 node-linker / verify-deps-before-run 是 pnpm 专用配置，
# 用 npm 跑会刷一堆 "Unknown project config" 警告（还会按 npm 的规则解析 node_modules）。
set -e
cd "$(dirname "$0")/mgmt" || exit 1

# 首次运行先安装依赖（node_modules 已存在则跳过）
if [ ! -d node_modules ]; then
  pnpm install
fi

exec pnpm run dev
