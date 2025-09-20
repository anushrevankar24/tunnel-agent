package internal

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Config holds the agent configuration
type Config struct {
	// Connection Configuration
	ServerURL  string
	LocalURL   string
	ConfigFile string

	// Agent Configuration
	LocalPort int
	Timeout   int
	DebugMode bool

	// Security Configuration
	AuthToken string
}

// ParseConfig parses command line flags and environment variables
func ParseConfig() (*Config, error) {
	config := &Config{
		// Default values
		LocalURL:  "http://127.0.0.1:9000",
		LocalPort: 5050,
		Timeout:   30,
		DebugMode: false,
	}

	// Parse command line flags
	flag.StringVar(&config.ServerURL, "server", "", "WebSocket server URL (e.g., ws://localhost:8080/agent/connect?id=my-test)")
	flag.StringVar(&config.LocalURL, "local", config.LocalURL, "Local service URL to tunnel")
	flag.StringVar(&config.ConfigFile, "config", "", "Optional configuration file path")
	flag.IntVar(&config.LocalPort, "port", config.LocalPort, "Local API port")
	flag.IntVar(&config.Timeout, "timeout", config.Timeout, "Request timeout in seconds")
	flag.BoolVar(&config.DebugMode, "debug", config.DebugMode, "Enable debug mode")

	flag.Parse()

	// Override with environment variables if set
	if envServerURL := os.Getenv("AGENT_SERVER_URL"); envServerURL != "" {
		config.ServerURL = envServerURL
	}

	if envLocalURL := os.Getenv("AGENT_LOCAL_URL"); envLocalURL != "" {
		config.LocalURL = envLocalURL
	}

	if envLocalPort := os.Getenv("AGENT_LOCAL_PORT"); envLocalPort != "" {
		if parsed, err := strconv.Atoi(envLocalPort); err == nil {
			config.LocalPort = parsed
		}
	}

	if envTimeout := os.Getenv("AGENT_TIMEOUT"); envTimeout != "" {
		if parsed, err := strconv.Atoi(envTimeout); err == nil {
			config.Timeout = parsed
		}
	}

	if envDebugMode := os.Getenv("DEBUG_MODE"); envDebugMode != "" {
		config.DebugMode = strings.ToLower(envDebugMode) == "true"
	}

	if envAuthToken := os.Getenv("AUTH_TOKEN"); envAuthToken != "" {
		config.AuthToken = envAuthToken
	}

	// Validation
	if config.ServerURL == "" {
		return nil, fmt.Errorf("server URL is required (use -server flag or AGENT_SERVER_URL env var)")
	}

	// Validate server URL
	if _, err := url.Parse(config.ServerURL); err != nil {
		return nil, fmt.Errorf("invalid server URL: %v", err)
	}

	// Validate local URL
	if _, err := url.Parse(config.LocalURL); err != nil {
		return nil, fmt.Errorf("invalid local URL: %v", err)
	}

	// Print configuration
	config.PrintConfig()

	return config, nil
}

// PrintConfig prints the current configuration
func (c *Config) PrintConfig() {
	log.Printf("=== Tunnel Agent Configuration ===")
	log.Printf("Server URL: %s", c.ServerURL)
	log.Printf("Local URL: %s", c.LocalURL)
	log.Printf("Local API Port: %d", c.LocalPort)
	log.Printf("Timeout: %d seconds", c.Timeout)
	log.Printf("Debug Mode: %t", c.DebugMode)
	if c.AuthToken != "" {
		log.Printf("Auth Token: [HIDDEN]")
	}
	log.Printf("===================================")
}

// GetLocalAPIAddress returns the local API address
func (c *Config) GetLocalAPIAddress() string {
	return fmt.Sprintf("127.0.0.1:%d", c.LocalPort)
}
