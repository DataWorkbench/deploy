# Copyright 2020 The DataWorkbench Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM golang:1.15.5-alpine3.12 as builder

ARG WORKSPACE=/opt

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache git curl openssl

WORKDIR ${WORKSPACE}

# the protoc can't run on alpine,
# we only need the protobuf's stdarnd library in the `/protoc/include`.
RUN mkdir -p ./protoc && \
    wget https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-linux-x86_64.zip && \
	unzip -d ./protoc protoc-3.14.0-linux-x86_64.zip

RUN wget https://github.com/upx/upx/releases/download/v3.95/upx-3.95-amd64_linux.tar.xz && \
    tar -Jxf upx-3.95-amd64_linux.tar.xz


FROM golang:1.15.5-alpine3.12
ARG WORKSPACE=/opt
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache git protobuf make curl openssl jq rsync bash

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w GOCACHE=/go/cache

COPY --from=builder ${WORKSPACE}/protoc/include /usr/local/include
COPY --from=builder ${WORKSPACE}/upx-*/upx /usr/bin/
