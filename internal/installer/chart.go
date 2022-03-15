package installer

import "encoding/json"

// ***************************************************************
// ChartMeta Zone
// ***************************************************************
// implement Chart interface
type ChartMeta struct {
	// for pod label
	ChartName   string
	ReleaseName string

	WaitingReady bool

	values *MetaConfig `json:",omitempty"`
}

func (m *ChartMeta) updateFromConfig(c Config) error {
	if c.Image == nil {
		return nil
	}

	if m.values == nil {
		m.values = &MetaConfig{}
	}
	if m.values.Image == nil {
		m.values.Image = &ImageConfig{}
	}
	m.values.Image.updateFromConfig(c.Image)
	return nil
}

func (m ChartMeta) parseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	if m.values == nil || m.values.Image == nil {
		return v, nil
	}

	bytes, err := json.Marshal(m.values)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bytes, &v)
	return v, err
}

func (m ChartMeta) getChartName() string {
	return m.ChartName
}

func (m ChartMeta) getReleaseName() string {
	return m.ReleaseName
}

func (m ChartMeta) getLabels() map[string]string {
	return map[string]string{
		InstanceLabelKey: m.ReleaseName,
	}
}

func (m ChartMeta) waitingReady() bool {
	return m.WaitingReady
}

type MetaConfig struct {
	Image *ImageConfig `json:"image,omitempty"`
}