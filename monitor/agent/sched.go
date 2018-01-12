package main

/*
	这是一个调度模块，将本端的数据进行采集，然后将数据交个网络模块，最后由网络模块将数据传输给server端.
*/

import (
	"time"
	"yinzhengjie/monitor/common"
)

type MetricFunc func() []*common.Metric //可以返回多个函数指标（你可以理解是监控指标。）

type Sched struct {
	ch chan *common.Metric
}

func NewSched(ch chan *common.Metric) *Sched {
	return &Sched{
		ch: ch,
	}
}

func (s *Sched) AddMetric(collecter MetricFunc, step time.Duration) { //该函数的功能是间隔“step”时间周期会调用一次“collecter”函数
	ticker := time.NewTicker(step)
	for range ticker.C {
		metrics := collecter()
		for _, metric := range metrics {
			s.ch <- metric
		}
	}
}
