module github.com/dashpole/systemd_exporter

go 1.13

require (
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/godbus/dbus v0.0.0-20190402143921-271e53dc4968 // indirect
	github.com/prometheus/client_golang v1.7.1
	k8s.io/component-base v0.18.5
	k8s.io/klog v1.0.0
)
