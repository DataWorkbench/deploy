package config

import (
	"fmt"
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/pkg/errors"
)

//TODO: add backup-config
type MysqlConfig struct {
	TimeoutSecond int       `json:"-" yaml:"timeoutSecond,omitempty"`
	Image         *Image    `json:"image,omitempty" yaml:"image,omitempty"`
	Pxc           *Workload `json:"pxc" yaml:"pxc"`
}

func (m MysqlConfig) Validate(releaseName string) error {
	if len(m.Pxc.Persistent.LocalPv.Nodes) != 3 {
		return errors.Errorf(common.OnlyNodeNumFmt, 3, releaseName)
	}
	return nil
}

// ***********************************************************
type MysqlClient struct {
	ExternalHost string `json:"externalHost" yaml:"-"`
	SecretName   string `json:"secretName" yaml:"-"`

	LogLevel        int8   `json:"logLevel,omitempty"        yaml:"logLevel,omitempty"`
	MaxIdleConn     int32  `json:"maxIdleConn,omitempty"     yaml:"maxIdleConn,omitempty"`
	MaxOpenConn     int32  `json:"maxOpenConn,omitempty"     yaml:"maxOpenConn,omitempty"`
	ConnMaxLifetime string `json:"connMaxLifetime,omitempty" yaml:"connMaxLifetime,omitempty"`
	SlowThreshold   string `json:"slowThreshold,omitempty"   yaml:"slowThreshold,omitempty"`
}

func (c *MysqlClient) Update(releaseName string) {
	c.ExternalHost = fmt.Sprintf(common.MysqlExternalHostFmt, releaseName)
	c.SecretName = fmt.Sprintf(common.MysqlSecretNameFmt, releaseName)
}
