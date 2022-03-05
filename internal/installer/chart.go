package installer

import "encoding/json"

// ***************************************************************
// ChartMeta Zone
// ***************************************************************
// implement Chart interface
type ChartMeta struct {
	ChartName   string
	ReleaseName string

	WaitingReady bool

	values ValuesConfig
}

func NewChartMeta(chartName, releaseName string, waittingReady bool) *ChartMeta {
	return &ChartMeta{
		ChartName:    chartName,
		ReleaseName:  releaseName,
		WaitingReady: waittingReady,
	}
}

func (m *ChartMeta) setMeta(chartName, releaseName string, waittingReady bool) {
	m.ChartName = chartName
	m.ReleaseName = releaseName
	m.WaitingReady = waittingReady
}

func (m ChartMeta) updateFromConfig(config Config) error {
	return nil
}

func (m ChartMeta) parseValues() (Values, error) {
	var v Values = map[string]interface{}{}
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

func (m ChartMeta) waitingReady() bool {
	return m.WaitingReady
}
