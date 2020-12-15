# Copyright 2020 The DataWorkbench Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM dhoer/flyway:5.1.4-mysql-8.0.11-alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache mysql-client

COPY ./deploy/build/db/ddl_init.sh /flyway/sql/ddl/
COPY ./deploy/build/db/ddl/*.sql /flyway/sql/ddl/
COPY ./spacemanager/scripts/sql/schema/*.sql /flyway/sql/
