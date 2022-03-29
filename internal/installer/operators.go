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
	h.WaitingReady = true

	_ = h.updateFromConfig(c)
	return h
}


type MysqlOperatorChart struct {
	helm.ChartMeta
}

func NewMysqlOperatorChart(releaseName string, c common.Config) *MysqlOperatorChart {
	m := &MysqlOperatorChart{}
	m.ChartName = common.MysqlOptChart
	m.ReleaseName = releaseName
	m.WaitingReady = true
	_ = m.updateFromConfig(c)
	return m
}


type RedisOperatorChart struct {
	helm.ChartMeta
}

func NewRedisOperatorChart(releaseName string, c common.Config) *RedisOperatorChart {
	r := &RedisOperatorChart{}
	r.ChartName = common.RedisOptChart
	r.ReleaseName = releaseName
	r.WaitingReady = true
	_ = r.updateFromConfig(c)
	return r
}
