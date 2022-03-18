# Copyright 2020 The Dataomnis Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM flyway/flyway:8.0.0-alpine

USER root

RUN apk add --no-cache mysql-client

RUN mkdir -p /flyway/ddl
COPY ./ddl/* /flyway/ddl/
COPY ./sql/*.sql  /flyway/sql/
