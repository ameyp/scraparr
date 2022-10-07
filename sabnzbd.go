package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Sabnzbd struct {
	Endpoint string
	ApiKey string
}

type SabnzbdHealth struct {
	have_warnings int
}

func (sabnzbd Sabnzbd) getHealth() (*SabnzbdHealth, error) {
	url := fmt.Sprintf("%s/api?mode=status&skip_dashboard=0&output=json&apikey=%s", sabnzbd.Endpoint, sabnzbd.ApiKey)
	response, err := MakeRequest(url)

	if err != nil {
		log.Printf("Error while fetching %s: %v", url, err)
		return nil, err
	}

	var health SabnzbdHealth
	err = json.Unmarshal(response, &health)
	if err != nil {
		log.Printf("Error while unmarshalling %s: %v", response, err)
		return nil, err
	}

	return &health, nil
}

func (sabnzbd Sabnzbd) EmitMetrics(namespace string, subsystem string, frequency int) {
	systemHealthUnreachable := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: "sabnzbd_system_health_unreachable_count",
		Help: "Could not scrape health metrics from the service",
	})
	healthGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: "sabnzbd_system_warnings_count",
		Help: "Number of outstanding warnings reported by Sabnzbd",
	})
	prometheus.MustRegister(healthGauge)

	// Periodically record some sample latencies for the three services.
	go func() {
		for {
			health, err := sabnzbd.getHealth()
			if err != nil {
				systemHealthUnreachable.Add(1.0)
			} else {
				systemHealthUnreachable.Add(0.0)
				healthGauge.Set(float64(health.have_warnings))
			}
			time.Sleep(time.Duration(frequency) * time.Minute)
		}
	}()
}

