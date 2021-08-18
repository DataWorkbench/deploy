# Copyright 2020 The DataWorkbench Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM dataworkbench/builder as builder

WORKDIR /go/src/DataWorkbench
COPY . .

ARG DW_BIN=/dataworkbench/bin
ARG DW_CONF=/dataworkbench/conf
ARG COMPILE_CMD=./deploy/build/scripts/compile.sh
RUN mkdir -p ${DW_BIN}
RUN mkdir -p ${DW_CONF}

# install grpc_health_probe for status probe of dataworkbench service on k8s
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.4 && \
    wget -qO ${DW_BIN}/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x ${DW_BIN}/grpc_health_probe

RUN ${COMPILE_CMD} -s apiproxy -o ${DW_BIN}/
COPY ./apiproxy/config/config.yaml ${DW_CONF}/apiproxy.yaml

RUN ${COMPILE_CMD} -s apiserver -o ${DW_BIN}/
COPY ./apiserver/config/config.yaml ${DW_CONF}/apiserver.yaml

RUN ${COMPILE_CMD} -s spacemanager -o ${DW_BIN}/
COPY ./spacemanager/config/config.yaml ${DW_CONF}/spacemanager.yaml

RUN ${COMPILE_CMD} -s flowmanager -o ${DW_BIN}/
COPY ./flowmanager/config/config.yaml ${DW_CONF}/flowmanager.yaml

RUN ${COMPILE_CMD} -s jobdeveloper -o ${DW_BIN}/
COPY ./jobdeveloper/config/config.yaml ${DW_CONF}/jobdeveloper.yaml

RUN ${COMPILE_CMD} -s jobmanager -o ${DW_BIN}/
COPY ./jobmanager/config/config.yaml ${DW_CONF}/jobmanager.yaml

RUN ${COMPILE_CMD} -s jobwatcher -o ${DW_BIN}/
COPY ./jobwatcher/config/config.yaml ${DW_CONF}/jobwatcher.yaml

RUN ${COMPILE_CMD} -s notifier -o ${DW_BIN}/
COPY ./notifier/config/config.yaml ${DW_CONF}/notifier.yaml

RUN ${COMPILE_CMD} -s scheduler -o ${DW_BIN}/
COPY ./scheduler/config/config.yaml ${DW_CONF}/scheduler.yaml

RUN ${COMPILE_CMD} -s sourcemanager -o ${DW_BIN}/
COPY ./sourcemanager/config/config.yaml ${DW_CONF}/sourcemanager.yaml

RUN ${COMPILE_CMD} -s udfmanager -o ${DW_BIN}/
COPY ./udfmanager/config/config.yaml ${DW_CONF}/udfmanager.yaml

RUN ${COMPILE_CMD} -s zeppelinscale -o ${DW_BIN}/
COPY ./zeppelinscale/config/config.yaml ${DW_CONF}/zeppelinscale.yaml

RUN ${COMPILE_CMD} -s filemanager -o ${DW_BIN}/
COPY ./filemanager/config/config.yaml ${DW_CONF}/filemanager.yaml

# compress cmds (do not need to un-compress while run)
RUN find ${DW_BIN} -type f -exec upx {} \;


FROM alpine:3.12
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --update ca-certificates && update-ca-certificates

ENV DATAWORKBENCH_CONF=/etc/dataworkbench
RUN mkdir -p ${DATAWORKBENCH_CONF}

ARG DW_BIN=/dataworkbench/bin
ARG DW_CONF=/dataworkbench/conf

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder ${DW_BIN}/* /usr/local/bin/
COPY --from=builder ${DW_CONF}/* ${DATAWORKBENCH_CONF}/

CMD ["sh"]