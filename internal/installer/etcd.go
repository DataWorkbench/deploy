package installer

import (
	"encoding/json"
)

type EtcdConfig struct {
	Image *ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`

	// TODO: support resource configurations
	WorkloadConfig `json:",inline" yaml:",inline"`
}

// EtcdChart for etcd-cluster, implement Chart
type EtcdChart struct {
	ChartMeta `json:",inline" yaml:",inline"`
	values    *EtcdConfig `json:"config,omitempty" yaml:"config,omitempty"`
}

// update each field value from global Config if that is ZERO
func (e *EtcdChart) updateFromConfig(c Config) error {
	if c.Etcd != nil {
		e.values = c.Etcd
	}

	if c.Image != nil {
		if e.values.Image == nil {
			e.values.Image = &ImageConfig{}
			e.values.Image.updateFromConfig(c.Image)
		}
	}

	return e.values.Persistent.updateLocalPv(c.LocalPVHome, c.Nodes)
}

func (e EtcdChart) initLocalPvHome() error {
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

func NewEtcdChart(release string) *EtcdChart {
	e := &EtcdChart{}
	e.ChartName = EtcdClusterChart
	e.ReleaseName = release
	e.WaitingReady = true
	e.values = &EtcdConfig{}
	return e
}
