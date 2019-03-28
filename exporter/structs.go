package exporter

import (
	"github.com/nfrush/prometheus_service_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	ServiceMetrics map[string]*prometheus.Desc
	PidMetrics map[string]*prometheus.Desc
	config.Config
}

type Service struct {
	Name string
	State string
	Substate string
}

func (s *Service) IsActive() float64 {
	if s.State == "active" {
		return float64(1)
	}
	return float64(0)
}

type PidFile struct {
	Name string
	State string
	PID  string
}