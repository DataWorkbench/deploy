# Copyright 2020 The Dataomnis Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.
FROM dockerhub.qingcloud.com/dataomnis/builder as builder

ARG DATAOMNIS_CONF=/etc/dataomnis
RUN mkdir -p ${DATAOMNIS_CONF}

ARG SERVICE
COPY ./deploy/tmp/conf/${SERVICE}.yaml ${DATAOMNIS_CONF}/
COPY ./deploy/tmp/bin/${SERVICE} /usr/local/bin/
RUN upx /usr/local/bin/${SERVICE}


###############################################################
FROM alpine:3.12

RUN apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ENV DATAOMNIS_CONF=/etc/dataomnis
RUN mkdir -p ${DATAOMNIS_CONF}

ARG SERVICE
COPY --from=builder /usr/local/bin/grpc_health_probe /usr/local/bin/
COPY --from=builder ${DATAOMNIS_CONF}/${SERVICE}.yaml ${DATAOMNIS_CONF}/
COPY --from=builder /usr/local/bin/${SERVICE} /usr/local/bin/

CMD ["sh"]
