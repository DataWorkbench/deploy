package installer

import (
	"encoding/json"
	v1 "k8s.io/api/core/v1"
)

type EtcdCluster struct {
	Nodes      []string    `structs:"nodes,omitempty" yaml:"nodes" validate:"min=3"`
	Image      *ImageConfig `structs:"image,omitempty" yaml:"image"`
	Persistent Storage     `structs:"persistent" yaml:"storage"`
}

type HdfsCluster struct {
	Image    ImageConfig `structs:"image,omitempty" yaml:"image"`
	Nodes    []string    `structs:"nodes,omitempty" yaml:"nodes"`
	HdfsHome string      `structs:"hdfsHome,omitempty" yaml:"hdfsHome"`

	Namenode    HdfsNodeConfig `structs:"namenode,omitempty" yaml:"namenode"`
	Datanode    HdfsNodeConfig `structs:"datanode,omitempty" yaml:"datanode"`
	Journalnode HdfsNodeConfig `structs:"journalnode,omitempty" yaml:"journalnode"`

	Zookeeper HdfsNodeConfig `structs:"zookeeper,omitempty" yaml:"zookeeper"`
}

type MysqlCluster struct {
	Nodes   []string          `structs:"nodes,omitempty" yaml:"nodes"`
	Image   ImageConfig       `structs:"image,omitempty" yaml:"image"`
	Storage map[string]string `structs:"storage,omitempty" yaml:"storage"`

	// for pxc node
	Resources v1.ResourceRequirements `structs:"resources,omitempty" yaml:"resources"`
	// TODO: support backup configuration for s3
}

type RedisCluster struct {
	Nodes   []string          `structs:"nodes,omitempty" yaml:"nodes"`
	Image   ImageConfig       `structs:"image,omitempty" yaml:"image"`
	Storage map[string]string `structs:"storage,omitempty" yaml:"storage"`
	// for redis node
	Resources v1.ResourceRequirements `structs:"resources,omitempty" yaml:"resources"`
}

func (e EtcdCluster) parseValues(c Config) Values {
	if e.Image != nil {
		e.Image.overwritten(c.Image)
	}

	json.Marshal(e)

}
