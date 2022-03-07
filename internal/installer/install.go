package installer

import (
	"context"
	"errors"
	"github.com/DataWorkbench/glog"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

const (
	DefaultOperatorNamespace = "dataomnis-operator"
	DefaultSystemNamespace   = "dataomnis-system"

	// flink helm chart name
	FlinkChart = "flink-0.1.6.tgz"
)

func InitConfiguration() {

}

func Install(configFile string, debug bool) error {
	ctx := context.Background()
	logger := glog.NewDefault()

	var err error
	// check configuration file
	_, err = os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Printf("The configuration file: %s not exist!", configFile)
			logger.Error().Msg(msg).Fire()
			err = errors.New(msg)
			return err
		}
		err = nil
	}

	conf := &Config{}
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Error().String("Failed to read configuration file", configFile).Error("error", err).Fire()
		logger.Error().Msg("please make sure the file is YAML format.").Fire()
		return err
	}
	if err = yaml.Unmarshal(bytes, conf); err != nil {
		logger.Error().Msg("Failed to parse bytes from the configuration to yaml!").Fire()
		return err
	}

	// install operators
	if err = installOperators(ctx, logger, debug); err != nil {
		return err
	}


	// install etcd-cluster

	return nil
}

func installOperators(ctx context.Context, logger *glog.Logger, debug bool) error {
	var helm *Proxy
	var err error
	if helm, err = NewProxy(ctx, DefaultOperatorNamespace, debug, *logger); err != nil {
		logger.Error().Error("create helm proxy to install operators error", err).Fire()
		return err
	}

	hdfsOperator := NewChartMeta(HdfsOperatorChart, HdfsOperatorName, true)
	if err = helm.install(hdfsOperator); err != nil {
		logger.Error().Error("install hdfs-operator error", err).Fire()
		return err
	}

	mysqlOperator := NewChartMeta(MysqlOperatorChart, MysqlOperatorName, true)
	if err = helm.install(mysqlOperator); err != nil {
		logger.Error().Error("install mysql-operator error", err).Fire()
		return err
	}

	redisOperator := NewChartMeta(RedisOperatorChart, RedisOperatorName, true)
	if err = helm.install(redisOperator); err != nil {
		logger.Error().Error("install redis-operator error", err).Fire()
		return err
	}
	return nil
}