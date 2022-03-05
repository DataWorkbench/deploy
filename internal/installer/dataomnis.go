package installer

import (
	"github.com/DataWorkbench/common/lib/iaas"
)

type ServiceConfig struct {
	LogLevel         int8 `json:"logLevel,omitempty" yaml:"logLevel"`
	GrpcLogLevel     int8 `json:"grpcLogLevel,omitempty" yaml:"grpcLogLevel"`
	GrpcLogVerbosity int8 `json:"grpcLogVerbosity,omitempty" yaml:"grpcLogVerbosity"`

	ServiceName    bool `json:"serviceName,omitempty" yaml:"serviceName"`
	MetricsEnabled bool `json:"metricsEnabled,omitempty" yaml:"metricsEnabled"`

	WorkloadConfig `json:",omitempty,inline" yaml:",inline"`

	Envs []map[string]string `json:"envs,omitempty" yaml:"envs,flow"`
}

type Webservice struct {
	Enabled bool `yaml:"enabled"`
}

type ApiGlobalConfig struct {
	Enabled       bool     `yaml:"enabled" yaml:"enabled"`
	Regions       []Region `json:"regions,omitempty" yaml:"regions,flow"`
	ServiceConfig `json:",,omitemptyinline" yaml:",inline"`
}

type Region struct {
	Hosts string `json:"hosts,omitempty" yaml:"hosts"`
	EnUs  string `json:"en_us,omitempty" yaml:"en_us_name"`
	ZhCn  string `json:"zh_cn,omitempty" yaml:"zh_cn_name"`
}

type ServiceMonitorConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Interval string `yaml:"interval"`
}

type MysqlClientConfig struct {
	LogLevel        int8   `json:"logLevel,omitempty" yaml:"logLevel"`
	MaxIdleConn     int8   `json:"maxIdleConn,omitempty" yaml:"maxIdleConn"`
	MaxOpenConn     int8   `json:"maxOpenConn,omitempty" yaml:"maxOpenConn"`
	ConnMaxLifetime string `json:"connMaxLifetime,omitempty" yaml:"connMaxLifetime"`
	SlowTshreshold  string `json:"slowTshreshold,omitempty" yaml:"slowTshreshold"`
}

type DataomnisChart struct {
	ChartMeta `json:",inline" yaml:",inline"`

	// dataomnis version
	Version string `json:"version" yaml:"version"`

	Domain string `json:"domain"  yaml:"domain"`
	Port   string `json:"port"    yaml:"port"`

	// global configurations for all service as default
	Image         ImageConfig `json:"image" yaml:"image"`
	Workload      `json:",inline" yaml:",inline"`
	ServiceConfig `json:",inline" yaml:",inline"`

	MysqlClient MysqlClientConfig `json:"mysql"   yaml:"mysqlCluster"`

	Iaas iaas.Config `json:"iaas" yaml:"iaas"`

	WebService      Webservice      `json:"webservice" yaml:"webservice"`
	Apiglobal       ApiGlobalConfig `json:"apiglobal" yaml:"apiGlobal"`
	Apiserver       ServiceConfig   `json:"apiserver" yaml:"apiserver"`
	Account         ServiceConfig   `json:"account" yaml:"account"`
	Developer       ServiceConfig   `json:"developer" yaml:"developer"`
	Enginemanager   ServiceConfig   `json:"enginemanager" yaml:"enginemanager"`
	Jobmanager      ServiceConfig   `json:"jobmanager" yaml:"jobmanager"`
	Resourcemanager ServiceConfig   `json:"resourcemanager" yaml:"resourcemanager"`
	Scheduler       ServiceConfig   `json:"scheduler" yaml:"scheduler"`
	Spacemanager    ServiceConfig   `json:"spacemanager" yaml:"spacemanager"`

	Jaeger         Workload             `json:"jaeger" yaml:"jaeger"`
	ServiceMonitor ServiceMonitorConfig `json:"serviceMonitor" yaml:"serviceMonitor"`
}
