package installer

const LocalHomeFmt = "%s/dataomnis"

const (
	// operators:
	// helm chart name
	HdfsOptChart  = "hdfs-operator-0.1.0.tgz"
	RedisOptChart = "redis-cluster-operator-0.1.0.tgz"
	MysqlOptChart = "pxc-operator-1.9.1.tgz"
	// release name
	HdfsOptName  = "hdfs-operator"
	RedisOptName = "redis-cluster-operator"
	MysqlOptName = "mysql-operator"

	// Dependence-Service:
	// helm chart name
	EtcdClusterChart  = "etcd-cluster-1.0.0.tgz"
	HdfsClusterChart  = "hdfs-cluster-0.1.1.tgz"
	RedisClusterChart = "redis-cluster-0.1.0.tgz"
	MysqlClusterChart = "pxc-db-1.9.1.tgz"
	// release name
	EtcdClusterName  = "etcd-cluster"
	HdfsClusterName  = "hdfs-cluster"
	RedisClusterName = "redis-cluster"
	MysqlClusterName = "mysql-cluster"

	// dataomnis
	DataomnisSystemChart = "dataomnis-1.0.0.tgz"
	DataomnisSystemName = "dataomnis"
)

const (
	MysqlExternalHostFmt = "%s-pxc-db-haproxy"
	MysqlSecretNameFmt   = "%s-pxc-db"
	HdfsConfigMapFmt     = "%s-common-config"
	RedisAddressFmt      = "rfs-%s"
)

const (
	DefaultOperatorNamespace = "dataomnis-operator"
	DefaultSystemNamespace   = "dataomnis-system"

	// flink helm chart name
	FlinkChart = "flink-0.1.6.tgz"

	InstanceLabelKey           = "app.kubernetes.io/instance"
	MysqlInstanceLabelValueFmt = "%s-pxc-db"
	HdfsInstanceLabelKey = "qy.dataworkbench.com/cluster-name"
)
