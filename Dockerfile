# Copyright 2020 The DataWorkbench Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM dataworkbench/builder as builder

WORKDIR /go/src/DataWorkbench
COPY . .

ARG DW_BIN_IN_BUILDER=/dataworkbench/bin
ARG DW_CONF_IN_BUILDER=/dataworkbench/conf
ARG COMPILE_CMD=./deploy/build/scripts/compile.sh
RUN mkdir -p ${DW_BIN_IN_BUILDER}
RUN mkdir -p ${DW_CONF_IN_BUILDER}

# install grpc_health_probe for status probe of dataworkbench service on k8s
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.4 && \
    wget -qO ${DW_BIN_IN_BUILDER}/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x ${DW_BIN_IN_BUILDER}/grpc_health_probe

RUN ${COMPILE_CMD} -o ${DW_BIN_IN_BUILDER} -c ${DW_CONF_IN_BUILDER}
# compress cmds (do not need to un-compress while run)
RUN find ${DW_BIN_IN_BUILDER} -type f -exec upx {} \;


FROM alpine:3.12
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --update ca-certificates && update-ca-certificates

ENV DATABENCH_CONF=/etc/dataworkbench
RUN mkdir -p ${DATABENCH_CONF}

ARG DW_BIN_IN_BUILDER=/dataworkbench/bin
ARG DW_CONF_IN_BUILDER=/dataworkbench/conf

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder ${DW_BIN_IN_BUILDER}/* /usr/local/bin/
COPY --from=builder ${DW_CONF_IN_BUILDER}/* ${DATABENCH_CONF}/

CMD ["sh"]
