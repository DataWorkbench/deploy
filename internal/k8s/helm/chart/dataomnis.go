package chart

import (
	"encoding/json"
	"fmt"
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/DataWorkbench/deploy/internal/config"
	"github.com/DataWorkbench/deploy/internal/ssh"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)


type DataomnisChart struct {
	ChartMeta

	Conf *config.DataomnisConfig
}

// update each field value from global Config if that is ZERO
func (d *DataomnisChart) UpdateFromConfig(c config.Config) error {
	d.Conf = c.Dataomnis

	if c.Image != nil {
		if d.Conf.Image == nil {
			d.Conf.Image = &config.Image{}
		}
		d.Conf.Image.Copy(c.Image)
	}
	d.Conf.Image.Tag = d.Conf.Version

	if d.Conf.MysqlClient == nil {
		d.Conf.MysqlClient = &config.MysqlClient{}
	}
	d.Conf.MysqlClient.Update(common.MysqlClusterName)

	if d.Conf.RedisClient == nil {
		d.Conf.RedisClient = &config.RedisClient{Mode: common.RedisClusterModeCluster}
	}
	d.Conf.RedisClient.GenerateAddr(common.RedisClusterName, 3)

	d.Conf.EtcdClient = &config.EtcdClient{
		Endpoint: common.EtcdClusterName,
	}

	d.Conf.HdfsClient = &config.HdfsClient{
		ConfigmapName: fmt.Sprintf(common.HdfsConfigMapFmt, common.HdfsClusterName),
	}

	// update hostPath(log-dir) for values
	if d.Conf.Persistent == nil {
		d.Conf.Persistent = &config.Persistent{}
	}
	d.Conf.Persistent.HostPath = fmt.Sprintf(common.DataomnisHostPathFmt, c.LocalPVHome, d.GetReleaseName())

	d.Conf.Nodes = c.Nodes

	if d.Conf.Apiglobal.Enabled {
		d.Conf.Apiglobal.UpdateRegion()
		d.Conf.Apiglobal.UpdateAuthentication()
	}

	if config.Debug {
		data, err := yaml.Marshal(d.Conf)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(common.TmpValuesFile, data, 0777)
	}

	return nil
}

// TODO: rsync flink-xxx.tgz to helmCache
func (d DataomnisChart) InitLocalDir() error {
	localLogDir := fmt.Sprintf("%s/log/{account,apiglobal,apiserver,enginemanager,resourcemanager,scheduler,spacemanager,notifier}", d.Conf.Persistent.HostPath)
	helmCacheDir := fmt.Sprintf(common.DefaultHelmRepoCacheFmt, "/root")

	var host *ssh.Host
	var conn *ssh.Connection
	var err error
	for _, node := range d.Conf.Nodes {
		host = &ssh.Host{Address: node}
		conn, err = ssh.NewConnection(host)
		if err != nil {
			return errors.Wrap(err, "new connection failed")
		}
		if _, err := conn.Mkdir(localLogDir); err != nil {
			return err
		}
		if _, err := conn.Mkdir(helmCacheDir); err != nil {
			return err
		}
	}
	return nil
}

func (d DataomnisChart) ParseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	bytes, err := json.Marshal(d.Conf)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}

func NewDataomnisChart(release string) *DataomnisChart {
	d := &DataomnisChart{}
	d.ChartName = common.DataomnisSystemChart
	d.ReleaseName = release
	d.Waiting = true
	return d
}
