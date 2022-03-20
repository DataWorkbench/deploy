package installer

import (
	"encoding/json"
	"fmt"
)

type MysqlConfig struct {
	Image *ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`

	Pxc *WorkloadConfig `json:"pxc" yaml:"pxc"`
}

// MysqlChart for etcd-cluster, implement Chart
type MysqlChart struct {
	ChartMeta `json:",inline" yaml:",inline"`

	values *MysqlConfig `json:"config,omitempty" yaml:"config,omitempty"`
}

// update each field value from global Config if that is ZERO
func (m *MysqlChart) updateFromConfig(c Config) error {
	if c.Mysql != nil {
		m.values = c.Mysql
	}

	if c.Image != nil {
		if m.values.Image == nil {
			m.values.Image = &ImageConfig{}
			m.values.Image.updateFromConfig(c.Image)
		}
	}

	return m.values.Pxc.Persistent.updateLocalPv(c.LocalPVHome, c.Nodes)
}

func (m *MysqlChart) initLocalPvHome() error {
	localPvHome := fmt.Sprintf("%s/%s/{data,log,mysql-bin}", m.values.Pxc.Persistent.LocalPv.Home, MysqlClusterName)
	for _, node := range m.values.Pxc.Persistent.LocalPv.Nodes {
		if err := CreateRemoteDir(node, localPvHome); err != nil {
			return err
		}
	}
	return nil
}

func (m *MysqlChart) parseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	bytes, err := json.Marshal(m.values)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}

func (m MysqlChart) getLabels() map[string]string {
	return map[string]string{
		InstanceLabelKey: fmt.Sprintf(MysqlInstanceLabelValueFmt, m.ReleaseName),
	}
}

func NewMysqlChart(release string, c Config) *MysqlChart {
	m := &MysqlChart{}
	m.ChartName = MysqlClusterChart
	m.ReleaseName = release
	m.WaitingReady = true

	if c.Mysql != nil {
		m.values = c.Mysql
	} else {
		m.values = &MysqlConfig{}
	}
	return m
}
