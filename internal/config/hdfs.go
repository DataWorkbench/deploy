package config

import (
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/pkg/errors"
)

const RoleNameNode = "Namenode"
const RoleDataNode = "Datanode"
const RoleJournalNode = "Journalnode"
const RoleZookeeper = "Zookeeper"

// common config for datanode / journalnode / namenode / zookeeper
type HdfsNodeConfig struct {
	Role string
	// alias for Persistent, updated from configuration
	Storage *Persistent `json:"storage,omitempty" yaml:"persistent,omitempty"`
}

func (node *HdfsNodeConfig) Update(releaseName string, c *HdfsConfig) error {
	if node.Role == RoleNameNode {
		if len(node.Storage.LocalPv.Nodes) == 0 { // update Namenode
			if len(c.Nodes) < 3 {
				return errors.Errorf(common.AtLeaseNodeNumFmt, 3, releaseName)
			}
			node.Storage.LocalPv.Nodes = c.Nodes[len(c.Nodes)-2:]
		} else if len(node.Storage.LocalPv.Nodes) != 2 { // check: only 2 nodes for namenode
			return errors.Errorf(common.OnlyNodeNumFmt, 2, node.Role)
		}
	} else {
		if len(node.Storage.LocalPv.Nodes) == 0 { // update others
			if len(c.Nodes) < 3 {
				return errors.Errorf(common.AtLeaseNodeNumFmt, 3, releaseName)
			}
			node.Storage.LocalPv.Nodes = c.Nodes[:3]
		} else if len(node.Storage.LocalPv.Nodes) < 3 { // check: at lease 3 nodes for others
			return errors.Errorf(common.AtLeaseNodeNumFmt, 3, node.Role)
		}
	}
	return nil
}

// HdfsConfig for hdfs-cluster
type HdfsConfig struct {
	TimeoutSecond int `json:"-" yaml:"timeoutSecond,omitempty"`

	Image    *Image   `json:"image,omitempty" yaml:"image,omitempty"`
	Nodes    []string `json:"nodes,omitempty" yaml:"nodes,omitempty" validate:"eq=0|min=3"`
	HdfsHome string   `json:"hdfsHome"        yaml:"-"`

	Zookeeper   *HdfsNodeConfig `json:"zookeeper,omitempty"   yaml:"zookeeper"   validate:"required"`
	Journalnode *HdfsNodeConfig `json:"journalnode,omitempty" yaml:"journalnode" validate:"required"`
	Namenode    *HdfsNodeConfig `json:"namenode,omitempty"    yaml:"namenode"    validate:"required"`
	Datanode    *HdfsNodeConfig `json:"datanode,omitempty"    yaml:"datanode"    validate:"required"`
}

// ***********************************************************
type HdfsClient struct {
	ConfigmapName string `json:"configmapName" yaml:"-"`
}
