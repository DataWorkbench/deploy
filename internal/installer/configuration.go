package installer

import (
	"errors"
	"fmt"
	"github.com/DataWorkbench/glog"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var LeastNodeErr = errors.New("at least 3 nodes are required for helm release")

// ***************************************************************
// K8s Configuration Zone
// ***************************************************************

// firstAtAll, create docker registry secret by kubectl:
// kubectl create secret docker-registry my-docker-registry-secret
//                                       --docker-server=<your-registry-server>
//                                       --docker-username=<your-name>
//                                       --docker-password=<your-pword>
//                                       --docker-email=<your-email>
type ImageConfig struct {
	// TODO: check pull secrets
	Registry    string   `json:"registry,omitempty" yaml:"registry,omitempty"`
	PullSecrets []string `json:"pullSecrets,omitempty" yaml:"pullSecrets,omitempty"`
	PullPolicy  string   `json:"pullPolicy,omitempty" yaml:"pullPolicy,omitempty"`

	Tag string `json:",omitempty" yaml:"-"`
}

// update echo field-value from Config
func (i *ImageConfig) updateFromConfig(source *ImageConfig) {
	if source == nil {
		return
	}
	if i.Registry == "" && source.Registry != "" {
		i.Registry = source.Registry
	}
	if len(i.PullSecrets) == 0 && len(source.PullSecrets) > 0 {
		i.PullSecrets = source.PullSecrets
	}
	if i.PullPolicy == "" && source.PullPolicy != "" {
		i.PullPolicy = source.PullPolicy
	}
}

type Resource struct {
	Cpu    string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory string `json:"memory,omitempty" yaml:"memory,omitempty"`
}

type ResourceConfig struct {
	Limits   Resource `json:"limits,omitempty" yaml:"limits,omitempty"`
	Requests Resource `json:"requests,omitempty" yaml:"requests,omitempty"`
}

// k8s workload(a deployment / statefulset / .. in a Chart) configurations
type WorkloadConfig struct {
	Replicas       int8   `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	UpdateStrategy string `json:"updateStrategy,omitempty" yaml:"updateStrategy,omitempty"`

	Resource *ResourceConfig `json:"resources,omitempty" yaml:"resources,omitempty"`

	Persistent *PersistentConfig `json:"persistent,omitempty" yaml:"persistent,omitempty"`
}

type LocalPvConfig struct {
	Nodes []string `json:"nodes" yaml:"nodes"`
	Home  string   `json:"home" yaml:"-"`
}

type PersistentConfig struct {
	// for local pv, eg: 10Gi
	Size    string         `json:"size,omitempty" yaml:"size"`
	LocalPv *LocalPvConfig `json:"localPv,omitempty" yaml:"localPv"`
}

func (p *PersistentConfig) updateLocalPv(localPvHome string, nodes []string) error {
	// TODO: check if localPv exist and start with localPvHome
	p.LocalPv.Home = fmt.Sprintf(LocalHomeFmt, localPvHome)

	if len(p.LocalPv.Nodes) == 0 {
		if len(nodes) < 3 {
			return LeastNodeErr
		}
		p.LocalPv.Nodes = nodes
	}
	return nil
}

// ***************************************************************
// All Configuration From dataomnis-conf.yaml Zone
// ***************************************************************

// TODO: save the dataomnis-conf.yaml to k8s as configmap for backup
type Config struct {
	// kube nodes from k8s apiserver
	Nodes []string `yaml:"-"`

	// Local PV home
	LocalPVHome string `yaml:"localPvHome" validate:"required"`

	Image *ImageConfig `yaml:"image"`

	// dependent service
	Etcd  *EtcdConfig  `yaml:"etcdCluster"`
	Hdfs  *HdfsConfig  `yaml:"hdfsCluster"`
	Mysql *MysqlConfig `yaml:"mysqlCluster"`
	Redis *RedisConfig `yaml:"redisCluster"`

	// dataomnis version
	Dataomnis DataomnisConfig `yaml:"dataomnis"`
}

func (c *Config) Read(file string, logger glog.Logger) error {
	var err error
	// check configuration file
	_, err = os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("the configuration file: %s not exist!", file)
			logger.Error().Msg(msg).Fire()
			err = errors.New(msg)
			return err
		}
		err = nil
	}

	logger.Info().String("read configuration file", file).Fire()
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Error().String("failed to read configuration file", file).Error("error", err).Fire()
		logger.Error().Msg("please make sure the file is YAML format.").Fire()
		return err
	}
	logger.Info().Msg("parse content from configuration file to Config..").Fire()
	if err = yaml.Unmarshal(bytes, c); err != nil {
		logger.Error().Error("parse bytes from the configuration to yaml error", err).Fire()
		return err
	}

	logger.Debug().Any("Configuration", c).Fire()
	// validate
	logger.Info().Msg("validate Config..").Fire()
	if err = validator.New().Struct(c); err != nil {
		logger.Error().Error("validate configuration error", err).Fire()
		return err
	}
	return nil
}

func (c Config) checkImage() error {
	return nil
}
