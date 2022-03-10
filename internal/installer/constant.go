package installer

const LocalHomeFmt = "%s/dataomnis"

const (
	// operators:
	// helm chart name
	HdfsOptChart  = "hdfs-operator-0.1.0.tgz"
	RedisOptChart = "redis-operator-3.1.0.tgz"
	MysqlOptChart = "pxc-operator-1.9.1.tgz"
	// release name
	HdfsOptName  = "hdfs-operator"
	RedisOptName = "redis-operator"
	MysqlOptName = "mysql-operator"

	// Dependence-Service:
	// helm chart name
	EtcdClusterChart  = "etcd-cluster-1.0.0.tgz"
	HdfsClusterChart  = "hdfs-cluster-0.1.1.tgz"
	RedisClusterChart = "redis-operator-3.1.0.tgz"
	MysqlClusterChart = "pxc-operator-1.9.1.tgz"
	// release name
	EtcdClusterName  = "etcd-cluster"
	HdfsClusterName  = "hdfs-cluster"
	RedisClusterName = "redis-cluster"
	MysqlClusterName = "mysql-cluster"
)

const (
	MysqlExternalHostFmt = "%s-pxc-db-haproxy"
	MysqlSecretNameFmt   = "%s-pxc-db"
	HdfsConfigMapFmt     = "%s-common-config"
	RedisAddressFmt      = "rfs-%s"
)
