package config

import (
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/pkg/errors"
)

type EtcdConfig struct {
	Image *Image `json:"image,omitempty" yaml:"image,omitempty"`

	Workload `json:",inline" yaml:",inline"`
}

func (e EtcdConfig) Validate(releaseName string) error {
	if len(e.Persistent.LocalPv.Nodes) != 3 {
		return errors.Errorf(common.OnlyNodeNumFmt, 3, releaseName)
	}
	return nil
}


// ***********************************************************
type EtcdClient struct {
	Endpoint string `json:"endpoint" yaml:"-"`
}
