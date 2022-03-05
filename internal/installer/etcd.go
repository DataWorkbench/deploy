package installer

import (
	"encoding/json"
	"errors"
	"fmt"
)

type EtcdValuesConfig struct {
	ValuesConfig `json:",inline" yaml:",inline"`

	Nodes      []string           `json:"nodes,omitempty" yaml:"nodes,omitempty" validate:"eq=0|min=3"`
	Persistent LocalStorageConfig `json:"persistent" yaml:"storage"`
}

// EtcdChart for etcd-cluster, implement Chart
type EtcdChart struct {
	ChartMeta `json:",inline" yaml:",inline"`
	values EtcdValuesConfig `json:"config,omitempty" yaml:"config,omitempty"`
}

// update each field value from global Config if that is ZERO
func (e *EtcdChart) updateFromConfig(c Config) error {
	if c.Image != nil {
		if e.values.Image == nil {
			e.values.Image = &ImageConfig{}
		}
		e.values.Image.updateFromConfig(c.Image)
	}

	// TODO: check if localPv exist and start with c.LocalPVHome
	e.values.Persistent.LocalHome = fmt.Sprintf(LocalHomeFmt, c.LocalPVHome)

	if e.values.Nodes == nil {
		if len(c.Nodes) < 3 {
			return errors.New("at least 3 nodes are required for etcd-cluster")
		}
		// Default: select pre-three nodes to install etcd-cluster
		e.values.Nodes = c.Nodes[:3]
	}
	return nil
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
