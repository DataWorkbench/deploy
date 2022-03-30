package config

import (
	"fmt"
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/pkg/errors"
	"strings"
)

// RedisConfig for redis-cluster
type RedisConfig struct {
	MasterSize int    `json:"masterSize"      yaml:"-"`
	Image      *Image `json:"image,omitempty" yaml:"image,omitempty"`

	Workload `json:",inline" yaml:",inline"`
}

func (v RedisConfig) Validate(releaseName string) error {
	if len(v.Persistent.LocalPv.Nodes) != 3 {
		return errors.Errorf(common.OnlyNodeNumFmt, 3, releaseName)
	}
	return nil
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

func (r *RedisClient) GenerateAddr(releaseName string, size int) {
	var addrs []string
	if r.Mode == common.RedisClusterModeCluster && r.ClusterAddr == "" { // internal redis-cluster
		for i := 0; i < size; i++ {
			addrs = append(addrs, fmt.Sprintf(common.RedisClusterAddrFmt, releaseName, i, common.RedisClusterPort))
		}
		r.ClusterAddr = strings.Join(addrs, ",")
	}
}
