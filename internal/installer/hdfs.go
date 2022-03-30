package installer

import (
	"encoding/json"
	"fmt"
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/DataWorkbench/deploy/internal/k8s/helm"
	"github.com/DataWorkbench/deploy/internal/ssh"
	"github.com/pkg/errors"
	"time"
)

const RoleNameNode = "namenode"

// common config for datanode / journalnode / namenode / zookeeper
type HdfsNodeConfig struct {
	// alias for Persistent, updated from configuration
	Storage common.Persistent `json:"storage,omitempty" yaml:"persistent,omitempty"`
}

func (node *HdfsNodeConfig) updateFromHdfsConfig(c *HdfsConfig, role string) error {
	if len(node.Storage.LocalPv.Nodes) < 3 {
		if len(c.Nodes) < 3 {
			return common.LeastNodeErr
		}

		if len(node.Storage.LocalPv.Nodes) < 1 {
			if role == RoleNameNode {
				node.Storage.LocalPv.Nodes = c.Nodes[len(c.Nodes)-2:]
			} else {
				node.Storage.LocalPv.Nodes = c.Nodes[:3]
			}
		}
	}
	return nil
}

// HdfsConfig for hdfs-cluster
type HdfsConfig struct {
	TimeoutSecond int `json:"-" yaml:"timeoutSecond,omitempty"`

	Image    *common.Image `json:"image,omitempty" yaml:"image,omitempty"`
	Nodes    []string      `json:"nodes,omitempty" yaml:"nodes,omitempty" validate:"eq=0|min=3"`
	HdfsHome string        `json:"hdfsHome"        yaml:"-"`

	Namenode    *HdfsNodeConfig `json:"namenode,omitempty"    yaml:"namenode"    validate:"eq=0|eq=2"`
	Datanode    *HdfsNodeConfig `json:"datanode,omitempty"    yaml:"datanode"    validate:"eq=0|min=3"`
	Journalnode *HdfsNodeConfig `json:"journalnode,omitempty" yaml:"journalnode" validate:"eq=0|min=3"`

	Zookeeper *HdfsNodeConfig `json:"zookeeper,omitempty"     yaml:"zookeeper"   validate:"eq=0|min=3"`
}

// TODO: validate the yaml and nodes == 2 of namenode
func (v HdfsConfig) validate() error {
	return nil
}

// ********************************************
// HdfsChart for hdfs-cluster, implement Chart
type HdfsChart struct {
	helm.ChartMeta
	Conf *HdfsConfig
}

// update each field value from global Config if that is ZERO
func (h *HdfsChart) UpdateFromConfig(c common.Config) error {
	if c.Hdfs != nil {
		h.Conf = c.Hdfs
	}

	if c.Image != nil {
		if h.Conf.Image == nil {
			h.Conf.Image = &common.Image{}
			h.Conf.Image.Copy(c.Image)
		}
	}

	// TODO: check if localPv exist and start with c.LocalPVHome
	h.Conf.HdfsHome = fmt.Sprintf(common.LocalHomeFmt, c.LocalPVHome)

	// update datanode / namenode / journalnode / zookeeper conf
	var err error
	if err = h.Conf.Datanode.updateFromHdfsConfig(h.Conf, ""); err != nil {
		return err
	}
	if err = h.Conf.Journalnode.updateFromHdfsConfig(h.Conf, ""); err != nil {
		return err
	}
	if err = h.Conf.Namenode.updateFromHdfsConfig(h.Conf, RoleNameNode); err != nil {
		return err
	}
	err = h.Conf.Zookeeper.updateFromHdfsConfig(h.Conf, "")
	return err
}

func (h HdfsChart) InitLocalDir() error {
	dnLocalPvDir := fmt.Sprintf("%s/%s/datanode", h.Conf.HdfsHome, common.HdfsClusterName)
	jnLocalPvDir := fmt.Sprintf("%s/%s/journalnode", h.Conf.HdfsHome, common.HdfsClusterName)
	nnLocalPvDir := fmt.Sprintf("%s/%s/namenode", h.Conf.HdfsHome, common.HdfsClusterName)
	zkLocalPvDir := fmt.Sprintf("%s/%s/zookeeper", h.Conf.HdfsHome, common.HdfsClusterName)
	var host *ssh.Host
	var conn *ssh.Connection
	var err error
	for _, node := range h.Conf.Datanode.Storage.LocalPv.Nodes {
		host = &ssh.Host{Address: node}
		conn, err = ssh.NewConnection(host)
		if err != nil {
			return errors.Wrap(err, "new connection failed")
		}
		if _, err := conn.Mkdir(dnLocalPvDir); err != nil {
			return err
		}
	}
	for _, node := range h.Conf.Journalnode.Storage.LocalPv.Nodes {
		host = &ssh.Host{Address: node}
		conn, err = ssh.NewConnection(host)
		if err != nil {
			return errors.Wrap(err, "new connection failed")
		}
		if _, err := conn.Mkdir(jnLocalPvDir); err != nil {
			return err
		}
	}
	for _, node := range h.Conf.Namenode.Storage.LocalPv.Nodes {
		host = &ssh.Host{Address: node}
		conn, err = ssh.NewConnection(host)
		if err != nil {
			return errors.Wrap(err, "new connection failed")
		}
		if _, err := conn.Mkdir(nnLocalPvDir); err != nil {
			return err
		}
	}
	for _, node := range h.Conf.Zookeeper.Storage.LocalPv.Nodes {
		host = &ssh.Host{Address: node}
		conn, err = ssh.NewConnection(host)
		if err != nil {
			return errors.Wrap(err, "new connection failed")
		}
		if _, err := conn.Mkdir(zkLocalPvDir); err != nil {
			return err
		}
	}
	return nil
}

func (h HdfsChart) ParseValues() (helm.Values, error) {
	var v helm.Values = map[string]interface{}{}
	bytes, err := json.Marshal(h.Conf)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}

func (h HdfsChart) GetLabels() map[string]string {
	return map[string]string{
		common.HdfsInstanceLabelKey: h.ReleaseName,
	}
}

func (h HdfsChart) GetTimeoutSecond() time.Duration {
	if h.Conf.TimeoutSecond == 0 {
		return h.ChartMeta.GetTimeoutSecond()
	}
	return time.Duration(h.Conf.TimeoutSecond) * time.Second
}

func NewHdfsChart(release string, c common.Config) *HdfsChart {
	h := &HdfsChart{}
	h.ChartName = common.HdfsClusterChart
	h.ReleaseName = release
	h.Waiting = true

	if c.Hdfs != nil {
		h.Conf = c.Hdfs
	} else {
		h.Conf = &HdfsConfig{}
	}
	return h
}


// ***********************************************************
type HdfsClient struct {
	ConfigmapName string `json:"configmapName" yaml:"-"`
}
