package installer

import (
	"encoding/json"
	"errors"
	"fmt"
)

const RoleNameNode = "namenode"
var LeastNodeErr = errors.New("at least 3 nodes are required for helm release")


// common config for datanode / journalnode / namenode / zookeeper
type HdfsNodeConfig struct {
	Nodes        []string           `json:"nodes,omitempty" yaml:"nodes,omitempty" validate:"eq=0|min=3"`
	LocalStorage LocalStorageConfig `json:"storage,omitempty" yaml:"storage,omitempty"`
}

func (node HdfsNodeConfig) updateFromHdfsConfig(c HdfsChart, role string) error {
	if node.LocalStorage.Size != "" {
		// update storage.capacity from storage.size
		node.LocalStorage.Capacity = node.LocalStorage.Size
		node.LocalStorage.Size = ""
	}

	if len(node.Nodes) < 3 {
		if len(c.values.Nodes) < 3 {
			return LeastNodeErr
		}

		if role == RoleNameNode {
			node.Nodes = c.values.Nodes
		} else {
			node.Nodes = c.values.Nodes[:2]
		}
	}
	return nil
}

// HdfsValuesConfig for hdfs-cluster
type HdfsValuesConfig struct {
	ValuesConfig `json:",inline" yaml:",inline"`

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
func (h HdfsChart) updateFromConfig(c Config) error {
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
	if err = h.values.Datanode.updateFromHdfsConfig(h, ""); err != nil {
		return err
	}
	if err = h.values.Journalnode.updateFromHdfsConfig(h, ""); err != nil {
		return err
	}
	if err = h.values.Namenode.updateFromHdfsConfig(h, RoleNameNode); err != nil {
		return err
	}
	err = h.values.Zookeeper.updateFromHdfsConfig(h, "")
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
