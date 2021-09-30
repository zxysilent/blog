# 构建：使用golang:1.16 latest版本
FROM golang:1.16 as builder
MAINTAINER https://github.com/zxysilent
# 容器环境变量添加，会覆盖默认的变量值
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
# 设置工作区
# GOPATH "/go"
WORKDIR /go/cache
# 复制文件下载依赖
ADD go.mod .
ADD go.sum .
RUN go mod download

# 把全部文件添加到/ao/app目录
# 设置工作区
# GOPATH "/go"
WORKDIR /go/build

# ADD src dest(loacl->docker)
ADD . .

# 编译：把main.go编译成可执行的二进制文件，命名为app
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o app main.go

# 运行：使用scratch作为基础镜像
FROM scratch as runner

# 在build阶段复制时区到
# COPY src dest
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 在build阶段复制文件
# exe
COPY --from=builder /go/build/app /
# vue
COPY --from=builder /go/build/dist/ /dist/
# views
COPY --from=builder /go/build/views/ /views/
# static
COPY --from=builder /go/build/static/ /static/
# conf
COPY --from=builder /go/build/conf/conf.toml /conf/
EXPOSE 8085
# 启动服务
CMD ["/app"]