package installer

import (
	"encoding/json"
	"fmt"
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/DataWorkbench/deploy/internal/k8s/helm"
	"github.com/DataWorkbench/deploy/internal/ssh"
	"github.com/pkg/errors"
	"strings"
)

// RedisConfig for hdfs-cluster
type RedisConfig struct {
	MasterSize int           `json:"masterSize"      yaml:"-"`
	Image      *common.Image `json:"image,omitempty" yaml:"image,omitempty"`

	common.Workload `json:",inline" yaml:",inline"`
}

// TODO: validate the yaml and nodes == 3 by default
func (v RedisConfig) validate() error {
	return nil
}

// ********************************************
// RedisChart for hdfs-cluster, implement Chart
type RedisChart struct {
	helm.ChartMeta
	Conf *RedisConfig
}

// update each field value from global Config
func (r *RedisChart) UpdateFromConfig(c common.Config) error {
	if c.Redis != nil {
		r.Conf = c.Redis
	}

	if c.Image != nil {
		if r.Conf.Image == nil {
			r.Conf.Image = &common.Image{}
			r.Conf.Image.Copy(c.Image)
		}
	}

	if err := r.Conf.Persistent.UpdateLocalPv(c.LocalPVHome, c.Nodes); err != nil {
		return err
	}

	r.Conf.MasterSize = len(r.Conf.Persistent.LocalPv.Nodes)
	return nil
}

func (r RedisChart) InitLocalPvDir() error {
	localPvHome := fmt.Sprintf("%s/%s/{01,02}", r.Conf.Persistent.LocalPv.Home, common.RedisClusterName)
	var host *ssh.Host
	var conn *ssh.Connection
	var err error
	for _, node := range r.Conf.Persistent.LocalPv.Nodes {
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

func (r RedisChart) ParseValues() (helm.Values, error) {
	var v helm.Values = map[string]interface{}{}
	bytes, err := json.Marshal(r.Conf)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}

func (r RedisChart) GetLabels() map[string]string {
	return map[string]string{
		common.RedisInstanceLabelKey: r.ReleaseName,
	}
}

func (r RedisChart) GetTimeoutSecond() int {
	if r.Conf.TimeoutSecond == 0 {
		return r.ChartMeta.GetTimeoutSecond()
	}
	return r.Conf.TimeoutSecond
}

func NewRedisChart(release string, c common.Config) *RedisChart {
	r := &RedisChart{}
	r.ChartName = common.RedisClusterChart
	r.ReleaseName = release
	r.Waiting = true

	if c.Redis != nil {
		r.Conf = c.Redis
	} else {
		r.Conf = &RedisConfig{}
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
	if r.Mode == common.RedisClusterModeCluster && r.ClusterAddr == "" { // internal redis-cluster
		for i := 0; i < size; i++ {
			addrs = append(addrs, fmt.Sprintf(common.RedisClusterAddrFmt, releaseName, i, common.RedisClusterPort))
		}
		r.ClusterAddr = strings.Join(addrs, ",")
	}
}
