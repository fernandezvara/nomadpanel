package api

import (
	"encoding/json"
	"net/http"

	"github.com/fernandezvara/nomadpanel/usage"
	"github.com/gorilla/mux"
)

func getUsageByDC(c *Context, w http.ResponseWriter, r *http.Request) {
	dc := mux.Vars(r)["dc"]

	if _, exists := c.usage.Datacenters[dc]; exists == false {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c.usage.Current(dc))
}

func getUsageByDCNodes(c *Context, w http.ResponseWriter, r *http.Request) {
	dc := mux.Vars(r)["dc"]

	if _, exists := c.usage.Datacenters[dc]; exists == false {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	nodes := make(map[string]*usage.Snapshot)
	for _, node := range c.usage.Datacenters[dc].Nodes {
		nodes[node.ID] = c.usage.NodeCurrent(dc, node.ID)
	}
	json.NewEncoder(w).Encode(nodes)
}
