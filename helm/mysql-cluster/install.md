1 进入mysql-operator目录，执行命令   helm install mysql-operator .  -n dataomnis
2 在三个work节点创建localpv目录 mkdir -p /data/dataomnis/mysql-cluster  mysql-cluster 跟步骤三helm名称保持一致
3 进入 mysql-cluster目录，执行命令    helm install mysql-cluster . -n dataomnis