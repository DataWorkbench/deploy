package config


type Metrics struct {
	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

type GrpcLog struct {
	Level     int8 `json:"level,omitempty"     yaml:"level,omitempty"`
	Verbosity int8 `json:"verbosity,omitempty" yaml:"verbosity,omitempty"`
}

type Service struct {
	LogLevel  int8     `json:"logLevel,omitempty"  yaml:"logLevel"`
	LogOutput string   `json:"logOutput,omitempty" yaml:"logOutput"`
	GrpcLog   *GrpcLog `json:"grpcLog,omitempty"   yaml:"grpcLog,omitempty"`
	Metrics   *Metrics `json:"metrics,omitempty"   yaml:"metrics,omitempty"`

	Workload `json:",omitempty,inline" yaml:",omitempty,inline"`

	Envs map[string]string `json:"envs,omitempty" yaml:"envs,flow"`
}

type Webservice struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// TODO: update from webservice enable
// TODO: generate default region
type Apiglobal struct {
	Enabled    bool        `json:"enabled" yaml:"enabled" yaml:"enabled"`
	HttpServer *HttpServer `json:"httpServer,omitempty"   yaml:"httpServer,omitempty"`

	Regions      []Region               `json:"-"       yaml:"regions,flow,omitempty"`
	RegionsValue map[string]RegionValue `json:"regions" yaml:"-"`

	// Authentication: for helm values.yaml
	// IdentityProviders: for user configuration
	Authentication    *Authentication    `json:"authentication,omitempty" yaml:"-"`
	IdentityProviders []IdentityProvider `json:"-"                        yaml:"identityProviders,flow,omitempty"`
	HttpProxy         string             `json:"httpProxy,omitempty"      yaml:"httpProxy"`

	Service `json:",omitempty,inline" yaml:",inline"`
}

func (a *Apiglobal) UpdateRegion() {
	if len(a.Regions) > 0 {
		a.RegionsValue = map[string]RegionValue{}
		for _, r := range a.Regions {
			rv := RegionValue{
				Host: r.Host,
				Name: Names{
					ZhCn: r.ZhCn,
					EnUs: r.EnUs,
				},
			}
			a.RegionsValue[r.EnUs] = rv
		}
	}
}

func (a *Apiglobal) UpdateAuthentication() {
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
	HttpServer *HttpServer `json:"httpServer,omitempty"   yaml:"httpServer,omitempty"`
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
	Flink   *Flink      `json:"flink,omitempty"   yaml:"flink,omitempty"`
}

type HelmClient struct {
	RepoCachePath string `json:"repoCachePath,omitempty"   yaml:"RepoCachePath,omitempty"`
	Debug         bool   `json:"debug,omitempty"           yaml:"debug,omitempty"`
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
	Background    string `json:"background,omitempty"    yaml:"background,omitempty"    validate:"oneof=s3 hdfs"`
	HadoopConfDir string `json:"hadoopConfDir,omitempty" yaml:"hadoopConfDir,omitempty" validate:"required_if=Background hdfs"`
	S3            *S3    `json:"s3,omitempty"            yaml:"s3,omitempty"            validate:"required_if=Background s3"`
}

type S3 struct {
	Endpoint        string `json:"endpoint"        yaml:"endpoint"        validate:"required"`
	Region          string `json:"region"          yaml:"region"          validate:"required"`
	Bucket          string `json:"bucket"          yaml:"bucket"          validate:"required"`
	AccessKeyId     string `json:"accessKeyId"     yaml:"accessKeyId"     validate:"required"`
	SecretAccessKey string `json:"secretAccessKey" yaml:"secretAccessKey" validate:"required"`
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

// ***********************************************************
type IaasApiConfig struct {
	Zone            string `json:"zone"            yaml:"zone"            validate:"required"`
	Host            string `json:"host"            yaml:"host"            validate:"required"`
	Port            int    `json:"port"            yaml:"port"            validate:"required"`
	Protocol        string `json:"protocol"        yaml:"protocol"        validate:"required"`
	Timeout         int    `json:"timeout"         yaml:"timeout"         validate:"required"`
	Uri             string `json:"uri"             yaml:"uri"             validate:"required"`
	AccessKeyId     string `json:"accessKeyId"     yaml:"accessKeyId"     validate:"required"`
	SecretAccessKey string `json:"secretAccessKey" yaml:"secretAccessKey" validate:"required"`
}

type DataomnisConfig struct {
	// dataomnis version
	Version string `json:"-" yaml:"version"`

	Domain string `json:"domain"  yaml:"domain"`
	Port   string `json:"port"    yaml:"port,omitempty"`

	// global configurations for all service as default
	Image *Image `json:"image,omitempty" yaml:"image,omitempty"`

	MysqlClient *MysqlClient `json:"mysql" yaml:"mysql"`
	EtcdClient  *EtcdClient  `json:"etcd"  yaml:"-"`
	HdfsClient  *HdfsClient  `json:"hdfs"  yaml:"-"`
	RedisClient *RedisClient `json:"redis" yaml:"redisClient,omitempty"`

	Persistent *Persistent `json:"persistent" yaml:"-"`
	// all kube nodes from k8s apiserver to mkdir log-dir / helm-repo-cache
	Nodes []string `yaml:"nodes"`

	Iaas *IaasApiConfig `json:"iaas,omitempty" yaml:"iaas,omitempty" validate:"omitempty"`

	Common *Service `json:"common" yaml:"global"`

	WebService      *Webservice      `json:"webservice"                yaml:"webservice"`
	Apiglobal       *Apiglobal       `json:"apiglobal"                 yaml:"apiglobal"`
	Apiserver       *Apiserver       `json:"apiserver,omitempty"       yaml:"apiserver,omitempty"`
	Account         *Account         `json:"account,omitempty"         yaml:"account,omitempty"`
	Enginemanager   *Enginemanager   `json:"enginemanager"             yaml:"enginemanager"`
	Resourcemanager *Resourcemanager `json:"resourcemanager,omitempty" yaml:"resourcemanager,omitempty"`
	Scheduler       *Scheduler       `json:"scheduler,omitempty"       yaml:"scheduler,omitempty"`
	Spacemanager    *Service         `json:"spacemanager,omitempty"    yaml:"spacemanager,omitempty"`
	Developer       *Service         `json:"developer,omitempty"       yaml:"developer,omitempty"`

	Jaeger         *Workload `json:"jaeger,omitempty" yaml:"jaeger,omitempty"`
	ServiceMonitor *ServiceMonitor  `json:"serviceMonitor"   yaml:"serviceMonitor"`
}
