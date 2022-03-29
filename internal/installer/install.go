package installer

import (
	"context"
	"errors"
	"fmt"
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/DataWorkbench/deploy/internal/k8s/helm"
	"github.com/DataWorkbench/glog"
)

var Operators = []string{
	common.HdfsOptName,
	common.MysqlOptName,
	common.RedisOptName,
}

var DependencyServices = []string{
	common.EtcdClusterName,
	common.MysqlClusterName,
	common.RedisClusterName,
	common.HdfsClusterName,
}

var AllServices []string

func init() {
	AllServices = append(AllServices, Operators...)
	AllServices = append(AllServices, DependencyServices...)
	AllServices = append(AllServices, common.DataomnisSystemName)
}

func Install(configFile string, services *[]string, debug, dryRun bool) error {
	ctx := context.Background()
	logger := glog.NewDefault()
	if debug {
		logger = logger.WithLevel(glog.DebugLevel)
	}

	// check
	for _, s := range *services {
		if !common.StrSliceContains(AllServices, s) {
			msg := fmt.Sprintf("The service:%s cat not be installed.", s)
			logger.Error().Msg(msg).Fire()
			return errors.New(msg)
		}
	}

	conf := &common.Config{Debug: debug}
	err := conf.Read(configFile, *logger)
	if err != nil {
		return err
	}


	// install operators
	for _, service := range *services {
		if common.StrSliceContains(Operators, service) {
			if err = installOperator(ctx, service, *conf, logger, debug, dryRun); err != nil{
				return err
			}
			logger.Info().Msg("**************************************************************").Fire()
		}
	}

	// install dependency service
	for _, service := range *services {
		if common.StrSliceContains(DependencyServices, service) {
			if err = installDependencyService(ctx, service, *conf, logger, debug, dryRun); err != nil{
				return err
			}
			logger.Info().Msg("**************************************************************").Fire()
		}
	}

	// install dataomnis
	if common.StrSliceContains(*services, common.DataomnisSystemName) {
		err = installDataomnis(ctx, common.DataomnisSystemName, *conf, logger, debug, dryRun)
		logger.Info().Msg("**************************************************************").Fire()
	}
	return nil
}


func installOperator(ctx context.Context, name string, c common.Config, logger *glog.Logger, debug, dryRun bool) error {
	var proxy *helm.Proxy
	var err error
	logger.Info().String("new helm proxy with namespace", common.DefaultOperatorNamespace).Fire()
	if proxy, err = helm.NewProxy(ctx, common.DefaultOperatorNamespace, logger, debug, dryRun); err != nil {
		logger.Error().Error("create helm proxy to install operators error", err).Fire()
		return err
	}

	var installed bool
	installed, err = proxy.Exist(name)
	if err != nil {
		logger.Debug().Error("get release error", err).Fire()
		return err
	}
	if installed {
		logger.Warn().String("operator", name).Msg("was installed, ignore.").Fire()
		return nil
	}

	var chart helm.Chart
	switch name {
	case common.HdfsOptName:
		chart = NewHdfsOperatorChart(name, c)
	case common.MysqlOptName:
		chart = NewMysqlOperatorChart(name, c)
	case common.RedisOptName:
		chart = NewRedisOperatorChart(name, c)
	default:
		return errors.New(fmt.Sprintf("the service %s can not be installed", name))
	}

	logger.Info().String("install operator", name).Msg("..").Fire()
	if err = proxy.Install(chart); err != nil {
		logger.Error().Error("install operator error", err).Fire()
		return err
	}
	logger.Info().String("install operator", name).Msg(", done.").Fire()
	return nil
}


func installDependencyService(ctx context.Context, name string, c common.Config, logger *glog.Logger, debug, dryRun bool) error {
	var proxy *helm.Proxy
	var err error
	logger.Info().String("install dependency service", name).Msg("..").Fire()
	if proxy, err = helm.NewProxy(ctx, common.DefaultSystemNamespace, logger, debug, dryRun); err != nil {
		logger.Error().Error("create helm proxy error", err).Fire()
		return err
	}

	var installed bool
	installed, err = proxy.Exist(name)
	if err != nil {
		logger.Debug().Error("get release error", err).Fire()
		return err
	}
	if installed {
		logger.Warn().String("dependency service", name).Msg("was installed, ignore.").Fire()
		return nil
	}

	var chart helm.Chart
	switch name {
	case common.EtcdClusterName:
		chart = NewEtcdChart(name, c)
	case common.HdfsClusterName:
		chart = NewHdfsChart(name, c)
	case common.MysqlClusterName:
		chart = NewMysqlChart(name, c)
	case common.RedisClusterName:
		chart = NewRedisChart(name, c)
	default:
		return errors.New(fmt.Sprintf("the service %s can not be installed", name))
	}

	if err = chart.UpdateFromConfig(c); err != nil {
		logger.Error().Error("update values from Config error", err).Fire()
		return err
	}
	logger.Debug().Any("chart values", chart).Fire()

	if err = chart.InitLocalPvDir(); err != nil {
		logger.Error().Error("init local pv home error", err).Fire()
		return err
	}
	if err = proxy.Install(chart); err != nil {
		logger.Error().Error("helm install dependency service error", err).Fire()
		return err
	}
	logger.Info().String("install dependency service", name).Msg(", done.").Fire()
	return nil
}


func installDataomnis(ctx context.Context, name string, c common.Config, logger *glog.Logger, debug, dryRun bool) error {
	var proxy *helm.Proxy
	var err error
	logger.Info().Msg("install dataomnis ..").Fire()
	if proxy, err = helm.NewProxy(ctx, common.DefaultSystemNamespace, logger, debug, dryRun); err != nil {
		logger.Error().Error("create helm proxy error", err).Fire()
		return err
	}

	var installed bool
	installed, err = proxy.Exist(name)
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
	if err = proxy.Install(chart); err != nil {
		logger.Error().Error("helm install dataomnis error", err).Fire()
		return err
	}
	logger.Info().Msg("install dataomnis, done.").Fire()
	return nil
}
