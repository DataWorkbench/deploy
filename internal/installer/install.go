package installer

import (
	"context"
	"errors"
	"github.com/DataWorkbench/glog"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
)

const (
	DefaultOperatorNamespace = "dataomnis-operator"
	DefaultSystemNamespace   = "dataomnis-system"

	// flink helm chart name
	FlinkChart = "flink-0.1.6.tgz"
)

func InitConfiguration() {

}

func Install(configFile string, debug bool) (err error) {
	ctx := context.Background()
	logger := glog.NewDefault()

	// check configuration file
	_, err = os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Printf("The configuration file: %s not exist!", configFile)
			logger.Error().Msg(msg).Fire()
			err = errors.New(msg)
			return
		}
		err = nil
	}

	conf := &Config{}
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Error().String("Failed to read configuration file", configFile).Error(", error", err).Fire()
		return
	}
	if err = yaml.Unmarshal(bytes, conf); err != nil {
		logger.Error().Msg("Failed to parse bytes from the configuration to yaml!").Fire()
		return
	}

	var helm *Proxy
	if helm, err = NewProxy(ctx, DefaultOperatorNamespace, debug, *logger); err != nil {
		logger.Error().Error("create helm proxy error", err).Fire()
		return
	}

	// install operators
	operatorValues := c.
	if err = helm.install(HdfsOperatorName, HdfsOperatorChart, imageValuesStr, true); err != nil {
		logger.Error().Msg("Failed to install hdfs-operator.").Fire()
		return
	}
	if err = helm.install(MysqlOperatorName, MysqlOperatorChart, imageValuesStr, true); err != nil {
		logger.Error().Error("install mysql-operator error", err).Fire()
		return
	}
	if err = helm.install(RedisOperatorName, RedisOperatorChart, imageValuesStr, true); err != nil {
		logger.Error().Error("install redis-operator error", err).Fire()
		return
	}

	// install etcd-cluster

}