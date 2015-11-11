package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	nomadapi "github.com/hashicorp/nomad/api"
)

func getJobs(c *Context, w http.ResponseWriter, r *http.Request) {
	q := &nomadapi.QueryOptions{}
	jobs, _, err := c.nomadClient.Jobs().List(q)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(er{err.Error()})
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(jobs)
}

func getJob(c *Context, w http.ResponseWriter, r *http.Request) {
	jobID := mux.Vars(r)["id"]
	q := &nomadapi.QueryOptions{}
	job, _, err := c.nomadClient.Jobs().Info(jobID, q)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(er{err.Error()})
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(job)
}
