package prometheus

import (
	"k8s.io/component-base/metrics"
)

type systemd struct {
	metrics.BaseStableCollector
}

func NewSystemdCollector() metrics.StableCollector {
	return &systemd{}
}

func (s *systemd) DescribeWithStability(ch chan<- *metrics.Desc) {}

func (s *systemd) CollectWithStability(ch chan<- metrics.Metric) {}
