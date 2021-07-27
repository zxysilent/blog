#!/bin/bash

# 跟随系统build
# 任何语句的执行结果不是true则应该退出
set -e
apphash=`git log --pretty=format:'%h' -1`
if [ $? -ne 0 ]; then
    apphash="zxysilent"
fi
appname="blog"
echo "git log $apphash"
# ldflags="-X 'main.Var=$apphashn'"
# -w    去掉调试信息
# -s    去掉符号表
flags=""
# -a    强制重新构建
# -n	打印编译时会用到的所有命令，但不真正执行
# -x	打印编译时会用到的所有命令
# -race	开启竞态检测

buildProd(){
    go build -tags=prod -ldflags "$flags" main.go
}
buildDev(){
    go build -ldflags "$flags" main.go
}
usage() {
  echo "Usage: $0 [-p <prod>] [-d <dev>]" 1>&2;
}
while getopts "pd" o; do
  case "${o}" in
    p)
        echo "build prod"
        buildProd
        ;;
    d)
        echo "build dev"
        buildDev
        ;;
    *)
        usage
        exit 0
        ;;
  esac
done
if [ "$#" -eq 0 ] ; then
    usage
    echo "build dev"
    buildDev
    exit 0
fi