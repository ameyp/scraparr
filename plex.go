package main

import (
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Plex struct {
	Endpoint string
	ApiKey string
}

func (self Plex) isUp() bool {
	url := fmt.Sprintf("%s?X-Plex-Token=%s", self.Endpoint, self.ApiKey)
	_, err := MakeRequest(url)

	if err != nil {
		log.Printf("Error while fetching %s: %v", url, err)
	}

	return err == nil
}

func (self Plex) EmitMetrics(namespace string, subsystem string, frequency int) {
	systemHealthUnreachable := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: "plex_system_health_unreachable_count",
		Help: "Could not scrape health metrics from the service",
	})
	prometheus.MustRegister(systemHealthUnreachable)

	go func() {
		for {
			if !self.isUp() {
				systemHealthUnreachable.Add(1.0)
			} else {
				systemHealthUnreachable.Add(0.0)
			}

			time.Sleep(time.Duration(frequency) * time.Minute)
		}
	}()
}


