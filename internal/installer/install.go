package installer

import (
	"context"
	"github.com/DataWorkbench/glog"
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

	conf := &Config{}
	err := conf.Read(configFile, *logger)
	if err != nil {
		return err
	}

	// install operators
	logger.Info().Msg("install operators ..").Fire()
	if err = installOperators(ctx, logger, debug, *conf); err != nil {
		return err
	}

	logger.Info().Msg("install database-services ..").Fire()
	if err = installDatabases(ctx, logger, debug, *conf); err != nil {
		return err
	}
	return nil
}

func installOperators(ctx context.Context, logger *glog.Logger, debug bool, c Config) error {
	var helm *Proxy
	var err error
	logger.Info().String("new helm proxy with namespace", DefaultOperatorNamespace).Fire()
	if helm, err = NewProxy(ctx, DefaultOperatorNamespace, logger, debug); err != nil {
		logger.Error().Error("create helm proxy to install operators error", err).Fire()
		return err
	}

	logger.Info().Msg("install hdfs-operator ..").Fire()
	hdfsOperator := NewHdfsOperatorChart(HdfsOptName, c)
	if err = helm.install(hdfsOperator); err != nil {
		logger.Error().Error("install hdfs-operator error", err).Fire()
		return err
	}
	logger.Info().String("install hdfs-operator, done.", HdfsOptName).Fire()

	logger.Info().Msg("install mysql-operator ..").Fire()
	mysqlOperator := NewMysqlOperatorChart(MysqlOptName, c)
	if err = helm.install(mysqlOperator); err != nil {
		logger.Error().Error("install mysql-operator error", err).Fire()
		return err
	}
	logger.Info().Msg("install mysql-operator, done.").Fire()

	//logger.Info().Msg("install redis-operator ..").Fire()
	//redisOperator := NewRedisOperatorChart(RedisOptName, c)
	//if err = helm.install(redisOperator); err != nil {
	//	logger.Error().Error("install redis-operator error", err).Fire()
	//	return err
	//}
	//logger.Info().Msg("install redis-operator, done.").Fire()
	return nil
}

func installDatabases(ctx context.Context, logger *glog.Logger, debug bool, c Config) error {
	var (
	 	helm *Proxy
		err error
	)
	logger.Info().String("new helm proxy with namespace", DefaultSystemNamespace).Fire()
	if helm, err = NewProxy(ctx, DefaultSystemNamespace, logger, debug); err != nil {
		logger.Error().Error("create helm proxy to install operators error", err).Fire()
		return err
	}

	logger.Info().Msg("install etcd-cluster ..").Fire()
	etcd := NewEtcdChart(EtcdClusterName)
	if err = etcd.updateFromConfig(c); err != nil {
		logger.Error().Error("update etcd values from Config error", err).Fire()
		return err
	}
	if err = helm.install(etcd); err != nil {
		logger.Error().Error("helm install etcd-cluster error", err).Fire()
		return err
	}
	logger.Info().Msg("install etcd-cluster, done.").Fire()

	logger.Info().Msg("install hdfs-cluster ..").Fire()
	hdfs := NewHdfsChart(HdfsClusterName)
	if err = hdfs.updateFromConfig(c); err != nil {
		logger.Error().Error("update hdfs values from Config error", err).Fire()
		return err
	}
	if err = helm.install(hdfs); err != nil {
		logger.Error().Error("helm install hdfs-cluster error", err).Fire()
		return err
	}
	logger.Info().Msg("install hdfs-cluster, done.").Fire()

	logger.Info().Msg("install mysql-cluster ..").Fire()
	mysql := NewMysqlChart(MysqlClusterName)
	if err = mysql.updateFromConfig(c); err != nil {
		logger.Error().Error("update mysql values from Config error", err).Fire()
		return err
	}
	if err = helm.install(mysql); err != nil {
		logger.Error().Error("helm install mysql-cluster error", err).Fire()
		return err
	}
	logger.Info().Msg("install mysql-cluster, done.").Fire()
	return nil
}
