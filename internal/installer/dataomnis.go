package installer

import (
	"encoding/json"
	"fmt"
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
	LogLevel    int8           `json:"logLevel,omitempty" yaml:"logLevel"`
	GrpcLog     *GrpcLogConfig `json:"grpcLog,omitempty" yaml:"grpcLog,omitempty"`
	Metrics     *MetricsConfig `json:"metrics,omitempty" yaml:"metrics,omitempty"`

	WorkloadConfig `json:",omitempty,inline" yaml:",omitempty,inline"`

	Envs []map[string]string `json:"envs,omitempty" yaml:"envs,flow"`
}

type Webservice struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// TODO: update from webservice enable
// TODO: generate default region
type ApiGlobalConfig struct {
	Enabled bool     `json:"enabled" yaml:"enabled" yaml:"enabled"`
	Regions []Region `json:"-" yaml:"regions,flow,omitempty"`

	RegionValues []map[string]RegionValue `json:"regions" yaml:"-"`

	ServiceConfig `json:",omitempty,inline" yaml:",inline"`
}

func (a ApiGlobalConfig) updateRegion() {
	for _, r := range a.Regions {
		rv := RegionValue{
			Host: r.Host,
			Name: Names{
				ZhCn: r.ZhCn,
				EnUs: r.EnUs,
			},
		}
		a.RegionValues = append(a.RegionValues, map[string]RegionValue{
			r.EnUs: rv,
		})
	}
}

// Region for configurations, parse to RegionValue.
type Region struct {
	Host string `json:"-" yaml:"host"`
	EnUs string `json:"-" yaml:"enUsName"`
	ZhCn string `json:"-" yaml:"zhCnName"`
}

// ***********************************************************

// ***********************************************************
// RegionValue for apiglobal values.yaml, updated from Region.
type Names struct {
	ZhCn string `json:"zh_cn" yaml:"-"`
	EnUs string `json:"en_us" yaml:"-"`
}

type RegionValue struct {
	Host string `json:"hosts" yaml:"-"`
	Name Names  `json:"names"`
}

// ***********************************************************

type ServiceMonitorConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Interval string `yaml:"interval"`
}

type MysqlClientConfig struct {
	ExternalHost string `json:"externalHost" yaml:"-"`
	SecretName   string `json:"secretName" yaml:"-"`

	LogLevel        int8   `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	MaxIdleConn     int32  `json:"maxIdleConn,omitempty" yaml:"maxIdleConn,omitempty"`
	MaxOpenConn     int32  `json:"maxOpenConn,omitempty" yaml:"maxOpenConn,omitempty"`
	ConnMaxLifetime string `json:"connMaxLifetime,omitempty" yaml:"connMaxLifetime,omitempty"`
	SlowTshreshold  string `json:"slowTshreshold,omitempty" yaml:"slowTshreshold,omitempty"`
}

type EtcdClientConfig struct {
	Endpoint string `json:"endpoint" yaml:"-"`
}

type HdfsClientConfig struct {
	ConfigmapName string `json:"configmapName" yaml:"-"`
}

type RedisClientConfig struct {
	Address string `json:"address" yaml:"-"`
}

type DataomnisConfig struct {
	// dataomnis version
	Version string `json:"version" yaml:"version"`

	Domain string `json:"domain"  yaml:"domain"`
	Port   string `json:"port"    yaml:"port,omitempty"`

	// global configurations for all service as default
	Image *ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`

	MysqlClient *MysqlClientConfig `json:"mysql"   yaml:"mysql"`
	EtcdClient  *EtcdClientConfig  `json:"etcd" yaml:"-"`
	HdfsClient  *HdfsClientConfig  `json:"hdfs" yaml:"-"`
	RedisClient *RedisClientConfig `json:"redis" yaml:"-"`

	Iaas *iaas.Config `json:"iaas,omitempty" yaml:"iaas,omitempty" validate:"omitempty"`

	Common *ServiceConfig `json:"common,inline" yaml:",inline"`

	WebService      *Webservice      `json:"webservice" yaml:"webservice"`
	Apiglobal       *ApiGlobalConfig `json:"apiglobal" yaml:"apiGlobal"`
	Apiserver       *ServiceConfig   `json:"apiserver" yaml:"apiserver"`
	Account         *ServiceConfig   `json:"account" yaml:"account"`
	Developer       *ServiceConfig   `json:"developer" yaml:"developer"`
	Enginemanager   *ServiceConfig   `json:"enginemanager" yaml:"enginemanager"`
	Jobmanager      *ServiceConfig   `json:"jobmanager" yaml:"jobmanager"`
	Resourcemanager *ServiceConfig   `json:"resourcemanager" yaml:"resourcemanager"`
	Scheduler       *ServiceConfig   `json:"scheduler" yaml:"scheduler"`
	Spacemanager    *ServiceConfig   `json:"spacemanager" yaml:"spacemanager"`

	Jaeger         *ServiceConfig        `json:"jaeger" yaml:"jaeger"`
	ServiceMonitor *ServiceMonitorConfig `json:"serviceMonitor" yaml:"serviceMonitor"`
}

type DataomnisChart struct {
	ChartMeta

	values *DataomnisConfig
}

// update each field value from global Config if that is ZERO
func (d *DataomnisChart) updateFromConfig(c Config) error {
	if c.Image != nil {
		if d.values.Image == nil {
			d.values.Image = &ImageConfig{}
		}
		d.values.Image.updateFromConfig(c.Image)
	}
	d.values.Image.Tag = d.values.Version
	
	d.values.MysqlClient.ExternalHost = fmt.Sprintf(MysqlExternalHostFmt, MysqlClusterName)
	d.values.EtcdClient.Endpoint = EtcdClusterName
	d.values.HdfsClient.ConfigmapName = fmt.Sprintf(HdfsConfigMapFmt, HdfsClusterName)
	d.values.RedisClient.Address = fmt.Sprintf(RedisAddressFmt, RedisClusterName)

	d.values.Apiglobal.updateRegion()
	return nil
}

func (d *DataomnisChart) parseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	bytes, err := json.Marshal(d.values)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}

func NewDataomnisChart(release string, c Config) *DataomnisChart {
	d := &DataomnisChart{}
	d.ChartName = DataomnisSystemChart
	d.ReleaseName = release
	d.WaitingReady = true

	if c.Dataomnis != nil {
		d.values = c.Dataomnis
	} else {
		d.values = &DataomnisConfig{}
	}
	return d
}
