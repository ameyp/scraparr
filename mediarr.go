package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Radarr, Sonarr and Readarr all have the same API for our purposes
type Mediarr struct {
	Endpoint string
	ApiKey string
	ApiVersion string
	Prefix string
}

type MediarrHealth struct {
	Source string
	Type string
	Message string
}

type MediarrQueueStatus struct {
	unknownCount int
	errors bool
	warnings bool
	unknownErrors bool
	unknownWarnings bool
}

func (self Mediarr) getHealth() ([]MediarrHealth, error) {
	url := fmt.Sprintf("%s/api/%s/health?apikey=%s", self.Endpoint, self.ApiVersion, self.ApiKey)
	response, err := MakeRequest(url)

	if err != nil {
		log.Printf("Error while fetching %s: %v", url, err)
		return nil, err
	}

	var health []MediarrHealth
	err = json.Unmarshal(response, &health)

	if err != nil {
		log.Printf("Error while unmarshalling %s: %v", response, err)
		return nil, err
	}

	return health, nil
}

func (self Mediarr) getQueueStatus() (*MediarrQueueStatus, error) {
	url := fmt.Sprintf("%s/api/%s/queue/status?apikey=%s", self.Endpoint, self.ApiVersion, self.ApiKey)
	response, err := MakeRequest(url)

	if err != nil {
		return nil, err
	}

	var queueStatus MediarrQueueStatus
	err = json.Unmarshal(response, &queueStatus)

	if err != nil {
		return nil, err
	}

	return &queueStatus, nil
}

func (self Mediarr) EmitMetrics(namespace string, subsystem string, frequency int) {
	systemHealthUnreachable := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: self.Prefix + "system_health_unreachable_count",
		Help: "Could not scrape health metrics from the service",
	})
	systemStatus := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: self.Prefix + "system_status_count",
		Help: "TODO Figure out what this means when it shows up",
	})
	queueUnreachable := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: self.Prefix + "queue_unreachable_count",
		Help: "Could not scrape queue metrics from the service",
	})
	queueUnknownCount := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: self.Prefix + "queue_unknown_count",
		Help: "Number of unknown items in the queue",
	})
	queueErrorCount := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: self.Prefix + "queue_error_count",
		Help: "Number of errors with the queue",
	})
	queueWarningCount := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: self.Prefix + "queue_warning_count",
		Help: "Number of warnings with the queue",
	})
	queueUnknownErrorCount := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: self.Prefix + "queue_unknown_error_count",
		Help: "Number of unknown errors with the queue",
	})
	queueUnknownWarningCount := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name: self.Prefix + "queue_unknown_warning_count",
		Help: "Number of unknown warnings with the queue",
	})
	prometheus.MustRegister(systemHealthUnreachable)
	prometheus.MustRegister(systemStatus)
	prometheus.MustRegister(queueUnreachable)
	prometheus.MustRegister(queueUnknownCount)
	prometheus.MustRegister(queueErrorCount)
	prometheus.MustRegister(queueWarningCount)
	prometheus.MustRegister(queueUnknownErrorCount)
	prometheus.MustRegister(queueUnknownWarningCount)

	// Periodically record some sample latencies for the three services.
	go func() {
		for {
			health, err := self.getHealth()
			if err != nil {
				systemHealthUnreachable.Set(1.0)
			} else {
				systemHealthUnreachable.Set(0.0)
				systemStatus.Set(float64(len(health)))
			}

			queueStatus, err := self.getQueueStatus()
			if err != nil {
				queueUnreachable.Set(1.0)
			} else {
				queueUnreachable.Set(0.0)
				queueUnknownCount.Set(float64(queueStatus.unknownCount))

				if queueStatus.unknownErrors {
					queueUnknownErrorCount.Set(1.0)
				} else {
					queueUnknownErrorCount.Set(0.0)
				}

				if queueStatus.unknownWarnings {
					queueUnknownWarningCount.Set(1.0)
				} else {
					queueUnknownWarningCount.Set(0.0)
				}

				if queueStatus.errors {
					queueErrorCount.Set(1.0)
				} else {
					queueErrorCount.Set(0.0)
				}

				if queueStatus.warnings {
					queueWarningCount.Set(1.0)
				} else {
					queueWarningCount.Set(0.0)
				}
			}

			time.Sleep(time.Duration(frequency) * time.Minute)
		}
	}()
}

