package installer

import (
	"context"
	"errors"
	"fmt"
	"github.com/DataWorkbench/glog"
)

var Operators = []string{
	HdfsOptName,
	MysqlOptName,
	RedisOptName,
}

var DependencyServices = []string{
	EtcdClusterName,
	HdfsClusterName,
	MysqlClusterName,
	RedisClusterName,
}

var AllServices []string

func init() {
	AllServices = append(AllServices, Operators...)
	AllServices = append(AllServices, DependencyServices...)
}

func Install(configFile string, debug bool, services *[]string) error {
	ctx := context.Background()
	logger := glog.NewDefault()
	if debug {
		logger = logger.WithLevel(glog.DebugLevel)
	}

	// check
	for _, s := range *services {
		if !StrContains(AllServices, s) {
			msg := fmt.Sprintf("The service:%s cat not be installed.", s)
			logger.Error().Msg(msg).Fire()
			return errors.New(msg)
		}
	}

	conf := &Config{}
	err := conf.Read(configFile, *logger)
	if err != nil {
		return err
	}


	// install operators
	for _, service := range *services {
		if StrContains(Operators, service) {
			if err = installOperator(ctx, service, logger, debug, *conf); err != nil{
				return err
			}
		}
	}

	// install dependency service
	for _, service := range *services {
		if StrContains(DependencyServices, service) {
			if err = installDependencyService(ctx, service, logger, debug, *conf); err != nil{
				return err
			}
		}
	}
	return nil
}


func installOperator(ctx context.Context, name string, logger *glog.Logger, debug bool, c Config) error {
	var helm *Proxy
	var err error
	logger.Info().String("new helm proxy with namespace", DefaultOperatorNamespace).Fire()
	if helm, err = NewProxy(ctx, DefaultOperatorNamespace, logger, debug); err != nil {
		logger.Error().Error("create helm proxy to install operators error", err).Fire()
		return err
	}

	var chart Chart
	switch name {
	case HdfsOptName:
		chart = NewHdfsOperatorChart(name, c)
	case MysqlOptName:
		chart = NewMysqlOperatorChart(name, c)
	case RedisOptName:
		chart = NewRedisOperatorChart(name, c)
	}

	logger.Info().String("install operator", name).Msg("..").Fire()
	if err = helm.install(chart); err != nil {
		logger.Error().Error("install operator error", err).Fire()
		return err
	}
	logger.Info().String("install operator", name).Msg(", done.").Fire()
	return nil
}


func installDependencyService(ctx context.Context, name string, logger *glog.Logger, debug bool, c Config) error {
	var helm *Proxy
	var err error
	logger.Info().String("install dependency service", name).Msg("..").Fire()
	if helm, err = NewProxy(ctx, DefaultSystemNamespace, logger, debug); err != nil {
		logger.Error().Error("create helm proxy error", err).Fire()
		return err
	}

	var chart Chart
	switch name {
	case EtcdClusterName:
		chart = NewEtcdChart(name)
	case HdfsClusterName:
		chart = NewHdfsChart(name)
	case MysqlClusterName:
		chart = NewMysqlChart(name)
	case RedisClusterName:
		chart = NewRedisChart(name)
	}
	if err = chart.updateFromConfig(c); err != nil {
		logger.Error().Error("update values from Config error", err).Fire()
		return err
	}
	if err = helm.install(chart); err != nil {
		logger.Error().Error("helm install dependency service error", err).Fire()
		return err
	}
	logger.Info().String("install dependency service", name).Msg(", done.").Fire()
	return nil
}

func StrContains(ss []string, s string) bool {
	for _, _s := range ss {
		if _s == s {
			return true
		}
	}
	return false
}
