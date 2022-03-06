package installer

import (
	"encoding/json"
	"fmt"
)

const RoleNameNode = "namenode"

// common config for datanode / journalnode / namenode / zookeeper
type HdfsNodeConfig struct {
	// alias for Persistent
	Storage PersistentConfig `json:"storage,omitempty" yaml:"Persistent,omitempty"`
}

func (node HdfsNodeConfig) updateFromHdfsConfig(c HdfsValuesConfig, role string) error {
	if len(node.Storage.LocalPv.Nodes) < 3 {
		if len(c.Nodes) < 3 {
			return LeastNodeErr
		}

		if role == RoleNameNode {
			node.Storage.LocalPv.Nodes = c.Nodes
		} else {
			node.Storage.LocalPv.Nodes = c.Nodes[:2]
		}
	}
	return nil
}

// HdfsValuesConfig for hdfs-cluster
type HdfsValuesConfig struct {
	Image *ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`

	Nodes []string `json:"nodes,omitempty" yaml:"nodes,omitempty" validate:"eq=0|min=3"`

	HdfsHome string `json:"hdfsHome" yaml:"-"`

	Namenode    HdfsNodeConfig `json:"namenode,omitempty" yaml:"namenode" validate:"eq=0|eq=2"`
	Datanode    HdfsNodeConfig `json:"datanode,omitempty" yaml:"datanode" validate:"eq=0|min=3"`
	Journalnode HdfsNodeConfig `json:"journalnode,omitempty" yaml:"journalnode" validate:"eq=0|min=3"`

	Zookeeper HdfsNodeConfig `json:"zookeeper,omitempty" yaml:"zookeeper" validate:"eq=0|min=3"`
}

// TODO: validate the yaml and nodes == 2 of namenode
func (v HdfsValuesConfig) validate() error {
	return nil
}

// ********************************************
// HdfsChart for hdfs-cluster, implement Chart
type HdfsChart struct {
	ChartMeta

	values HdfsValuesConfig
}

// update each field value from global Config if that is ZERO
func (h HdfsChart) updateConfig(c Config) error {
	if c.Image != nil {
		if h.values.Image == nil {
			h.values.Image = &ImageConfig{}
		}
		h.values.Image.updateFromConfig(c.Image)
	}

	// TODO: check if localPv exist and start with c.LocalPVHome
	h.values.HdfsHome = fmt.Sprintf(LocalHomeFmt, c.LocalPVHome)

	// update datanode / namenode / journalnode / zookeeper conf
	var err error
	if err = h.values.Datanode.updateFromHdfsConfig(h.values, ""); err != nil {
		return err
	}
	if err = h.values.Journalnode.updateFromHdfsConfig(h.values, ""); err != nil {
		return err
	}
	if err = h.values.Namenode.updateFromHdfsConfig(h.values, RoleNameNode); err != nil {
		return err
	}
	err = h.values.Zookeeper.updateFromHdfsConfig(h.values, "")
	return err
}

func (h HdfsChart) parseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	bytes, err := json.Marshal(h.values)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}
