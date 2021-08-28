package handlers

import (
	"encoding/json"
	"gitlab.com/pbobby001/bobpos_api/pkg/logger"
	"net/http"
)

type HealthCheck struct {
	ServerName string `json:"server_name"`
	Author     string `json:"author"`
	Version    string `json:"version"`
	Health     string `json:"health"`
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	health := &HealthCheck{
		ServerName: "Post It API",
		Author:     "Prince Bobby",
		Version:    "1.0.0",
		Health:     "Alive",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(health)
	if err != nil {
		_ = logger.Logger.Error("unable to check health of server")
	}
}