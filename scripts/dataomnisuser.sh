#/bin/bash
mysqlclustername=mysql-cluster
password=pa88w0rd
namespace=dataomnis
kubectl exec -it $mysqlclustername-pxc-db-pxc-0 -c pxc -n $namespace -- mysql -uroot -p$password -e 'CREATE USER "dataomnis"@"%" IDENTIFIED BY "dataomnis"'
kubectl exec -it $mysqlclustername-pxc-db-pxc-0 -c pxc -n $namespace -- mysql -uroot -p$password -e 'CREATE USER "flyway"@"%" IDENTIFIED BY "dataomnis"'
kubectl exec -it $mysqlclustername-pxc-db-pxc-0 -c pxc -n $namespace -- mysql -uroot -p$password -e 'GRANT INSERT,SELECT,UPDATE,DELETE ON dataomnis.* TO "dataomnis"@"%"'
kubectl exec -it $mysqlclustername-pxc-db-pxc-0 -c pxc -n $namespace -- mysql -uroot -p$password -e 'GRANT EXECUTE,INSERT,SELECT,UPDATE,DELETE,ALTER ON dataomnis.* TO "flyway"@"%"'