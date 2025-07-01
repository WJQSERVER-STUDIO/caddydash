FROM alpine:latest AS BUILD

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

# 更新并安装所需的软件包
RUN apk update && apk add --no-cache sed wget curl tar

# 创建文件夹结构
RUN mkdir -p /data/caddy/config
RUN mkdir -p /data/caddy/config.d
RUN mkdir -p /usr/local/bin

# 下载并解压caddy
RUN VERSION=$(curl -s https://raw.githubusercontent.com/WJQSERVER/caddy/main/TEST-VERSION) && \
    echo "VERSION: ${VERSION}" && \
    wget -O /data/caddy/caddy.tar.gz https://github.com/WJQSERVER/caddy/releases/download/${VERSION}/caddy-${TARGETOS}-${TARGETARCH}-pages.tar.gz
RUN tar -xzvf /data/caddy/caddy.tar.gz -C /data/caddy 
RUN rm /data/caddy/caddy.tar.gz 
RUN chmod +x /data/caddy/caddy

RUN DASHVERSION=$(curl -s https://raw.githubusercontent.com/WJQSERVER-STUDIO/caddydash/main/VERSION) && \
    echo "VERSION: ${DASHVERSION}" && \
    wget -O /data/caddy/caddydash.tar.gz https://github.com/WJQSERVER-STUDIO/caddydash/releases/download/${DASHVERSION}/caddydash-${TARGETOS}-${TARGETARCH}.tar.gz
RUN tar -xzvf /data/caddy/caddydash.tar.gz -C /data/caddy 
RUN rm /data/caddy/caddydash.tar.gz 
RUN chmod +x /data/caddy/caddydash

# 构建镜像
FROM alpine:latest

# 创建必要的目录
RUN mkdir -p /data/caddy/config 
RUN mkdir -p /data/caddy/config.d
RUN mkdir -p /data/caddy/log
RUN mkdir -p /data/caddy/db

COPY frontend /data/caddy/
COPY gtmpl /data/caddy/
COPY tmpl /data/caddy/

# 复制构建阶段的文件
COPY --from=BUILD /data/caddy /data/caddy
COPY --from=BUILD /usr/local/bin/init.sh /usr/local/bin/init.sh

# 配置权限
RUN chmod +x /data/caddy/caddy 
RUN chmod +x /data/caddy/caddydash

# 启动命令
WORKDIR /data/caddy/
ENTRYPOINT ["/data/caddy/caddydash", "-c", "/data/caddy/config/config.toml"]