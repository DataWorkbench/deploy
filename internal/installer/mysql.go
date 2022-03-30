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

//TODO: add backup-config
type MysqlConfig struct {
	TimeoutSecond int              `json:"-" yaml:"timeoutSecond,omitempty"`
	Image         *common.Image    `json:"image,omitempty" yaml:"image,omitempty"`
	Pxc           *common.Workload `json:"pxc" yaml:"pxc"`
}

// MysqlChart for etcd-cluster, implement Chart
type MysqlChart struct {
	helm.ChartMeta `json:",inline" yaml:",inline"`

	Conf *MysqlConfig `json:"config,omitempty" yaml:"config,omitempty"`
}

// update each field value from global Config if that is ZERO
func (m *MysqlChart) UpdateFromConfig(c common.Config) error {
	if c.Mysql != nil {
		m.Conf = c.Mysql
	}

	if c.Image != nil {
		if m.Conf.Image == nil {
			m.Conf.Image = &common.Image{}
			m.Conf.Image.Copy(c.Image)
		}
	}

	return m.Conf.Pxc.Persistent.UpdateLocalPv(c.LocalPVHome, c.Nodes)
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

func (m MysqlChart) ParseValues() (helm.Values, error) {
	var v helm.Values = map[string]interface{}{}
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

func NewMysqlChart(release string, c common.Config) *MysqlChart {
	m := &MysqlChart{}
	m.ChartName = common.MysqlClusterChart
	m.ReleaseName = release
	m.Waiting = true

	if c.Mysql != nil {
		m.Conf = c.Mysql
	} else {
		m.Conf = &MysqlConfig{}
	}
	return m
}


// ***********************************************************
type MysqlClient struct {
	ExternalHost string `json:"externalHost" yaml:"-"`
	SecretName   string `json:"secretName" yaml:"-"`

	LogLevel        int8   `json:"logLevel,omitempty"        yaml:"logLevel,omitempty"`
	MaxIdleConn     int32  `json:"maxIdleConn,omitempty"     yaml:"maxIdleConn,omitempty"`
	MaxOpenConn     int32  `json:"maxOpenConn,omitempty"     yaml:"maxOpenConn,omitempty"`
	ConnMaxLifetime string `json:"connMaxLifetime,omitempty" yaml:"connMaxLifetime,omitempty"`
	SlowThreshold   string `json:"slowThreshold,omitempty"   yaml:"slowThreshold,omitempty"`
}

func (c *MysqlClient) update(releaseName string)  {
	c.ExternalHost = fmt.Sprintf(common.MysqlExternalHostFmt, releaseName)
	c.SecretName = fmt.Sprintf(common.MysqlSecretNameFmt, releaseName)
}
