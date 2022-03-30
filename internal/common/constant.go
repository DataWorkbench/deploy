package common

const (
	DefaultHelmRepoConfigFmt = "%s/.config/helm/repositories.yaml"
	DefaultHelmRepoCacheFmt  = "%s/.cache/helm/repository"

	DefaultOperatorNamespace = "dataomnis-operator"
	DefaultSystemNamespace   = "dataomnis-system"

	LocalHomeFmt         = "%s/dataomnis"
	DataomnisHostPathFmt = "%s/dataomnis/%s"

	DefaultTimeoutSecond = 600

	// flink helm chart name
	FlinkChart = "flink-0.1.6.tgz"
)

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
	DataomnisSystemChart = "dataomnis-0.8.0.tgz"
	DataomnisSystemName  = "dataomnis"
)

const (
	MysqlExternalHostFmt = "%s-pxc-db-haproxy"
	MysqlSecretNameFmt   = "%s-pxc-db"
	HdfsConfigMapFmt     = "%s-common-config"

	RedisClusterPort        = 6379
	RedisClusterAddrFmt     = "%s-%d:%d"
	RedisClusterModeCluster = "cluster"

	InstanceLabelKey           = "app.kubernetes.io/instance"
	MysqlInstanceLabelValueFmt = "%s-pxc-db"
	HdfsInstanceLabelKey       = "dataomnis.io/cluster-name"
	RedisInstanceLabelKey      = "redis.kun/name"
)

const TmpValuesFile = "/tmp/dataomnis-values.yaml"

// Error Fmt
const (
	OnlyNodeNumFmt    = "only %d nodes are required for %s local-pv"
	AtLeaseNodeNumFmt = "at lease %d nodes are required for %s local-pv"
)
