package main

type Sonarr struct {
	Endpoint string
	ApiKey string
}

func (self Sonarr) EmitMetrics(namespace string, subsystem string, frequency int) {
	mediarr := Mediarr{Endpoint: self.Endpoint, ApiKey: self.ApiKey, ApiVersion: "v3", Prefix: "sonarr_"}
	mediarr.EmitMetrics(namespace, subsystem, frequency)
}

