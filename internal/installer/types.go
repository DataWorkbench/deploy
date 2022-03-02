package installer

import (
	v1 "k8s.io/api/core/v1"
)

// ***************************************************************
// K8s Configuration Zone
// ***************************************************************

type ImageConfig struct {
	// TODO: check pull secrets
	Registry    string   `structs:"registry,omitempty" yaml:"registry"`
	PullSecrets []string `structs:"pullSecrets,omitempty" yaml:"pullSecret"`
	PullPolicy  string   `structs:"pullPolicy,omitempty" yaml:"pullPolicy"`
}

type Storage struct {
	Size int8 `structs:"size,omitempty" yaml:"size"`
}

// k8s workload configurations
type Workload struct {
	ReplicaCount   int8 `structs:"replicaCount,omitempty" yaml:"replicaCount"`
	UpdateStrategy int8 `structs:"updateStrategy,omitempty" yaml:"updateStrategy"`

	ReadinessProbe v1.Probe `structs:"readinessProbe,omitempty" yaml:"readinessProbe"`
	LivenessProbe  v1.Probe `structs:"livenessProbe,omitempty" yaml:"livenessProbe"`

	Resources v1.ResourceRequirements `structs:"resources,omitempty" yaml:"resources"`

	// for local pv, GB
	StorageSize int8 `structs:"storageSize,omitempty" yaml:"storageSize"`
}

// ***************************************************************
// Dataomnis Configuration Zone
// ***************************************************************

type ServiceConfig struct {
	LogLevel         int8 `structs:"logLevel,omitempty" yaml:"logLevel"`
	GrpcLogLevel     int8 `structs:"grpcLogLevel,omitempty" yaml:"grpcLogLevel"`
	GrpcLogVerbosity int8 `structs:"grpcLogVerbosity,omitempty" yaml:"grpcLogVerbosity"`

	ServiceName    bool `structs:"serviceName,omitempty" yaml:"serviceName"`
	MetricsEnabled bool `structs:"metricsEnabled,omitempty" yaml:"metricsEnabled"`

	Workload `structs:",,omitemptyinline" yaml:",inline"`

	Envs []map[string]string `structs:"envs,omitempty" yaml:"envs,flow"`
}

type Webservice struct {
	Enabled bool `yaml:"enabled"`
}

type ApiGlobalConfig struct {
	Enabled           bool     `yaml:"enabled" yaml:"enabled"`
	Regions           []Region `structs:"regions,omitempty" yaml:"regions,flow"`
	ServiceConfig `structs:",,omitemptyinline" yaml:",inline"`
}

type Region struct {
	Hosts string `structs:"hosts,omitempty" yaml:"hosts"`
	EnUs  string `structs:"en_us,omitempty" yaml:"en_us_name"`
	ZhCn  string `structs:"zh_cn,omitempty" yaml:"zh_cn_name"`
}

type ServiceMonitorConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Interval string `yaml:"interval"`
}

type MysqlClientConfig struct {
	LogLevel        int8   `structs:"logLevel,omitempty" yaml:"logLevel"`
	MaxIdleConn     int8   `structs:"maxIdleConn,omitempty" yaml:"maxIdleConn"`
	MaxOpenConn     int8   `structs:"maxOpenConn,omitempty" yaml:"maxOpenConn"`
	ConnMaxLifetime string `structs:"connMaxLifetime,omitempty" yaml:"connMaxLifetime"`
	SlowTshreshold  string `structs:"slowTshreshold,omitempty" yaml:"slowTshreshold"`
}

// ***************************************************************
// Dependence-Service Configuration Zone
// ***************************************************************

type LocalStorage struct {
	Storage `yaml:",inline"`
}

type HdfsNodeConfig struct {
	Nodes   []string          `structs:"nodes,omitempty" yaml:"nodes"`
	Storage map[string]string `structs:"storage,omitempty" yaml:"storage"`
}

// ***************************************************************
// All Configuration From dataomnis-conf.yaml Zone
// ***************************************************************

type Config struct {
	// Local PV home
	LocalPVHome string `structs:"localPVHome" yaml:"localPVHome" validate:"required"`

	Image ImageConfig `structs:"image,omitempty" yaml:"image"`

	// dataomnis version
	DataomnisConf Dataomnis `structs:"dataomnis" yaml:"dataomnis"`

	// dependent service
	EtcdConf  EtcdCluster  `structs:"etcdCluster" yaml:"etcdCluster"`
	HdfsConf  HdfsCluster  `structs:"hdfsCluster" yaml:"hdfsCluster"`
	MysqlConf MysqlCluster `structs:"mysqlCluster" yaml:"mysqlCluster"`
	RedisConf RedisCluster `structs:"redisCluster" yaml:"redisCluster"`
}
