package installer

import (
	v1 "k8s.io/api/core/v1"
)

// ***************************************************************
// K8s Configuration Zone
// ***************************************************************

type ValuesConfig struct {
	Image          *ImageConfig `json:"image,inline" yaml:"image,inline"`
	WorkloadConfig `json:",inline" yaml:",inline"`
}

type ImageConfig struct {
	// TODO: check pull secrets
	Registry    string   `json:"registry,omitempty" yaml:"registry"`
	PullSecrets []string `json:"pullSecrets,omitempty" yaml:"pullSecret"`
	PullPolicy  string   `json:"pullPolicy,omitempty" yaml:"pullPolicy"`
}

// update echo field-value from Config
func (i ImageConfig) updateFromConfig(source *ImageConfig) {
	if source == nil {
		return
	}
	if source.Registry != "" {
		i.Registry = source.Registry
	}
	if len(source.PullSecrets) > 0 {
		i.PullSecrets = source.PullSecrets
	}
	if source.PullPolicy != "" {
		i.PullPolicy = source.PullPolicy
	}
}

type StorageConfig struct {
	Size string `json:"size,omitempty" yaml:"size,omitempty"`
}

// k8s workload configurations
type WorkloadConfig struct {
	ReplicaCount   int8 `json:"replicaCount,omitempty" yaml:"replicaCount"`
	UpdateStrategy int8 `json:"updateStrategy,omitempty" yaml:"updateStrategy"`

	ReadinessProbe v1.Probe `json:"readinessProbe,omitempty" yaml:"readinessProbe"`
	LivenessProbe  v1.Probe `json:"livenessProbe,omitempty" yaml:"livenessProbe"`

	Resources v1.ResourceRequirements `json:"resources,omitempty" yaml:"resources"`

	// for local pv, eg: 10Gi
	Storage StorageConfig `json:"storage,omitempty" yaml:"storage"`
}

// ***************************************************************
// Dependence-Service Configuration Zone
// ***************************************************************

type LocalStorageConfig struct {
	StorageConfig `json:",inline,omitempty" yaml:",inline,omitempty"`

	// for etcd-cluster
	LocalHome string `json:"localHome,omitempty" yaml:"-"`

	// for hdfs
	Capacity string `json:"capacity,omitempty" yaml:"-"`
}

// ***************************************************************
// All Configuration From dataomnis-conf.yaml Zone
// ***************************************************************

// TODO: save the dataomnis-conf.yaml to k8s as configmap for backup
type Config struct {
	// kube nodes from k8s apiserver
	Nodes []string `yaml:"-"`

	// Local PV home
	LocalPVHome string `yaml:"localPVHome" validate:"required"`

	Image *ImageConfig `yaml:"image"`

	// dataomnis version
	Dataomnis DataomnisChart `yaml:"dataomnis"`

	// dependent service
	Etcd  EtcdValuesConfig  `yaml:"etcdCluster"`
	Hdfs  HdfsValuesConfig  `yaml:"hdfsCluster"`
	Mysql MysqlChart `yaml:"mysqlCluster"`
	Redis RedisChart `yaml:"redisCluster"`
}
