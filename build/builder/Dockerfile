# Copyright 2020 The DataWorkbench Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM golang:1.15.5-alpine3.12 as builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache git curl openssl

# the protoc can't run on alpine,
# we only need the protobuf's stdarnd library in the `/protoc/include`.
RUN mkdir -p /protoc && cd /protoc \
	&& wget https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-linux-x86_64.zip \
	&& unzip protoc-3.14.0-linux-x86_64.zip


FROM golang:1.15.5-alpine3.12

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache git protobuf make curl openssl jq rsync bash

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w GOCACHE=/go/cache

COPY --from=builder /protoc/include /usr/local/include
