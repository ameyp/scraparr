package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Sabnzbd struct {
		Endpoint string
		ApiKey string
	}
	Radarr struct {
		Endpoint string
		ApiKey string
	}
	Sonarr struct {
		Endpoint string
		ApiKey string
	}
	Readarr struct {
		Endpoint string
		ApiKey string
	}
	Prowlarr struct {
		Endpoint string
		ApiKey string
	}
	Plex struct {
		Endpoint string
		ApiKey string
	}
}

func readConfig() Config {
	filePath := os.Getenv("CONFIG_FILE")

	if filePath == "" {
		return Config{}
	}

	fil, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(fil, &config)

	if err != nil {
		log.Fatal(err)
	}

	return config
}

func main() {
	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))

	config := readConfig()

	if config.Sabnzbd.Endpoint != "" && config.Sabnzbd.ApiKey != "" {
		sabnzbd := Sabnzbd{
			Endpoint: config.Sabnzbd.Endpoint,
			ApiKey: config.Sabnzbd.ApiKey,
		}
		sabnzbd.EmitMetrics("wirywolf", "downloaders", 5)
	}

	if config.Radarr.Endpoint != "" && config.Radarr.ApiKey != "" {
		radarr := Radarr{
			Endpoint: config.Radarr.Endpoint,
			ApiKey: config.Radarr.ApiKey,
		}
		radarr.EmitMetrics("wirywolf", "downloaders", 5)
	}

	if config.Sonarr.Endpoint != "" && config.Sonarr.ApiKey != "" {
		sonarr := Sonarr{
			Endpoint: config.Sonarr.Endpoint,
			ApiKey: config.Sonarr.ApiKey,
		}
		sonarr.EmitMetrics("wirywolf", "downloaders", 5)
	}

	if config.Readarr.Endpoint != "" && config.Readarr.ApiKey != "" {
		readarr := Readarr{
			Endpoint: config.Readarr.Endpoint,
			ApiKey: config.Readarr.ApiKey,
		}
		readarr.EmitMetrics("wirywolf", "downloaders", 5)
	}

	if config.Prowlarr.Endpoint != "" && config.Prowlarr.ApiKey != "" {
		prowlarr := Prowlarr{
			Endpoint: config.Prowlarr.Endpoint,
			ApiKey: config.Prowlarr.ApiKey,
		}
		prowlarr.EmitMetrics("wirywolf", "downloaders", 5)
	}

	if config.Plex.Endpoint != "" && config.Plex.ApiKey != "" {
		plex := Plex{
			Endpoint: config.Plex.Endpoint,
			ApiKey: config.Plex.ApiKey,
		}
		plex.EmitMetrics("wirywolf", "downloaders", 5)
	}

	log.Println("Publishing metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
