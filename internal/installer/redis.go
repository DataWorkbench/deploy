package installer

import (
	"encoding/json"
	"fmt"
)


// RedisConfig for hdfs-cluster
type RedisConfig struct {
	Image *ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`

	WorkloadConfig `json:",inline" yaml:",inline"`

	MasterSize int `json:"masterSize,omitempty" yaml:"-"`
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

	if err := r.values.Persistent.updateLocalPv(c.LocalPVHome, c.Nodes); err != nil {
		return err
	}

	r.values.MasterSize = len(r.values.Persistent.LocalPv.Nodes)
	return nil
}

func (r *RedisChart) initLocalPvHome() error {
	localPvHome := fmt.Sprintf("%s/%s/{01,02}", r.values.Persistent.LocalPv.Home, RedisClusterName)
	for _, node := range r.values.Persistent.LocalPv.Nodes {
		if err := CreateRemoteDir(node, localPvHome); err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisChart) parseValues() (Values, error) {
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
