package main

import (
	"fmt"
	"net/http"

	"github.com/dashpole/systemd_exporter/pkg/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"k8s.io/component-base/metrics"
	"k8s.io/klog"
)

var (
	port     = 8080
	endpoint = "/metrics"
)

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	mux := http.NewServeMux()
	reg := metrics.NewKubeRegistry()

	collector, err := prometheus.NewSystemdCollector()
	if err != nil {
		klog.Fatalf("Failed to create Systemd collector: %v", err)
	}

	reg.CustomMustRegister(collector)

	mux.Handle(endpoint, promhttp.HandlerFor(reg, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))

	klog.Infof("Starting systemd exporter")

	klog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
