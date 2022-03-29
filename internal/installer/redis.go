package installer

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// RedisConfig for hdfs-cluster
type RedisConfig struct {
	MasterSize int    `json:"masterSize"      yaml:"-"`
	Image      *Image `json:"image,omitempty" yaml:"image,omitempty"`

	Workload `json:",inline" yaml:",inline"`
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

// update each field value from global Config
func (r *RedisChart) updateFromConfig(c Config) error {
	if c.Redis != nil {
		r.values = c.Redis
	}

	if c.Image != nil {
		if r.values.Image == nil {
			r.values.Image = &Image{}
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
	var host *Host
	var conn *Connection
	var err error
	for _, node := range r.values.Persistent.LocalPv.Nodes {
		host = &Host{Address: node}
		conn, err = NewConnection(host)
		if err != nil {
			return errors.Wrap(err, "new connection failed")
		}

		if _, err := conn.Mkdir(localPvHome); err != nil {
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

func (r *RedisChart) getLabels() map[string]string {
	return map[string]string{
		RedisInstanceLabelKey: r.ReleaseName,
	}
}

func (r *RedisChart) getTimeoutSecond() int {
	if r.values.TimeoutSecond == 0 {
		return r.ChartMeta.getTimeoutSecond()
	}
	return r.values.TimeoutSecond
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


// *************************************************************************************
type RedisClient struct {
	Mode         string `json:"mode,omitempty"         yaml:"mode,omitempty"`
	SentinelAddr string `json:"sentinelAddr,omitempty" yaml:"sentinelAddr,omitempty"`
	ClusterAddr  string `json:"clusterAddr,omitempty"  yaml:"clusterAddr,omitempty"`
	MasterName   string `json:"masterName,omitempty"   yaml:"masterName,omitempty"`
	Database     string `json:"database,omitempty"     yaml:"database,omitempty"`
	Username     string `json:"username,omitempty"     yaml:"username,omitempty"`
	Password     string `json:"password,omitempty"     yaml:"password,omitempty"`
}

func (r *RedisClient) generateAddr(releaseName string, size int) {
	var addrs []string
	if r.Mode == RedisClusterModeCluster && r.ClusterAddr == "" { // internal redis-cluster
		for i := 0; i < size; i++ {
			addrs = append(addrs, fmt.Sprintf(RedisClusterAddrFmt, releaseName, i, RedisClusterPort))
		}
		r.ClusterAddr = strings.Join(addrs, ",")
	}
}
