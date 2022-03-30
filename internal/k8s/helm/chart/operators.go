package chart

import (
	"github.com/DataWorkbench/deploy/internal/common"
)

// ***************************************************************
// ChartMeta Zone
// ***************************************************************
// implement Chart interface

type HdfsOperatorChart struct {
	ChartMeta
}

func NewHdfsOperatorChart(releaseName string) *HdfsOperatorChart {
	h := &HdfsOperatorChart{}
	h.ChartName = common.HdfsOptChart
	h.ReleaseName = releaseName
	h.Waiting = true
	return h
}


type MysqlOperatorChart struct {
	ChartMeta
}

func NewMysqlOperatorChart(releaseName string) *MysqlOperatorChart {
	m := &MysqlOperatorChart{}
	m.ChartName = common.MysqlOptChart
	m.ReleaseName = releaseName
	m.Waiting = true
	return m
}


type RedisOperatorChart struct {
	ChartMeta
}

func NewRedisOperatorChart(releaseName string) *RedisOperatorChart {
	r := &RedisOperatorChart{}
	r.ChartName = common.RedisOptChart
	r.ReleaseName = releaseName
	r.Waiting = true
	return r
}
