package api

import (
	"encoding/json"
	"net/http"
)

type regionDCs struct {
	Region      string   `json:"region"`
	Datacenters []string `json:"datacenters"`
}

func getDCs(c *Context, w http.ResponseWriter, r *http.Request) {
	dcs := make([]string, 0, len(c.usage.Datacenters))
	for k := range c.usage.Datacenters {
		dcs = append(dcs, k)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dcs)
}

type regionStruct struct {
	Region string `json:"region"`
}

func getRegion(c *Context, w http.ResponseWriter, r *http.Request) {
	region, err := c.nomadClient.Agent().Region()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(er{err.Error()})
		w.WriteHeader(500)
		return
	}
	json.NewEncoder(w).Encode(regionStruct{region})
}
