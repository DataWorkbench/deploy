package installer

import (
	"encoding/json"
)

type EtcdValuesConfig struct {
	Image *ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`

	Nodes      []string         `json:"nodes,omitempty" yaml:"nodes,omitempty" validate:"eq=0|min=3"`
	Persistent PersistentConfig `json:"persistent" yaml:"persistent"`
}

// EtcdChart for etcd-cluster, implement Chart
type EtcdChart struct {
	ChartMeta `json:",inline" yaml:",inline"`
	values    *EtcdValuesConfig `json:"config,omitempty" yaml:"config,omitempty"`
}

// update each field value from global Config if that is ZERO
func (e *EtcdChart) updateFromConfig(c Config) error {
	if c.Image != nil {
		if e.values.Image == nil {
			e.values.Image = &ImageConfig{}
		}
		e.values.Image.updateFromConfig(c.Image)
	}

	return e.values.Persistent.updateLocalPv(c.LocalPVHome, e.values.Nodes)
}

func (e *EtcdChart) parseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	bytes, err := json.Marshal(e.values)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}
