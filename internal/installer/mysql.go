package installer

import (
	"encoding/json"
)

type MysqlValuesConfig struct {
	Image *ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`

	Nodes []string `json:"-" yaml:"nodes,omitempty"`

	Pxc *WorkloadConfig `json:"pxc" yaml:"pxc"`
}

// MysqlChart for etcd-cluster, implement Chart
type MysqlChart struct {
	ChartMeta `json:",inline" yaml:",inline"`

	values *MysqlValuesConfig `json:"config,omitempty" yaml:"config,omitempty"`
}

// update each field value from global Config if that is ZERO
func (m *MysqlChart) updateFromConfig(c Config) error {
	if c.Image != nil {
		if m.values.Image == nil {
			m.values.Image = &ImageConfig{}
		}
		m.values.Image.updateFromConfig(c.Image)
	}

	return m.values.Pxc.Persistent.updateLocalPv(c.LocalPVHome, m.values.Nodes)
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
