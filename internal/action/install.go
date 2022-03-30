package action

import (
	"context"
	"errors"
	"fmt"
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/DataWorkbench/deploy/internal/config"
	"github.com/DataWorkbench/deploy/internal/k8s/helm"
	"github.com/DataWorkbench/deploy/internal/k8s/helm/chart"
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

func Install(ctx context.Context, configFile string, services *[]string, debug, dryRun bool) error {
	config.Debug = debug
	config.DryRun = dryRun

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

	conf := &config.Config{}
	err := conf.Read(configFile, *logger)
	if err != nil {
		return err
	}


	// install operators
	for _, service := range *services {
		if common.StrSliceContains(Operators, service) {
			if err = installOperator(ctx, service, *conf, logger); err != nil{
				return err
			}
			logger.Info().Msg("**************************************************************").Fire()
		}
	}

	// install dependency service
	for _, service := range *services {
		if common.StrSliceContains(DependencyServices, service) {
			if err = installDependencyService(ctx, service, *conf, logger); err != nil{
				return err
			}
			logger.Info().Msg("**************************************************************").Fire()
		}
	}

	// install dataomnis
	if common.StrSliceContains(*services, common.DataomnisSystemName) {
		err = installDataomnis(ctx, common.DataomnisSystemName, *conf, logger)
		logger.Info().Msg("**************************************************************").Fire()
	}
	return nil
}


func installOperator(ctx context.Context, name string, c config.Config, logger *glog.Logger) error {
	var proxy *helm.Proxy
	var err error
	logger.Info().String("install operator", name).Msg("..").Fire()
	if proxy, err = helm.NewProxy(common.DefaultOperatorNamespace, logger); err != nil {
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

	var hc helm.Chart
	switch name {
	case common.HdfsOptName:
		hc = chart.NewHdfsOperatorChart(name)
	case common.MysqlOptName:
		hc = chart.NewMysqlOperatorChart(name)
	case common.RedisOptName:
		hc = chart.NewRedisOperatorChart(name)
	default:
		return errors.New(fmt.Sprintf("the service %s can not be installed", name))
	}

	if err = hc.UpdateFromConfig(c); err != nil {
		logger.Error().Error("update from config error", err).Fire()
		return err
	}
	if err = proxy.Install(ctx, hc); err != nil {
		logger.Error().Error("helm install error", err).Fire()
		return err
	}
	logger.Info().String("install operator", name).Msg(", done.").Fire()
	return nil
}


func installDependencyService(ctx context.Context, name string, c config.Config, logger *glog.Logger) error {
	var proxy *helm.Proxy
	var err error
	logger.Info().String("install dependency service", name).Msg("..").Fire()
	if proxy, err = helm.NewProxy(common.DefaultSystemNamespace, logger); err != nil {
		logger.Error().Error("create helm proxy error", err).Fire()
		return err
	}

	// check if exist
	// TODO: support check if ready(make sure exist and ready)
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

	var hc helm.Chart
	switch name {
	case common.EtcdClusterName:
		hc = chart.NewEtcdChart(name)
	case common.HdfsClusterName:
		hc = chart.NewHdfsChart(name)
	case common.MysqlClusterName:
		hc = chart.NewMysqlChart(name)
	case common.RedisClusterName:
		hc = chart.NewRedisChart(name)
	default:
		return errors.New(fmt.Sprintf("the service %s can not be installed", name))
	}

	if err = hc.UpdateFromConfig(c); err != nil {
		logger.Error().Error("update values from Config error", err).Fire()
		return err
	}
	logger.Debug().Any("chart values", hc).Fire()

	if err = hc.InitLocalDir(); err != nil {
		logger.Error().Error("init local dir error", err).Fire()
		return err
	}
	if err = proxy.Install(ctx, hc); err != nil {
		logger.Error().Error("helm install dependency service error", err).Fire()
		return err
	}
	logger.Info().String("install dependency service", name).Msg(", done.").Fire()
	return nil
}


func installDataomnis(ctx context.Context, name string, c config.Config, logger *glog.Logger) error {
	var proxy *helm.Proxy
	var err error
	logger.Info().Msg("install dataomnis ..").Fire()
	if proxy, err = helm.NewProxy(common.DefaultSystemNamespace, logger); err != nil {
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

	hc := chart.NewDataomnisChart(name)
	if err = hc.UpdateFromConfig(c); err != nil {
		logger.Error().Error("update values from Config error", err).Fire()
		return err
	}
	if err = hc.InitLocalDir(); err != nil {
		logger.Error().Error("init hostPath dir error", err).Fire()
		return err
	}
	if err = proxy.Install(ctx, hc); err != nil {
		logger.Error().Error("helm install dataomnis error", err).Fire()
		return err
	}
	logger.Info().Msg("install dataomnis, done.").Fire()
	return nil
}
