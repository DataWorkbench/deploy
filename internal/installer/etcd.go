package installer

import (
	"encoding/json"
	"fmt"
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/DataWorkbench/deploy/internal/k8s/helm"
	"github.com/DataWorkbench/deploy/internal/ssh"
	"github.com/pkg/errors"
)

type EtcdConfig struct {
	Image *common.Image `json:"image,omitempty" yaml:"image,omitempty"`

	common.Workload `json:",inline" yaml:",inline"`
}

// EtcdChart for etcd-cluster, implement Chart
type EtcdChart struct {
	helm.ChartMeta `json:",inline" yaml:",inline"`
	Conf           *EtcdConfig `yaml:"config,omitempty"`
}

// update each field value from global Config if that is ZERO
func (e *EtcdChart) UpdateFromConfig(c common.Config) error {
	if c.Etcd != nil {
		e.Conf = c.Etcd
	}

	if c.Image != nil {
		if e.Conf.Image == nil {
			e.Conf.Image = &common.Image{}
			e.Conf.Image.Copy(c.Image)
		}
	}

	return e.Conf.Persistent.UpdateLocalPv(c.LocalPVHome, c.Nodes)
}

func (e EtcdChart) InitLocalPvDir() error {
	localPvHome := fmt.Sprintf("%s/%s", e.Conf.Persistent.LocalPv.Home, common.EtcdClusterName)
	var host *ssh.Host
	var conn *ssh.Connection
	var err error
	for _, node := range e.Conf.Persistent.LocalPv.Nodes {
		host = &ssh.Host{Address: node}
		conn, err = ssh.NewConnection(host)
		if err != nil {
			return errors.Wrap(err, "new connection failed")
		}
		if _, err := conn.Mkdir(localPvHome); err != nil {
			return err
		}
	}
	return nil
}

func (e EtcdChart) ParseValues() (helm.Values, error) {
	var v helm.Values = map[string]interface{}{}
	bytes, err := json.Marshal(e.Conf)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}

func (e EtcdChart) GetTimeoutSecond() int {
	if e.Conf.TimeoutSecond == 0 {
		return e.ChartMeta.GetTimeoutSecond()
	}
	return e.Conf.TimeoutSecond
}

func NewEtcdChart(release string, c common.Config) *EtcdChart {
	e := &EtcdChart{}
	e.ChartName = common.EtcdClusterChart
	e.ReleaseName = release
	e.Waiting = true

	if c.Etcd != nil {
		e.Conf = c.Etcd
	} else {
		e.Conf = &EtcdConfig{}
	}
	return e
}

// ***********************************************************
type EtcdClient struct {
	Endpoint string `json:"endpoint" yaml:"-"`
}
