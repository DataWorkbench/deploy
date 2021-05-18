# Copyright 2020 The DataWorkbench Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM dataworkbench/builder as builder

WORKDIR /go/src/DataWorkbench
COPY . .

RUN mkdir -p /dataworkbench_bin
RUN ./apiserver/scripts/compile.sh /dataworkbench_bin
RUN ./spacemanager/scripts/compile.sh /dataworkbench_bin

RUN find /dataworkbench_bin -type f -exec upx {} \;

FROM alpine:3.12
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --update ca-certificates && update-ca-certificates

RUN mkdir -p /dataworkbench/conf
COPY ./apiserver/config/config.yaml /dataworkbench/conf/apiserver.yaml
COPY ./spacemanager/config/config.yaml /dataworkbench/conf/spacemanager.yaml

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /dataworkbench_bin/* /usr/local/bin/

CMD ["sh"]