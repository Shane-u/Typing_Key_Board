FROM golang:1.23.0 as build

# 设置工作目录
WORKDIR /app

# 设置Go模块代理
ENV GOPROXY=https://goproxy.cn,direct

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . ./

# 构建应用，由于在Alpine Linux环境下编译的Go二进制文件无法在容器中执行，默认情况下，Go在Alpine环境中编译的二进制文件是动态链接的，依赖于特定的C库，所以在构建时添加静态链接标志。
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用更小的镜像运行应用
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates

# 创建数据目录
RUN mkdir -p /app/data

WORKDIR /app

# 从构建阶段复制编译好的应用
COPY --from=build /app/main /app/

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"] 