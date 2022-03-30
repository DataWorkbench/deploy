package installer

import (
	"github.com/DataWorkbench/deploy/internal/common"
	"github.com/DataWorkbench/deploy/internal/k8s/helm"
)

// ***************************************************************
// ChartMeta Zone
// ***************************************************************
// implement Chart interface

type HdfsOperatorChart struct {
	helm.ChartMeta
}

func NewHdfsOperatorChart(releaseName string, c common.Config) *HdfsOperatorChart {
	h := &HdfsOperatorChart{}
	h.ChartName = common.HdfsOptChart
	h.ReleaseName = releaseName
	h.Waiting = true

	_ = h.UpdateFromConfig(c)
	return h
}


type MysqlOperatorChart struct {
	helm.ChartMeta
}

func NewMysqlOperatorChart(releaseName string, c common.Config) *MysqlOperatorChart {
	m := &MysqlOperatorChart{}
	m.ChartName = common.MysqlOptChart
	m.ReleaseName = releaseName
	m.Waiting = true
	_ = m.UpdateFromConfig(c)
	return m
}


type RedisOperatorChart struct {
	helm.ChartMeta
}

func NewRedisOperatorChart(releaseName string, c common.Config) *RedisOperatorChart {
	r := &RedisOperatorChart{}
	r.ChartName = common.RedisOptChart
	r.ReleaseName = releaseName
	r.Waiting = true
	_ = r.UpdateFromConfig(c)
	return r
}
