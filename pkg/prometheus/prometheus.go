package prometheus

import (
	"strings"

	"github.com/coreos/go-systemd/dbus"

	"k8s.io/component-base/metrics"
	"k8s.io/klog"
)

var (
	unitStates = []string{"active", "activating", "deactivating", "inactive", "failed"}
	unitState  = metrics.NewDesc(
		"unit_state",
		"Systemd unit",
		[]string{"name", "state", "type"},
		nil,
		metrics.ALPHA,
		"",
	)
)

type systemd struct {
	metrics.BaseStableCollector
}

func NewSystemdCollector() metrics.StableCollector {
	return &systemd{}
}

func (s *systemd) DescribeWithStability(ch chan<- *metrics.Desc) {
	ch <- unitState
}

func (s *systemd) CollectWithStability(ch chan<- metrics.Metric) {
	conn, err := dbus.New()
	if err != nil {
		klog.Errorf("failed to get dbus connection: %v", err)
		return
	}
	defer conn.Close()
	units, err := conn.ListUnits()
	if err != nil {
		klog.Errorf("failed to list units: %v", err)
		return
	}
	for _, unit := range units {
		serviceType := ""
		if !strings.HasSuffix(unit.Name, ".service") {
			continue
		}
		serviceTypeProperty, err := conn.GetUnitTypeProperty(unit.Name, "Service", "Type")
		if err != nil {
			klog.Errorf("Failed to get unit type for unit %v: %v", unit.Name, err)
		} else {
			serviceType = serviceTypeProperty.Value.Value().(string)
		}
		for _, stateName := range unitStates {
			isActive := 0.0
			if stateName == unit.ActiveState {
				isActive = 1.0
			}
			ch <- metrics.NewLazyConstMetric(
				unitState, metrics.GaugeValue, isActive,
				unit.Name, stateName, serviceType)
		}
	}
}
