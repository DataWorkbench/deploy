package installer

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type EtcdConfig struct {
	Image *Image `json:"image,omitempty" yaml:"image,omitempty"`

	Workload `json:",inline" yaml:",inline"`
}

// EtcdChart for etcd-cluster, implement Chart
type EtcdChart struct {
	ChartMeta `json:",inline" yaml:",inline"`
	values    *EtcdConfig `yaml:"config,omitempty"`
}

// update each field value from global Config if that is ZERO
func (e *EtcdChart) updateFromConfig(c Config) error {
	if c.Etcd != nil {
		e.values = c.Etcd
	}

	if c.Image != nil {
		if e.values.Image == nil {
			e.values.Image = &Image{}
			e.values.Image.updateFromConfig(c.Image)
		}
	}

	return e.values.Persistent.updateLocalPv(c.LocalPVHome, c.Nodes)
}

func (e *EtcdChart) initLocalPvHome() error {
	localPvHome := fmt.Sprintf("%s/%s", e.values.Persistent.LocalPv.Home, EtcdClusterName)
	var host *Host
	var conn *Connection
	var err error
	for _, node := range e.values.Persistent.LocalPv.Nodes {
		host = &Host{Address: node}
		conn, err = NewConnection(host)
		if err != nil {
			return errors.Wrap(err, "new connection failed")
		}
		if _, err := conn.Mkdir(localPvHome); err != nil {
			return err
		}
	}
	return nil
}

func (e *EtcdChart) parseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	bytes, err := json.Marshal(e.values)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}

func (e *EtcdChart) getTimeoutSecond() int {
	if e.values.TimeoutSecond == 0 {
		return e.ChartMeta.getTimeoutSecond()
	}
	return e.values.TimeoutSecond
}

func NewEtcdChart(release string, c Config) *EtcdChart {
	e := &EtcdChart{}
	e.ChartName = EtcdClusterChart
	e.ReleaseName = release
	e.WaitingReady = true

	if c.Etcd != nil {
		e.values = c.Etcd
	} else {
		e.values = &EtcdConfig{}
	}
	return e
}
