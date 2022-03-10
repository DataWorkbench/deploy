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

	Image *ImageConfig `json:",omitempty"`
}

func (m ChartMeta) updateFromConfig(c Config) error {
	m.Image = c.Image
	return nil
}

func (m ChartMeta) parseValues() (Values, error) {
	var v Values = map[string]interface{}{}
	if m.Image == nil {
		return v, nil
	}

	bytes, err := json.Marshal(m.Image)
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

func (m ChartMeta) waitingReady() bool {
	return m.WaitingReady
}
