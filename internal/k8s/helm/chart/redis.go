package chart

import (
	"encoding/json"
	"fmt"
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/DataWorkbench/deploy/internal/config"
	"github.com/DataWorkbench/deploy/internal/ssh"
	"github.com/pkg/errors"
	"time"
)

// ********************************************
// RedisChart for redis-cluster, implement Chart
type RedisChart struct {
	ChartMeta
	Conf *config.RedisConfig
}

// update each field value from global Config
func (r *RedisChart) UpdateFromConfig(c config.Config) error {
	r.Conf = c.Redis

	if c.Image != nil {
		if r.Conf.Image == nil {
			r.Conf.Image = &config.Image{}
			r.Conf.Image.Copy(c.Image)
		}
	}

	r.Conf.Persistent.UpdateLocalPv(c.LocalPVHome, c.Nodes[:3])
	r.Conf.MasterSize = 3
	return r.Conf.Validate(r.ReleaseName)
}

func (r RedisChart) InitLocalDir() error {
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

func (r RedisChart) ParseValues() (Values, error) {
	var v Values = map[string]interface{}{}
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

func (r RedisChart) GetTimeoutSecond() time.Duration {
	if r.Conf.TimeoutSecond == 0 {
		return r.ChartMeta.GetTimeoutSecond()
	}
	return time.Duration(r.Conf.TimeoutSecond) * time.Second
}

func NewRedisChart(release string) *RedisChart {
	r := &RedisChart{}
	r.ChartName = common.RedisClusterChart
	r.ReleaseName = release
	r.Waiting = true
	return r
}
