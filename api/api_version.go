package api

import (
	"encoding/json"
	"net/http"
)

type version struct {
	Version string
}

func getVersion(c *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version{Version: c.version})
}
