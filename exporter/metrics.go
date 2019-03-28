package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
)

func AddMetrics(service_fqdn string, pid_enabled bool, pid_fqdn string) map[string]*prometheus.Desc {
	var prometheus_svc_namespace = service_fqdn
	if prometheus_svc_namespace == "" {
		prometheus_svc_namespace = prometheus.BuildFQName("common", "service", "state")
	}

	ServiceMetrics := make(map[string]*prometheus.Desc)

	ServiceMetrics["State"] = prometheus.NewDesc(
		prometheus_svc_namespace,
		"The current state of the service",
		[]string{"service", "state", "group"},
		nil,
		)

	if pid_enabled {
		var prometheus_pid_namespace = pid_fqdn
		if prometheus_pid_namespace == "" {
			prometheus_pid_namespace = prometheus.BuildFQName("common", "service", "pid")
		}

		ServiceMetrics["PID"] = prometheus.NewDesc(
			prometheus_pid_namespace,
			"The current value of the pid file",
			[]string{"filename", "state", "service", "group"},
			nil,
			)

	}

	return ServiceMetrics
}

func (e *Exporter) processMetrics(services []*Service, pids []*PidFile, ch chan<- prometheus.Metric) error {
	for _, x := range services {
		ch <- prometheus.MustNewConstMetric(e.ServiceMetrics["State"], prometheus.GaugeValue, x.IsActive(), x.Name, x.Substate, x.Group)
	}

	for _, y := range pids {
		parsed_value, err := strconv.ParseFloat(strings.Replace(y.PID, "\n", "", 1), 64)
		if err == nil {
		    ch <- prometheus.MustNewConstMetric(e.ServiceMetrics["PID"], prometheus.GaugeValue, parsed_value, y.Name, y.State, y.Service, y.Group)
		}
	}

	return nil
}
