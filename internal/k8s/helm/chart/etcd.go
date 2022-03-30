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


// EtcdChart for etcd-cluster, implement Chart
type EtcdChart struct {
	ChartMeta `json:",inline" yaml:",inline"`
	Conf      *config.EtcdConfig `yaml:"config,omitempty"`
}

// update each field value from global Config if that is ZERO
func (e *EtcdChart) UpdateFromConfig(c config.Config) error {
	e.Conf = c.Etcd

	if c.Image != nil {
		if e.Conf.Image == nil {
			e.Conf.Image = &config.Image{}
			e.Conf.Image.Copy(c.Image)
		}
	}

	e.Conf.Persistent.UpdateLocalPv(c.LocalPVHome, c.Nodes[:3])
	return e.Conf.Validate(e.ReleaseName)
}

func (e EtcdChart) InitLocalDir() error {
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

func (e EtcdChart) ParseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	bytes, err := json.Marshal(e.Conf)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}

func (e EtcdChart) GetTimeoutSecond() time.Duration {
	if e.Conf.TimeoutSecond == 0 {
		return e.ChartMeta.GetTimeoutSecond()
	}
	return time.Duration(e.Conf.TimeoutSecond) * time.Second
}

func NewEtcdChart(release string) *EtcdChart {
	e := &EtcdChart{}
	e.ChartName = common.EtcdClusterChart
	e.ReleaseName = release
	e.Waiting = true
	return e
}
