# Copyright 2020 The DataWorkbench Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

# NOTICE: run docker build at up-level current directory
FROM dataworkbench/builder as builder
ARG BIN_IN_BUILDER=/dataworkbench/bin
ARG CONF_IN_BUILDER=/dataworkbench/conf
ARG COMPILE_CMD=./deploy/build/scripts/compile.sh
ARG SERVICES=apiglobal,apiserver,spacemanager,flowmanager,jobmanager,jobdeveloper,jobwatcher,scheduler,sourcemanager,udfmanager,resourcemanager,notifier,observer,enginemanager,logmanager
WORKDIR /go/src/DataWorkbench

COPY . .
# compile service in databench
RUN ${COMPILE_CMD} -s ${SERVICES} -o ${BIN_IN_BUILDER} -c ${CONF_IN_BUILDER}
# compress cmds (do not need to un-compress while run)
RUN find ${BIN_IN_BUILDER} -type f -exec upx {} \;


FROM alpine:3.12
ARG BIN_IN_BUILDER=/dataworkbench/bin
ARG CONF_IN_BUILDER=/dataworkbench/conf
ENV DATABENCH_CONF=/etc/dataworkbench
RUN mkdir -p ${DATABENCH_CONF}

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /opt/bin/* /usr/local/bin/
COPY --from=builder ${BIN_IN_BUILDER}/* /usr/local/bin/
COPY --from=builder ${CONF_IN_BUILDER}/* ${DATABENCH_CONF}/

CMD ["sh"]
