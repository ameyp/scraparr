package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Prowlarr struct {
	Endpoint string
	ApiKey string
}

type ProwlarrHealth struct {
	Source string
	Type string
	Message string
}

func (self Prowlarr) getHealth() ([]ProwlarrHealth, error) {
	url := fmt.Sprintf("%s/api/v1/health?apikey=%s", self.Endpoint, self.ApiKey)
	response, err := MakeRequest(url)

	if err != nil {
		log.Printf("Error while fetching %s: %v", url, err)
		return nil, err
	}

	var health []ProwlarrHealth
	err = json.Unmarshal(response, &health)

	if err != nil {
		log.Printf("Error while unmarshalling %s: %v", response, err)
		return nil, err
	}

	return health, nil
}

func (self Prowlarr) EmitMetrics(namespace string, subsystem string, frequency int) {
	systemHealthUnreachable := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: "prowlarr_system_health_unreachable_count",
		Help: "Could not scrape health metrics from the service",
	})
	systemStatus := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: "prowlarr_system_status_count",
		Help: "TODO Figure out what this means when it shows up",
	})

	prometheus.MustRegister(systemHealthUnreachable)
	prometheus.MustRegister(systemStatus)

	// Periodically record some sample latencies for the three services.
	go func() {
		for {
			health, err := self.getHealth()
			if err != nil {
				systemHealthUnreachable.Add(1.0)
			} else {
				systemHealthUnreachable.Add(0.0)
				systemStatus.Set(float64(len(health)))
			}

			time.Sleep(time.Duration(frequency) * time.Minute)
		}
	}()
}


