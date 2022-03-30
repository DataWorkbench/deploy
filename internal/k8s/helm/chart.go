package helm

import (
	"encoding/json"
	"github.com/DataWorkbench/deploy/internal/common"
	"gopkg.in/yaml.v3"
	"time"
)

// **************************************************
// the type Values used for helm when install release
// **************************************************
type Values map[string]interface{}

func (v Values) parse() (string, error) {
	valueBytes, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(valueBytes), nil
}

func (v Values) isEmpty() bool {
	return len(v) == 0
}

// ***************************************************************
// ChartMeta Zone
// ***************************************************************
// implement Chart interface
type Meta struct {
	Image *common.Image `json:"image,omitempty"`
}

type ChartMeta struct {
	// for pod label
	ChartName   string
	ReleaseName string

	Waiting bool

	Conf *Meta `json:",omitempty"`
}

func (m *ChartMeta) UpdateFromConfig(c common.Config) error {
	if c.Image == nil {
		return nil
	}

	if m.Conf == nil {
		m.Conf = &Meta{}
	}
	if m.Conf.Image == nil {
		m.Conf.Image = &common.Image{}
	}
	m.Conf.Image.Copy(c.Image)
	return nil
}

func (m ChartMeta) InitLocalDir() error {
	return nil
}

func (m ChartMeta) ParseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	if m.Conf == nil || m.Conf.Image == nil {
		return v, nil
	}

	bytes, err := json.Marshal(m.Conf)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}

func (m ChartMeta) GetChartName() string {
	return m.ChartName
}

func (m ChartMeta) GetReleaseName() string {
	return m.ReleaseName
}

func (m ChartMeta) GetLabels() map[string]string {
	return map[string]string{
		common.InstanceLabelKey: m.ReleaseName,
	}
}

func (m ChartMeta) WaitingReady() bool {
	return m.Waiting
}

func (m ChartMeta) GetTimeoutSecond() time.Duration {
	return common.DefaultTimeoutSecond * time.Second
}
