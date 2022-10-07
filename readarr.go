package main

type Readarr struct {
	Endpoint string
	ApiKey string
}

func (self Readarr) EmitMetrics(namespace string, subsystem string, frequency int) {
	mediarr := Mediarr{Endpoint: self.Endpoint, ApiKey: self.ApiKey, ApiVersion: "v1", Prefix: "readarr_"}
	mediarr.EmitMetrics(namespace, subsystem, frequency)
}

