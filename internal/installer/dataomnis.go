package installer

import (
	"github.com/DataWorkbench/common/lib/iaas"
)

type MetricsConfig struct {
	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

type GrpcLogConfig struct {
	Level     int8 `json:"level,omitempty" yaml:"level,omitempty"`
	Verbosity int8 `json:"verbosity,omitempty" yaml:"verbosity,omitempty"`
}

type ServiceConfig struct {
	LogLevel    int8          `json:"logLevel,omitempty" yaml:"logLevel"`
	GrpcLog     GrpcLogConfig `json:"grpcLog,omitempty" yaml:"grpcLog,omitempty"`
	ServiceName bool          `json:"serviceName,omitempty" yaml:"serviceName"`
	Metrics     MetricsConfig `json:"metrics,omitempty" yaml:"metrics,omitempty"`

	WorkloadConfig `json:",omitempty,inline" yaml:",omitempty,inline"`

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
	ExternalHost string `json:"externalHost" yaml:"-"`
	SecretName   string `json:"secretName" yaml:"-"`

	LogLevel        int8   `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	MaxIdleConn     int8   `json:"maxIdleConn,omitempty" yaml:"maxIdleConn,omitempty"`
	MaxOpenConn     int8   `json:"maxOpenConn,omitempty" yaml:"maxOpenConn,omitempty"`
	ConnMaxLifetime string `json:"connMaxLifetime,omitempty" yaml:"connMaxLifetime,omitempty"`
	SlowTshreshold  string `json:"slowTshreshold,omitempty" yaml:"slowTshreshold,omitempty"`
}

type DataomnisConfig struct {
	// dataomnis version
	Version string `json:"version" yaml:"version"`

	Domain string `json:"domain"  yaml:"domain"`
	Port   string `json:"port"    yaml:"port,omitempty"`

	// global configurations for all service as default
	Image  ImageConfig   `json:"image,omitempty" yaml:"image,omitempty"`
	Common ServiceConfig `json:"common,inline" yaml:",inline"`

	MysqlClient MysqlClientConfig `json:"mysql"   yaml:"mysqlCluster"`

	Iaas iaas.Config `json:"iaas,omitempty" yaml:"iaas,omitempty"`

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

	Jaeger         ServiceConfig        `json:"jaeger" yaml:"jaeger"`
	ServiceMonitor ServiceMonitorConfig `json:"serviceMonitor" yaml:"serviceMonitor"`
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
