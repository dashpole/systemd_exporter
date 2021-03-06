package prometheus

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/coreos/go-systemd/dbus"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"github.com/opencontainers/runc/libcontainer/configs"

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
	unitCPUUsage = metrics.NewDesc(
		"unit_cpu_usage_seconds_total",
		"Systemd unit CPU usage in seconds",
		[]string{"name"},
		nil,
		metrics.ALPHA,
		"",
	)
	unitMemoryUsage = metrics.NewDesc(
		"unit_memory_usage_bytes",
		"Systemd unit memory usage in bytes",
		[]string{"name"},
		nil,
		metrics.ALPHA,
		"",
	)
	allowedUnits = map[string]bool{
		"kubelet":               true,
		"containerd":            true,
		"dockerd":               true,
		"node-problem-detector": true,
		"systemd-journald":      true,
	}
)

type systemd struct {
	metrics.BaseStableCollector
	mountPoints map[string]string
}

// NewSystemdCollector returns a prometheus collector for systemd
func NewSystemdCollector() (metrics.StableCollector, error) {
	allCgroups, err := cgroups.GetCgroupMounts(true)
	if err != nil {
		return nil, fmt.Errorf("Failed to get cgroup mounts: %v", err)
	}
	allMountPoints := map[string]string{}
	for _, mount := range allCgroups {
		for _, subsystem := range mount.Subsystems {
			allMountPoints[subsystem] = mount.Mountpoint
		}
	}
	klog.Infof("Initialized with mountpoints: %+v", allMountPoints)
	return &systemd{mountPoints: allMountPoints}, nil
}

func (s *systemd) DescribeWithStability(ch chan<- *metrics.Desc) {
	ch <- unitState
	ch <- unitCPUUsage
	ch <- unitMemoryUsage
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
		if strings.HasSuffix(unit.Name, ".service") {
			if _, ok := allowedUnits[strings.TrimSuffix(unit.Name, ".service")]; ok {
				s.collectUnitState(conn, unit, ch)
				s.collectUnitCgroupMetrics(conn, unit, ch)
			}
		}
	}
}

func (s *systemd) collectUnitCgroupMetrics(conn *dbus.Conn, unit dbus.UnitStatus, ch chan<- metrics.Metric) {
	sliceProperties, err := conn.GetUnitTypeProperties(unit.Name, "Service")
	if err != nil {
		klog.Warningf("Failed to get unit service for unit %v. No metrics will be collected: %v", unit.Name, err)
		return
	}
	cgroup, found := sliceProperties["ControlGroup"]
	if !found {
		klog.Warningf("ControlGroup property for unit %v not found. No metrics will be collected.", unit.Name)
		return
	}
	cgroupName, ok := cgroup.(string)
	if !ok {
		klog.Warningf("Failed to convert cgroup: %v to string.", cgroup)
		return
	}
	cgroupPaths := make(map[string]string, len(s.mountPoints))
	for k, v := range s.mountPoints {
		cgroupPaths[k] = path.Join(v, cgroupName)
	}

	manager := fs.NewManager(&configs.Cgroup{Name: cgroupName}, cgroupPaths, false)
	stats, err := manager.GetStats()
	if err != nil {
		klog.Warningf("Failed to get stats for cgroup %v: %v", cgroupName, err)
		return
	}
	if stats.CpuStats.CpuUsage.TotalUsage > 0 {
		ch <- metrics.NewLazyConstMetric(
			unitCPUUsage, metrics.CounterValue,
			float64(stats.CpuStats.CpuUsage.TotalUsage)/float64(time.Second),
			strings.TrimSuffix(unit.Name, ".service"))
	}
	ch <- metrics.NewLazyConstMetric(
		unitMemoryUsage, metrics.GaugeValue,
		float64(stats.MemoryStats.Usage.Usage),
		strings.TrimSuffix(unit.Name, ".service"))
}

func (s *systemd) collectUnitState(conn *dbus.Conn, unit dbus.UnitStatus, ch chan<- metrics.Metric) {
	serviceType := ""
	serviceTypeProperty, err := conn.GetUnitTypeProperty(unit.Name, "Service", "Type")
	if err != nil {
		klog.Warningf("Failed to get unit type for unit %v: %v", unit.Name, err)
		return
	}
	serviceType = serviceTypeProperty.Value.Value().(string)
	name := strings.TrimSuffix(unit.Name, ".service")
	for _, stateName := range unitStates {
		isActive := 0.0
		if stateName == unit.ActiveState {
			isActive = 1.0
		}
		ch <- metrics.NewLazyConstMetric(
			unitState, metrics.GaugeValue, isActive,
			name, stateName, serviceType)
	}
}
