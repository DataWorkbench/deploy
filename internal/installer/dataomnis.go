package installer

import (
	"encoding/json"
	"fmt"
	"github.com/DataWorkbench/common/lib/iaas"
	"github.com/pkg/errors"
	"strings"
)

type Metrics struct {
	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

type GrpcLog struct {
	Level     int8 `json:"level,omitempty"     yaml:"level,omitempty"`
	Verbosity int8 `json:"verbosity,omitempty" yaml:"verbosity,omitempty"`
}

type Service struct {
	LogLevel  int8     `json:"logLevel,omitempty" yaml:"logLevel"`
	LogOutput string   `json:"logOutput"          yaml:"logOutput"`
	GrpcLog   *GrpcLog `json:"grpcLog,omitempty"  yaml:"grpcLog,omitempty"`
	Metrics   *Metrics `json:"metrics,omitempty"  yaml:"metrics,omitempty"`

	WorkloadConfig `json:",omitempty,inline" yaml:",omitempty,inline"`

	Envs map[string]string `json:"envs,omitempty" yaml:"envs,flow"`
}

type Webservice struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// TODO: update from webservice enable
// TODO: generate default region
type Apiglobal struct {
	Enabled    bool       `json:"enabled" yaml:"enabled" yaml:"enabled"`
	HttpServer HttpServer `json:"httpServer,omitempty"   yaml:"httpServer,omitempty"`

	Regions      []Region                 `json:"-"       yaml:"regions,flow,omitempty"`
	RegionValues []map[string]RegionValue `json:"regions" yaml:"-"`

	// Authentication: for helm values.yaml
	// IdentityProviders: for user configuration
	Authentication    *Authentication    `json:"authentication"      yaml:"-"`
	IdentityProviders []IdentityProvider `json:"-"                   yaml:"identityProviders"`
	HttpProxy         string             `json:"httpProxy,omitempty" yaml:"httpProxy"`

	Service `json:",omitempty,inline" yaml:",inline"`
}

func (a *Apiglobal) updateRegion() {
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

func (a *Apiglobal) updateAuthentication() {
	if len(a.IdentityProviders) > 0 {
		pMap := map[string]IdentityProvider{}
		for _, p := range a.IdentityProviders {
			pMap[p.Name] = p
		}
		a.Authentication.IdentityProviders = pMap
	}
}

// HttpServer config
type HttpServer struct {
	ReadTimeout  string `json:"read_timeout,omitempty" yaml:"readTimeout,omitempty"`
	WriteTimeout string `json:"write_timeout,omitempty" yaml:"writeTimeout,omitempty"`
	IdleTimeout  string `json:"idle_timeout,omitempty" yaml:"idleTimeout,omitempty"`
	ExitTimeout  string `json:"exit_timeout,omitempty" yaml:"exitTimeout,omitempty"`
}

// Region for configuration by user, parsed to RegionValue.
type Region struct {
	Host string `json:"-" yaml:"host"`
	EnUs string `json:"-" yaml:"enUsName"`
	ZhCn string `json:"-" yaml:"zhCnName"`
}

// RegionValue for helm values.yaml, updated from Region.
type Names struct {
	ZhCn string `json:"zh_cn" yaml:"-"`
	EnUs string `json:"en_us" yaml:"-"`
}

type RegionValue struct {
	Host string `json:"hosts" yaml:"-"`
	Name Names  `json:"names"`
}

// IdentityProvider for authentication
type IdentityProvider struct {
	Name         string `json:"name"         yaml:"name"`
	ClientId     string `json:"clientId"     yaml:"clientId"`
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`
	TokenUrl     string `json:"tokenUrl"     yaml:"tokenUrl"`
	RedirectUrl  string `json:"redirectUrl"  yaml:"redirectUrl"`
}

type Authentication struct {
	IdentityProviders map[string]IdentityProvider `json:"identityProviders,omitempty" yaml:"-"`
}

// ***********************************************************
type Apiserver struct {
	Service    `json:",omitempty,inline" yaml:",inline"`
	HttpServer HttpServer `json:"httpServer,omitempty"   yaml:"httpServer,omitempty"`
}

// ***********************************************************
type Account struct {
	Service `json:",omitempty,inline" yaml:",inline"`
	Source  string `json:"source,omitempty"   yaml:"source,omitempty"`
}

// ***********************************************************
type Enginemanager struct {
	Service `json:",omitempty,inline" yaml:",inline"`
	Helm    *HelmClient `json:"helm,omitempty"   yaml:"helm,omitempty"`
	Flink   Flink       `json:"flink,omitempty"   yaml:"flink,omitempty"`
}

type HelmClient struct {
	RepoCachePath string `json:"repoCachePath,omitempty"   yaml:"RepoCachePath,omitempty"`
	Debug         bool   `json:"debug,omitempty"   yaml:"debug,omitempty"`
}

type Flink struct {
	RestServicePort    int8   `json:"restServicePort,omitempty"   yaml:"restServicePort,omitempty"`
	RestServiceNameFmt string `json:"restServiceNameFmt,omitempty"   yaml:"restServiceNameFmt,omitempty"`
	IngressClass       string `json:"ingressClass,omitempty"   yaml:"ingressClass,omitempty"`
	EnableMultus       bool   `json:"enableMultus,omitempty"   yaml:"enableMultus,omitempty"`
}

// ***********************************************************
type Resourcemanager struct {
	Service `json:",omitempty,inline" yaml:",inline"`
	Storage *Storage `json:"storage,omitempty"   yaml:"storage,omitempty"`
}

type Storage struct {
	Background    string `json:"background,omitempty"    yaml:"background,omitempty"`
	HadoopConfDir string `json:"hadoopConfDir,omitempty" yaml:"hadoopConfDir,omitempty"`
	S3            *S3    `json:"s3,omitempty"            yaml:"s3,omitempty"`
}

type S3 struct {
	Endpoint string `json:"endpoint,omitempty"   yaml:"endpoint,omitempty"`
	Region   string `json:"region,omitempty"   yaml:"region,omitempty"`
	Bucket   string `json:"bucket,omitempty"   yaml:"bucket,omitempty"`
}

// ***********************************************************
type Scheduler struct {
	Service         `json:",omitempty,inline" yaml:",inline"`
	EtcdDialTimeout string `json:"etcdDialTimeout,omitempty"   yaml:"etcdDialTimeout,omitempty"`
}

// ***********************************************************
type ServiceMonitor struct {
	Enabled  bool   `json:"enabled"  yaml:"enabled"`
	Interval string `json:"interval" yaml:"interval"`
}

type MysqlClient struct {
	ExternalHost string `json:"externalHost" yaml:"-"`
	SecretName   string `json:"secretName" yaml:"-"`

	LogLevel        int8   `json:"logLevel,omitempty"        yaml:"logLevel,omitempty"`
	MaxIdleConn     int32  `json:"maxIdleConn,omitempty"     yaml:"maxIdleConn,omitempty"`
	MaxOpenConn     int32  `json:"maxOpenConn,omitempty"     yaml:"maxOpenConn,omitempty"`
	ConnMaxLifetime string `json:"connMaxLifetime,omitempty" yaml:"connMaxLifetime,omitempty"`
	SlowThreshold   string `json:"slowThreshold,omitempty"   yaml:"slowThreshold,omitempty"`
}

type EtcdClient struct {
	Endpoint string `json:"endpoint" yaml:"-"`
}

type HdfsClient struct {
	ConfigmapName string `json:"configmapName" yaml:"-"`
}

type RedisClient struct {
	Mode         string `json:"mode,omitempty"         yaml:"mode,omitempty"`
	SentinelAddr string `json:"sentinelAddr,omitempty" yaml:"sentinelAddr,omitempty"`
	ClusterAddr  string `json:"clusterAddr,omitempty"  yaml:"clusterAddr,omitempty"`
	MasterName   string `json:"masterName,omitempty"   yaml:"masterName,omitempty"`
	Database     string `json:"database,omitempty"     yaml:"database,omitempty"`
	Username     string `json:"username,omitempty"     yaml:"username,omitempty"`
	Password     string `json:"password,omitempty"     yaml:"password,omitempty"`
}

func (r *RedisClient) generateAddr(releaseName string, size int) {
	var addrs []string
	if r.Mode == "cluster" && r.ClusterAddr == "" {
		if releaseName == RedisClusterChartName {
			for i := 0; i < size; i++ {
				addrs = append(addrs, fmt.Sprintf(RedisClusterAddrFmt, releaseName, i, RedisClusterPort))
			}
		}
		r.ClusterAddr = strings.Join(addrs, ",")
	}
}

type Dataomnis struct {
	// dataomnis version
	Version string `json:"version" yaml:"version"`

	Domain string `json:"domain"  yaml:"domain"`
	Port   string `json:"port"    yaml:"port,omitempty"`

	// global configurations for all service as default
	Image *ImageConfig `json:"image,omitempty" yaml:"image,omitempty"`

	MysqlClient *MysqlClient `json:"mysql" yaml:"mysql"`
	EtcdClient  *EtcdClient  `json:"etcd"  yaml:"-"`
	HdfsClient  *HdfsClient  `json:"hdfs"  yaml:"-"`
	RedisClient *RedisClient `json:"redis" yaml:"redisCluster,omitempty"`

	WorkloadConfig `json:",inline" yaml:"-"`

	Iaas *iaas.Config `json:"iaas,omitempty" yaml:"iaas,omitempty" validate:"omitempty"`

	Common *Service `json:"common" yaml:"global"`

	WebService      *Webservice      `json:"webservice"      yaml:"webservice"`
	Apiglobal       *Apiglobal       `json:"apiglobal"       yaml:"apiGlobal"`
	Apiserver       *Apiserver       `json:"apiserver"       yaml:"apiserver"`
	Account         *Account         `json:"account"         yaml:"account"`
	Enginemanager   *Enginemanager   `json:"enginemanager"   yaml:"enginemanager"`
	Resourcemanager *Resourcemanager `json:"resourcemanager" yaml:"resourcemanager"`
	Scheduler       *Scheduler       `json:"scheduler"       yaml:"scheduler"`
	Spacemanager    *Service         `json:"spacemanager"    yaml:"spacemanager"`
	Developer       *Service         `json:"developer"       yaml:"developer"`

	Jaeger         *WorkloadConfig `json:"jaeger"         yaml:"jaeger"`
	ServiceMonitor *ServiceMonitor `json:"serviceMonitor" yaml:"serviceMonitor"`
}

type DataomnisChart struct {
	ChartMeta

	values *Dataomnis
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
	d.values.RedisClient.generateAddr(RedisClusterName, 3)

	// update hostPath for log-dir
	d.values.Persistent.HostPath = fmt.Sprintf(DataomnisHostPathFmt, c.LocalPVHome, d.getReleaseName())
	d.values.Persistent.LocalPv = nil

	d.values.Apiglobal.updateRegion()
	d.values.Apiglobal.updateAuthentication()
	return nil
}

func (d *DataomnisChart) initHostPathDir(c Config) error {
	localPvHome := fmt.Sprintf("%s/log/{account,apiglobal,apiserver,enginemanager,resourcemanager,scheduler,spacemanager}", d.values.Persistent.HostPath)
	var host *Host
	var conn *Connection
	var err error
	for _, node := range c.Nodes {
		host = &Host{Address: node}
		conn, err = NewConnection(host)
		if err != nil {
			return errors.Wrap(err, "new connection failed")
		}
		if err := conn.Mkdir(localPvHome); err != nil {
			return err
		}
	}
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
		d.values = &Dataomnis{}
	}
	return d
}
