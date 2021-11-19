#!/bin/sh

cd /flyway/sql/ddl || exit

[ -n "$MYSQL_PASSWORD" ] && OPT="-p$(echo "$MYSQL_PASSWORD" | tr -d '\n')"

for F in $(ls *.sql)
do
    echo "Start process $F"
    mysql "-h$MYSQL_HOST" "-P$MYSQL_PORT" "-u$MYSQL_USER" "$OPT" "$@" < "$F"
    if [ $? -ne 0 ]; then
        echo "Process $F failed"
                return 1
        else
                echo "Process $F successful"
        fi
done
