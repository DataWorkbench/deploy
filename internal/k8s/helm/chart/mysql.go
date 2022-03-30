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

// MysqlChart for mysql-cluster, implement Chart
type MysqlChart struct {
	ChartMeta `json:",inline" yaml:",inline"`

	Conf *config.MysqlConfig `json:"config,omitempty" yaml:"config,omitempty"`
}

// update each field value from global Config if that is ZERO
func (m *MysqlChart) UpdateFromConfig(c config.Config) error {
	m.Conf = c.Mysql

	if c.Image != nil {
		if m.Conf.Image == nil {
			m.Conf.Image = &config.Image{}
			m.Conf.Image.Copy(c.Image)
		}
	}

	m.Conf.Pxc.Persistent.UpdateLocalPv(c.LocalPVHome, c.Nodes[:3])
	return m.Conf.Validate(m.ReleaseName)
}

func (m MysqlChart) InitLocalDir() error {
	localPvHome := fmt.Sprintf("%s/%s/{data,log,mysql-bin}", m.Conf.Pxc.Persistent.LocalPv.Home, common.MysqlClusterName)
	var host *ssh.Host
	var conn *ssh.Connection
	var err error
	for _, node := range m.Conf.Pxc.Persistent.LocalPv.Nodes {
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

func (m MysqlChart) ParseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	bytes, err := json.Marshal(m.Conf)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}

func (m MysqlChart) GetLabels() map[string]string {
	return map[string]string{
		common.InstanceLabelKey: fmt.Sprintf(common.MysqlInstanceLabelValueFmt, m.ReleaseName),
	}
}

func (m MysqlChart) GetTimeoutSecond() time.Duration {
	if m.Conf.TimeoutSecond == 0 {
		return m.ChartMeta.GetTimeoutSecond()
	}
	return time.Duration(m.Conf.TimeoutSecond) * time.Second
}

func NewMysqlChart(release string) *MysqlChart {
	m := &MysqlChart{}
	m.ChartName = common.MysqlClusterChart
	m.ReleaseName = release
	m.Waiting = true
	return m
}
