package installer

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

const RoleNameNode = "namenode"

// common config for datanode / journalnode / namenode / zookeeper
type HdfsNodeConfig struct {
	// alias for Persistent, updated from configuration
	Storage PersistentConfig `json:"storage,omitempty" yaml:"persistent,omitempty"`
}

func (node *HdfsNodeConfig) updateFromHdfsConfig(c *HdfsConfig, role string) error {
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

// HdfsConfig for hdfs-cluster
type HdfsConfig struct {
	Image *ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`

	Nodes []string `json:"nodes,omitempty" yaml:"nodes,omitempty" validate:"eq=0|min=3"`

	HdfsHome string `json:"hdfsHome" yaml:"-"`

	Namenode    *HdfsNodeConfig `json:"namenode,omitempty" yaml:"namenode" validate:"eq=0|eq=2"`
	Datanode    *HdfsNodeConfig `json:"datanode,omitempty" yaml:"datanode" validate:"eq=0|min=3"`
	Journalnode *HdfsNodeConfig `json:"journalnode,omitempty" yaml:"journalnode" validate:"eq=0|min=3"`

	Zookeeper *HdfsNodeConfig `json:"zookeeper,omitempty" yaml:"zookeeper" validate:"eq=0|min=3"`
}

// TODO: validate the yaml and nodes == 2 of namenode
func (v HdfsConfig) validate() error {
	return nil
}

// ********************************************
// HdfsChart for hdfs-cluster, implement Chart
type HdfsChart struct {
	ChartMeta

	values *HdfsConfig
}

// update each field value from global Config if that is ZERO
func (h HdfsChart) updateFromConfig(c Config) error {
	if c.Hdfs != nil {
		h.values = c.Hdfs
	}

	if c.Image != nil {
		if h.values.Image == nil {
			h.values.Image = &ImageConfig{}
			h.values.Image.updateFromConfig(c.Image)
		}
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

func (h HdfsChart) initLocalPvHome() error {
	localPvHome := fmt.Sprintf("%s/%s/{namenode,datanode,journalnode,zookeeper}", h.values.HdfsHome, HdfsClusterName)
	var host *Host
	var conn *Connection
	var err error
	for _, node := range h.values.Nodes {
		host = &Host{Address: node}
		conn, err = NewConnection(host)
		if err != nil {
			return errors.Wrap(err, "new connection failed")
		}
		if err := conn.Mkdir(localPvHome); err != nil {
			return err
		}
	}
	return nil
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

func (h *HdfsChart) getLabels() map[string]string {
	return map[string]string{
		HdfsInstanceLabelKey: h.ReleaseName,
	}
}

func NewHdfsChart(release string, c Config) *HdfsChart {
	h := &HdfsChart{}
	h.ChartName = HdfsClusterChart
	h.ReleaseName = release
	h.WaitingReady = true

	if c.Hdfs != nil {
		h.values = c.Hdfs
	} else {
		h.values = &HdfsConfig{}
	}
	return h
}
