package installer

import (
	"encoding/json"
)

// RedisConfig for hdfs-cluster
type RedisConfig struct {
	Image *ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`

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

	values *RedisConfig
}

// update each field value from global Config if that is ZERO
func (r *RedisChart) updateConfig(c Config) error {
	if c.Redis != nil {
		r.values = c.Redis
	}

	if c.Image != nil {
		if r.values.Image == nil {
			r.values.Image = &ImageConfig{}
			r.values.Image.updateFromConfig(c.Image)
		}
	}
	return r.values.Redis.Persistent.updateLocalPv(c.LocalPVHome, c.Nodes)
}

func (r RedisChart) parseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	bytes, err := json.Marshal(r.values)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}

func NewRedisChart(release string, c Config) *RedisChart {
	r := &RedisChart{}
	r.ChartName = RedisClusterChart
	r.ReleaseName = release
	r.WaitingReady = true

	if c.Redis != nil {
		r.values = c.Redis
	} else {
		r.values = &RedisConfig{}
	}
	return r
}
