# ---- 构建阶段 ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

# 使用国内代理加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/aidev .

# ---- 运行阶段 ----
FROM alpine:3.23

RUN apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/aidev .

# MySQL 连接配置 (运行时可通过 docker run -e 覆盖)
ENV MYSQL_USER=aivideotest
ENV MYSQL_PASS=aivideotest@1599
ENV MYSQL_HOST=mysql-3c212056e2fb-public.rds.volces.com
ENV MYSQL_PORT=3306
ENV MYSQL_DB=apiproxy

EXPOSE 8080

CMD ["./aidev"]
