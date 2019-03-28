package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, v := range e.ServiceMetrics {
		ch <- v
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if e.Config.Supervisor == "systemd" {
		var data, err = e.systemd()

		if err != nil {
			log.Errorf("Error gathering service metrics: %v", err)
			return
		}

		err = e.processMetrics(data, nil, ch)

		if err != nil {
			log.Error("Error processing service metrics: %v", err)
			return
		}
	}

	if e.Config.PIDCollection {
		var data, err = e.pid()

		if err != nil {
			log.Errorf("Error gathering pid metrics: %v", err)
			return
		}

		err = e.processMetrics(nil, data, ch)

		if err != nil {
			log.Error("Error processing pid metrics: %v", err)
			return
		}
	}

	log.Info("All metrics successfully collected")
}
