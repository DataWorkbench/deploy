# Copyright 2020 The Dataomnis Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM flyway/flyway:8.0.0-alpine

USER root

RUN apk add --no-cache mysql-client

RUN mkdir -p /flyway/ddl
COPY ./deploy/build/db/ddl_init.sh /flyway/ddl/
COPY ./deploy/build/db/ddl/*.sql /flyway/ddl/

COPY ./spacemanager/scripts/sql/schema/V0.0.0__space_manager_init.sql /flyway/sql/V0.0.0.0__space_manager_init.sql
COPY ./sourcemanager/scripts/sql/schema/V0.0.0__sourcemanager.sql /flyway/sql/V0.0.0.1__sourcemanager.sql
COPY ./jobmanager/scripts/sql/schema/job_manager.sql /flyway/sql/V0.0.0.2__jobmanager.sql
COPY ./udfmanager/scripts/sql/schema/V0.0.0__udfmanager.sql /flyway/sql/V0.0.0.3__udfmanager.sql
COPY ./resourcemanager/scripts/sql/schema/V0.0.0__resource.sql  /flyway/sql/V0.0.0.4__resource.sql
COPY ./notifier/scripts/sql/schema/V0.0.0__notifier.sql  /flyway/sql/V0.0.0.5__notifier.sql
COPY ./account/scripts/sql/schema/V0.0.0__account.sql  /flyway/sql/V0.0.0.6__account.sql
COPY ./enginemanager/scripts/sql/schema/V0.0.0__engine_manager_init.sql  /flyway/sql/V0.0.0.7__engine_manager_init.sql
