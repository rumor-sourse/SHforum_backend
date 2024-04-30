FROM golang:1.19.3

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 下载依赖信息
RUN go mod tidy

# 将我们的代码编译成二进制可执行文件 bubble
RUN go build -o shforum .


# 需要运行的命令
ENTRYPOINT ["./shforum"]
