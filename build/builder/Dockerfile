# Copyright 2020 The DataWorkbench Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM golang:1.15.5-alpine3.12 as builder

ARG WORKSPACE=/opt
WORKDIR ${WORKSPACE}

# compress cmds (do not need to un-compress while run)
RUN wget https://github.com/upx/upx/releases/download/v3.95/upx-3.95-amd64_linux.tar.xz && \
    tar -Jxf upx-3.95-amd64_linux.tar.xz && \
    mv upx-*/upx /usr/local/bin/

# install grpc_health_probe: v0.4.4 for status probe of dataworkbench service on k8s
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.4 && \
    wget -qO grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
#RUN wget -qO grpc_health_probe https://api.anybox.qingcloud.com/v1/downloads/VrZ76Pr4qxKLv3WjTJoqq00dJxwIg30Ie0XVh5fpQOnI2wGjCzd3OLuSyZwXE0d8PLkzCokx8YT4vmzBpNJ49HlUIggrPopCL1uTwHPikUkJTVt26hyRIUia1ZxcNZOg/grpc_health_probe-linux-amd64 && \
    chmod +x grpc_health_probe && \
    upx grpc_health_probe && \
    mv grpc_health_probe /usr/local/bin/


FROM golang:1.15.5-alpine3.12
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache bash git make openssl build-base

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w GOCACHE=/go/cache

COPY --from=builder /usr/local/bin/* /usr/local/bin/
