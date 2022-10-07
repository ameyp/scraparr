package main

type Radarr struct {
	Endpoint string
	ApiKey string
}

func (self Radarr) EmitMetrics(namespace string, subsystem string, frequency int) {
	mediarr := Mediarr{Endpoint: self.Endpoint, ApiKey: self.ApiKey, ApiVersion: "v3", Prefix: "radarr_"}
	mediarr.EmitMetrics(namespace, subsystem, frequency)
}

