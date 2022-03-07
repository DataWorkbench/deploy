package installer

import (
	"encoding/json"
)

// RedisConfig for hdfs-cluster
type RedisConfig struct {
	Image *ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`

	Nodes []string `json:"nodes,omitempty" yaml:"nodes,omitempty" validate:"eq=0|min=3"`

	Redis WorkloadConfig `json:"redis,omitempty" yaml:"redis,omitempty"`
}

// TODO: validate the yaml and nodes == 3 by default
func (v RedisConfig) validate() error {
	return nil
}

// ********************************************
// RedisChart for hdfs-cluster, implement Chart
type RedisChart struct {
	ChartMeta

	values RedisConfig
}

// update each field value from global Config if that is ZERO
func (h RedisChart) updateConfig(c Config) error {
	if c.Image != nil {
		if h.values.Image == nil {
			h.values.Image = &ImageConfig{}
		}
		h.values.Image.updateFromConfig(c.Image)
	}
	return h.values.Redis.Persistent.updateLocalPv(c.LocalPVHome, h.values.Nodes)
}

func (h RedisChart) parseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	bytes, err := json.Marshal(h.values)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}
