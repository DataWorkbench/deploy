package installer

import (
	"context"
	"errors"
	"fmt"
	"github.com/DataWorkbench/glog"
	"github.com/go-playground/validator/v10"
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
	if debug {
		logger = logger.WithLevel(glog.DebugLevel)
	}

	var err error
	// check configuration file
	_, err = os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("the configuration file: %s not exist!", configFile)
			logger.Error().Msg(msg).Fire()
			err = errors.New(msg)
			return err
		}
		err = nil
	}

	conf := &Config{}
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Error().String("failed to read configuration file", configFile).Error("error", err).Fire()
		logger.Error().Msg("please make sure the file is YAML format.").Fire()
		return err
	}
	if err = yaml.Unmarshal(bytes, conf); err != nil {
		logger.Error().Error("parse bytes from the configuration to yaml error", err).Fire()
		return err
	}
	logger.Debug().Any("Configuration", conf).Fire()
	// validate
	if err = validator.New().Struct(conf); err != nil {
		logger.Error().Error("validate configuration error", err).Fire()
		return err
	}

	// install operators
	if err = installOperators(ctx, logger, debug, *conf); err != nil {
		return err
	}


	// install etcd-cluster
	if err = installDatabases(ctx, logger, debug, *conf); err != nil {
		return err
	}
	return nil
}

func installOperators(ctx context.Context, logger *glog.Logger, debug bool, c Config) error {
	var helm *Proxy
	var err error
	if helm, err = NewProxy(ctx, DefaultOperatorNamespace, logger, debug); err != nil {
		logger.Error().Error("create helm proxy to install operators error", err).Fire()
		return err
	}

	hdfsOperator := NewChartMeta(HdfsOperatorChart, HdfsOperatorName, true)
	hdfsOperator.updateFromConfig(c)
	if err = helm.install(hdfsOperator); err != nil {
		logger.Error().Error("install hdfs-operator error", err).Fire()
		return err
	}

	mysqlOperator := NewChartMeta(MysqlOperatorChart, MysqlOperatorName, true)
	mysqlOperator.updateFromConfig(c)
	if err = helm.install(mysqlOperator); err != nil {
		logger.Error().Error("install mysql-operator error", err).Fire()
		return err
	}

	redisOperator := NewChartMeta(RedisOperatorChart, RedisOperatorName, true)
	redisOperator.updateFromConfig(c)
	if err = helm.install(redisOperator); err != nil {
		logger.Error().Error("install redis-operator error", err).Fire()
		return err
	}
	return nil
}

func installDatabases(ctx context.Context, logger *glog.Logger, debug bool, c Config) error {
	var helm *Proxy
	var err error
	if helm, err = NewProxy(ctx, DefaultSystemNamespace, logger, debug); err != nil {
		logger.Error().Error("create helm proxy to install operators error", err).Fire()
		return err
	}

	hdfs := &HdfsChart{}
	hdfs.setMeta(HdfsClusterChart, HdfsClusterName, debug)
	if err = hdfs.updateFromConfig(c); err != nil {
		logger.Error().Error("update hdfs values from Config error", err).Fire()
		return err
	}
	if err = helm.install(hdfs); err != nil {
		logger.Error().Error("helm install hdfs-cluster error", err).Fire()
		return err
	}
	return nil
}
