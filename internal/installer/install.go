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
	MysqlClusterName,
	RedisClusterName,
	HdfsClusterName,
}

var AllServices []string

func init() {
	AllServices = append(AllServices, Operators...)
	AllServices = append(AllServices, DependencyServices...)
	AllServices = append(AllServices, DataomnisSystemName)
}

func Install(configFile string, services *[]string, debug, dryRun bool) error {
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

	conf := &Config{Debug: debug}
	err := conf.Read(configFile, *logger)
	if err != nil {
		return err
	}


	// install operators
	for _, service := range *services {
		if StrContains(Operators, service) {
			if err = installOperator(ctx, service, *conf, logger, debug, dryRun); err != nil{
				return err
			}
			logger.Info().Msg("**************************************************************").Fire()
		}
	}

	// install dependency service
	for _, service := range *services {
		if StrContains(DependencyServices, service) {
			if err = installDependencyService(ctx, service, *conf, logger, debug, dryRun); err != nil{
				return err
			}
			logger.Info().Msg("**************************************************************").Fire()
		}
	}

	// install dataomnis
	if StrContains(*services, DataomnisSystemName) {
		err = installDataomnis(ctx, DataomnisSystemName, *conf, logger, debug, dryRun)
		logger.Info().Msg("**************************************************************").Fire()
	}
	return nil
}


func installOperator(ctx context.Context, name string, c Config, logger *glog.Logger, debug, dryRun bool) error {
	var helm *Proxy
	var err error
	logger.Info().String("new helm proxy with namespace", DefaultOperatorNamespace).Fire()
	if helm, err = NewProxy(ctx, DefaultOperatorNamespace, logger, debug, dryRun); err != nil {
		logger.Error().Error("create helm proxy to install operators error", err).Fire()
		return err
	}

	var installed bool
	installed, err = helm.exist(name)
	if err != nil {
		logger.Debug().Error("get release error", err).Fire()
		return err
	}
	if installed {
		logger.Warn().String("operator", name).Msg("was installed, ignore.").Fire()
		return nil
	}

	var chart Chart
	switch name {
	case HdfsOptName:
		chart = NewHdfsOperatorChart(name, c)
	case MysqlOptName:
		chart = NewMysqlOperatorChart(name, c)
	case RedisOptName:
		chart = NewRedisOperatorChart(name, c)
	default:
		return errors.New(fmt.Sprintf("the service %s can not be installed", name))
	}

	logger.Info().String("install operator", name).Msg("..").Fire()
	if err = helm.install(chart); err != nil {
		logger.Error().Error("install operator error", err).Fire()
		return err
	}
	logger.Info().String("install operator", name).Msg(", done.").Fire()
	return nil
}


func installDependencyService(ctx context.Context, name string, c Config, logger *glog.Logger, debug, dryRun bool) error {
	var helm *Proxy
	var err error
	logger.Info().String("install dependency service", name).Msg("..").Fire()
	if helm, err = NewProxy(ctx, DefaultSystemNamespace, logger, debug, dryRun); err != nil {
		logger.Error().Error("create helm proxy error", err).Fire()
		return err
	}

	var installed bool
	installed, err = helm.exist(name)
	if err != nil {
		logger.Debug().Error("get release error", err).Fire()
		return err
	}
	if installed {
		logger.Warn().String("dependency service", name).Msg("was installed, ignore.").Fire()
		return nil
	}

	var chart Chart
	switch name {
	case EtcdClusterName:
		chart = NewEtcdChart(name, c)
	case HdfsClusterName:
		chart = NewHdfsChart(name, c)
	case MysqlClusterName:
		chart = NewMysqlChart(name, c)
	case RedisClusterName:
		chart = NewRedisChart(name, c)
	default:
		return errors.New(fmt.Sprintf("the service %s can not be installed", name))
	}

	if err = chart.updateFromConfig(c); err != nil {
		logger.Error().Error("update values from Config error", err).Fire()
		return err
	}
	logger.Debug().Any("chart values", chart).Fire()

	if err = chart.initLocalPvHome(); err != nil {
		logger.Error().Error("init local pv home error", err).Fire()
		return err
	}
	if err = helm.install(chart); err != nil {
		logger.Error().Error("helm install dependency service error", err).Fire()
		return err
	}
	logger.Info().String("install dependency service", name).Msg(", done.").Fire()
	return nil
}


func installDataomnis(ctx context.Context, name string, c Config, logger *glog.Logger, debug, dryRun bool) error {
	var helm *Proxy
	var err error
	logger.Info().Msg("install dataomnis ..").Fire()
	if helm, err = NewProxy(ctx, DefaultSystemNamespace, logger, debug, dryRun); err != nil {
		logger.Error().Error("create helm proxy error", err).Fire()
		return err
	}

	var installed bool
	installed, err = helm.exist(name)
	if err != nil {
		logger.Debug().Error("get release error", err).Fire()
		return err
	}
	if installed {
		logger.Warn().Msg("dataomnis-system service was installed, ignore.").Fire()
		return nil
	}

	chart := NewDataomnisChart(name, c)
	if err = chart.updateFromConfig(c); err != nil {
		logger.Error().Error("update values from Config error", err).Fire()
		return err
	}
	if err = chart.initHostPathDir(c); err != nil {
		logger.Error().Error("init hostPath dir error", err).Fire()
		return err
	}
	if err = helm.install(chart); err != nil {
		logger.Error().Error("helm install dataomnis error", err).Fire()
		return err
	}
	logger.Info().Msg("install dataomnis, done.").Fire()
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
