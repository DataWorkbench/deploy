#!/bin/sh

[ -n "$MYSQL_HOST" ] && OPT="-h$(echo "$MYSQL_HOST" | tr -d '\n')"
[ -n "$MYSQL_USER" ] && OPT="${OPT} -u$(echo "$MYSQL_USER" | tr -d '\n')"
[ -n "$MYSQL_PASSWORD" ] && OPT="${OPT} -p$(echo "$MYSQL_PASSWORD" | tr -d '\n')"

cd /flyway/sql/ddl || exit

for F in $(ls *.sql)
do
    echo "Start process $F"
    mysql "$@" "$OPT" < "$F"
    if [ $? -ne 0 ]; then
        echo "Process $F failed"
                return 1
        else
                echo "Process $F successful"
        fi
done
