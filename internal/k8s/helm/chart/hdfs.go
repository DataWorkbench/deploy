package chart

import (
	"encoding/json"
	"fmt"
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/DataWorkbench/deploy/internal/config"
	"github.com/DataWorkbench/deploy/internal/ssh"
	"github.com/pkg/errors"
	"time"
)

// ********************************************
// HdfsChart for hdfs-cluster, implement Chart
type HdfsChart struct {
	ChartMeta
	Conf *config.HdfsConfig
}

// update each field value from global Config if that is ZERO
func (h *HdfsChart) UpdateFromConfig(c config.Config) error {
	h.Conf = c.Hdfs
	h.Conf.Namenode.Role = config.RoleNameNode
	h.Conf.Datanode.Role = config.RoleDataNode
	h.Conf.Journalnode.Role = config.RoleJournalNode
	h.Conf.Zookeeper.Role = config.RoleZookeeper

	if c.Image != nil {
		if h.Conf.Image == nil {
			h.Conf.Image = &config.Image{}
			h.Conf.Image.Copy(c.Image)
		}
	}

	// TODO: check if localPv exist and start with c.LocalPVHome
	h.Conf.HdfsHome = fmt.Sprintf(common.LocalHomeFmt, c.LocalPVHome)

	// update datanode / namenode / journalnode / zookeeper conf
	var err error
	if err = h.Conf.Datanode.Update(h.ReleaseName, h.Conf); err != nil {
		return err
	}
	if err = h.Conf.Journalnode.Update(h.ReleaseName, h.Conf); err != nil {
		return err
	}
	if err = h.Conf.Namenode.Update(h.ReleaseName, h.Conf); err != nil {
		return err
	}
	err = h.Conf.Zookeeper.Update(h.ReleaseName, h.Conf)
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

func (h HdfsChart) ParseValues() (Values, error) {
	var v Values = map[string]interface{}{}
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

func NewHdfsChart(release string) *HdfsChart {
	h := &HdfsChart{}
	h.ChartName = common.HdfsClusterChart
	h.ReleaseName = release
	h.Waiting = true
	return h
}
