package main

import (
	"fmt"
	"github.com/nfrush/prometheus_service_exporter/config"
	"github.com/nfrush/prometheus_service_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
)

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9101").String()
		metricsPath = kingpin.Flag("web.telemetry-path", "Path under which to expost metrics.").Default("/metrics").String()
		configPath = kingpin.Flag("config.file", "Path where configuration file lives").Default("").String()
		serviceMets map[string]*prometheus.Desc
	)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("process_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting process_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	cfg, err := config.LoadConfig(*configPath)

	if err != nil {
		log.Errorln("Config error", err.Error())
	}

	serviceMets = exporter.AddMetrics(cfg.ServiceMetricFQDN, cfg.PIDCollection, cfg.PIDMetricFQDN)

	exporter := exporter.Exporter {
		ServiceMetrics: serviceMets,
		Config: cfg,
	}

	fmt.Printf("Config: %v", cfg)

	prometheus.MustRegister(&exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Process Exporter</title></head>
			<body>
			<h1>Process Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatal(err)
	}
}