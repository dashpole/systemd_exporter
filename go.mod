module github.com/dashpole/systemd_exporter

go 1.13

require (
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/coreos/go-systemd/v22 v22.1.0 // indirect
	github.com/godbus/dbus v0.0.0-20190402143921-271e53dc4968 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/opencontainers/runc v1.0.0-rc90.0.20200616040943-82d2fa4eb069
	github.com/prometheus/client_golang v1.7.1
	google.golang.org/protobuf v1.24.0 // indirect
	k8s.io/component-base v0.18.5
	k8s.io/klog v1.0.0
)
