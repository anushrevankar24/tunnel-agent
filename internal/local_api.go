package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// LocalAPI provides a simple HTTP API for controlling the agent
type LocalAPI struct {
	wsClient *WSClient
	config   *Config
}

// NewLocalAPI creates a new local API server
func NewLocalAPI(wsClient *WSClient, config *Config) *LocalAPI {
	return &LocalAPI{
		wsClient: wsClient,
		config:   config,
	}
}

// Start starts the local API server
func (api *LocalAPI) Start() {
	http.HandleFunc("/status", api.handleStatus)
	http.HandleFunc("/reload", api.handleReload)

	addr := api.config.GetLocalAPIAddress()
	log.Printf("Starting local API server on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("Local API server failed: %v", err)
	}
}

// StatusResponse represents the status endpoint response
type StatusResponse struct {
	Connected bool   `json:"connected"`
	TunnelID  string `json:"tunnel_id"`
	Local     string `json:"local"`
	Server    string `json:"server"`
}

// handleStatus returns the current agent status
func (api *LocalAPI) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract tunnel ID from server URL
	tunnelID := "unknown"
	if api.config.ServerURL != "" {
		// Simple extraction - in production, use proper URL parsing
		if idx := len(api.config.ServerURL) - 1; idx >= 0 {
			// Look for the last query parameter
			if lastSlash := strings.LastIndex(api.config.ServerURL[:idx], "/"); lastSlash != -1 {
				if queryStart := strings.Index(api.config.ServerURL[lastSlash:], "id="); queryStart != -1 {
					tunnelID = api.config.ServerURL[lastSlash+queryStart+3:]
				}
			}
		}
	}

	status := StatusResponse{
		Connected: api.wsClient.IsConnected(),
		TunnelID:  tunnelID,
		Local:     api.config.LocalURL,
		Server:    api.config.ServerURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// handleReload reloads the configuration and reconnects
func (api *LocalAPI) handleReload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Reload requested via local API")

	// For now, just trigger a reconnection
	// In a full implementation, this would reload config from file
	api.wsClient.Stop()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reload initiated"))
}
