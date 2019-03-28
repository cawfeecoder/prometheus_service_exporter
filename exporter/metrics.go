package exporter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
)

func BuildFQDN(fqdn string) string {
	if fqdn != "" {
		parts := strings.Split(fqdn, "_")
		if len(parts) == 3 {
			return prometheus.BuildFQName(parts[0], parts[1], parts[2])
		} else {
			return ""
		}
	}
	return ""
}

func AddMetrics(service_fqdn string, pid_enabled bool, pid_fqdn string) map[string]*prometheus.Desc {
	var prometheus_svc_namespace = BuildFQDN(service_fqdn)
	if prometheus_svc_namespace == "" {
		prometheus_svc_namespace = prometheus.BuildFQName("common", "service", "state")
	}

	ServiceMetrics := make(map[string]*prometheus.Desc)

	ServiceMetrics["State"] = prometheus.NewDesc(
		prometheus_svc_namespace,
		"The current state of the service",
		[]string{"service", "state"},
		nil,
		)

	if pid_enabled {
		var prometheus_pid_namespace = BuildFQDN(pid_fqdn)
		if prometheus_pid_namespace == "" {
			prometheus_pid_namespace = prometheus.BuildFQName("common", "service", "pid")
		}

		ServiceMetrics["PID"] = prometheus.NewDesc(
			prometheus_pid_namespace,
			"The current value of the pid file",
			[]string{"filename", "state"},
			nil,
			)

	}

	return ServiceMetrics
}

func (e *Exporter) processMetrics(services []*Service, pids []*PidFile, ch chan<- prometheus.Metric) error {
	for _, x := range services {
		ch <- prometheus.MustNewConstMetric(e.ServiceMetrics["State"], prometheus.GaugeValue, x.IsActive(), x.Name, x.Substate)
	}

	for _, y := range pids {
		parsed_value, err := strconv.ParseFloat(y.PID, 64)
		fmt.Printf("PID Name: %v", y.Name)
		fmt.Printf("PID Value: %v", parsed_value)
		fmt.Printf("PID Unprased: %v", y.PID)
		fmt.Printf("PID State: %v", y.State)
		if err == nil {
		    ch <- prometheus.MustNewConstMetric(e.ServiceMetrics["PID"], prometheus.GaugeValue, parsed_value, y.Name, y.State)
		}
	}

	return nil
}
