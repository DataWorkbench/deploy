package installer

// ***************************************************************
// ChartMeta Zone
// ***************************************************************
// implement Chart interface

type HdfsOperatorChart struct {
	ChartMeta
}

func NewHdfsOperatorChart(releaseName string, c Config) *HdfsOperatorChart {
	h := &HdfsOperatorChart{}
	h.ChartName = HdfsOptChart
	h.ReleaseName = releaseName
	h.WaitingReady = true

	_ = h.updateFromConfig(c)
	return h
}


type MysqlOperatorChart struct {
	ChartMeta
}

func NewMysqlOperatorChart(releaseName string, c Config) *MysqlOperatorChart {
	m := &MysqlOperatorChart{}
	m.ChartName = MysqlOptChart
	m.ReleaseName = releaseName
	m.WaitingReady = true
	_ = m.updateFromConfig(c)
	return m
}


type RedisOperatorChart struct {
	ChartMeta
}

func NewRedisOperatorChart(releaseName string, c Config) *RedisOperatorChart {
	r := &RedisOperatorChart{}
	r.ChartName = RedisOptChart
	r.ReleaseName = releaseName
	r.WaitingReady = true
	_ = r.updateFromConfig(c)
	return r
}
