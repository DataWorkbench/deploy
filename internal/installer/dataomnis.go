package installer

import (
	"github.com/DataWorkbench/common/lib/iaas"
)

// global config from dataomnis-conf.yaml
type Dataomnis struct {
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

// overwritten by global configuration
func (i ImageConfig) overwritten(source ImageConfig) {
	if source.Registry != "" {
		i.Registry = source.Registry
	}
	if len(source.PullSecrets) > 0 {
		i.PullSecrets = source.PullSecrets
	}
	if source.PullPolicy != "" {
		i.PullPolicy = source.PullPolicy
	}
}
