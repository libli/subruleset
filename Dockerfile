# builder
FROM golang:1.20-bullseye as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download -x

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -x -o subruleset .

# runner
FROM debian:bullseye-slim
ENV TZ=Asia/Shanghai
ENV LANG=C.UTF-8

# 安装必要的依赖
RUN apt-get update && \
    apt-get install -y locales ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# 配置中文UTF-8编码
RUN echo "zh_CN.UTF-8 UTF-8" > /etc/locale.gen && \
    locale-gen zh_CN.UTF-8

WORKDIR /app
COPY --from=builder /build/subruleset /app/
CMD ["/app/subruleset"]