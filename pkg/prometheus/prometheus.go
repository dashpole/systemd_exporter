package prometheus

import (
	"fmt"

	"github.com/coreos/go-systemd/dbus"

	"k8s.io/component-base/metrics"
	"k8s.io/klog"
)

type systemd struct {
	metrics.BaseStableCollector

	units []dbus.UnitStatus
}

func NewSystemdCollector() metrics.StableCollector {
	return &systemd{}
}

func (s *systemd) UpdateUnits() error {
	conn, err := dbus.New()
	if err != nil {
		return fmt.Errorf("failed to get dbus connection: %v", err)
	}
	defer conn.Close()
	units, err := conn.ListUnits()
	if err != nil {
		return fmt.Errorf("failed to list units: %v", err)
	}
	s.units = units
	klog.Infof("updated units %+v", units)
	return nil
}

func (s *systemd) DescribeWithStability(ch chan<- *metrics.Desc) {
	if err := s.UpdateUnits(); err != nil {
		klog.Errorf("Error updating units: %+v", err)
	}
}

func (s *systemd) CollectWithStability(ch chan<- metrics.Metric) {
}
