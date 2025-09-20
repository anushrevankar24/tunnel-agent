package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RequestMsg matches the server's request message structure
type RequestMsg struct {
	Type    string            `json:"type"`
	ReqID   string            `json:"req_id"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Query   string            `json:"query"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"` // base64 encoded
}

// RespMsg matches the server's response message structure
type RespMsg struct {
	Type    string            `json:"type"`
	ReqID   string            `json:"req_id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"` // base64 encoded
	Error   string            `json:"error,omitempty"`
}

// WSClient manages the WebSocket connection to the server
type WSClient struct {
	config    *Config
	conn      *websocket.Conn
	writeMu   sync.Mutex
	localURL  string
	reconnect chan struct{}
	stop      chan struct{}
}

// NewWSClient creates a new WebSocket client
func NewWSClient(config *Config) *WSClient {
	return &WSClient{
		config:    config,
		localURL:  config.LocalURL,
		reconnect: make(chan struct{}, 1),
		stop:      make(chan struct{}),
	}
}

// Connect establishes a WebSocket connection to the server
func (c *WSClient) Connect() error {
	var err error
	c.conn, _, err = websocket.DefaultDialer.Dial(c.config.ServerURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	log.Printf("Connected to server at %s", c.config.ServerURL)

	// Start message handling
	go c.handleMessages()

	return nil
}

// handleMessages processes incoming messages from the server
func (c *WSClient) handleMessages() {
	for {
		select {
		case <-c.stop:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				log.Printf("Connection lost, attempting to reconnect...")
				c.triggerReconnect()
				return
			}

			// Parse the message
			var reqMsg RequestMsg
			if err := json.Unmarshal(message, &reqMsg); err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}

			// Handle request message
			if reqMsg.Type == "request" {
				go c.handleRequest(&reqMsg)
			}
		}
	}
}

// handleRequest processes a tunnel request and forwards it to the local service
func (c *WSClient) handleRequest(reqMsg *RequestMsg) {
	log.Printf("Handling request %s: %s %s", reqMsg.ReqID, reqMsg.Method, reqMsg.Path)

	// Decode request body
	var body io.Reader
	if reqMsg.Body != "" {
		bodyBytes, err := base64.StdEncoding.DecodeString(reqMsg.Body)
		if err != nil {
			log.Printf("Failed to decode request body: %v", err)
			c.sendErrorResponse(reqMsg.ReqID, "Failed to decode request body")
			return
		}
		body = strings.NewReader(string(bodyBytes))
	}

	// Build local URL
	localURL := c.localURL + reqMsg.Path
	if reqMsg.Query != "" {
		localURL += "?" + reqMsg.Query
	}

	// Create HTTP request to local service
	req, err := http.NewRequest(reqMsg.Method, localURL, body)
	if err != nil {
		log.Printf("Failed to create local request: %v", err)
		c.sendErrorResponse(reqMsg.ReqID, "Failed to create local request")
		return
	}

	// Copy headers
	for name, value := range reqMsg.Headers {
		req.Header.Set(name, value)
	}

	// Make request to local service with timeout
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make local request: %v", err)
		c.sendErrorResponse(reqMsg.ReqID, fmt.Sprintf("Local request failed: %v", err))
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.sendErrorResponse(reqMsg.ReqID, "Failed to read response body")
		return
	}

	// Send response back to server
	c.sendResponse(reqMsg.ReqID, resp.StatusCode, resp.Header, respBody)
}

// sendResponse sends a response message back to the server
func (c *WSClient) sendResponse(reqID string, status int, headers http.Header, body []byte) {
	respMsg := RespMsg{
		Type:    "response",
		ReqID:   reqID,
		Status:  status,
		Headers: make(map[string]string),
		Body:    base64.StdEncoding.EncodeToString(body),
	}

	// Copy headers
	for name, values := range headers {
		if len(values) > 0 {
			respMsg.Headers[name] = values[0]
		}
	}

	c.writeMu.Lock()
	err := c.conn.WriteJSON(respMsg)
	c.writeMu.Unlock()

	if err != nil {
		log.Printf("Failed to send response: %v", err)
	}
}

// sendErrorResponse sends an error response back to the server
func (c *WSClient) sendErrorResponse(reqID, errorMsg string) {
	respMsg := RespMsg{
		Type:   "response",
		ReqID:  reqID,
		Status: 500,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		Error: errorMsg,
	}

	c.writeMu.Lock()
	err := c.conn.WriteJSON(respMsg)
	c.writeMu.Unlock()

	if err != nil {
		log.Printf("Failed to send error response: %v", err)
	}
}

// triggerReconnect signals that a reconnection should be attempted
func (c *WSClient) triggerReconnect() {
	select {
	case c.reconnect <- struct{}{}:
	default:
	}
}

// StartReconnectLoop starts the reconnection loop with exponential backoff
func (c *WSClient) StartReconnectLoop() {
	go func() {
		backoff := time.Second
		maxBackoff := 30 * time.Second

		for {
			select {
			case <-c.stop:
				return
			case <-c.reconnect:
				log.Printf("Attempting to reconnect in %v...", backoff)
				time.Sleep(backoff)

				if err := c.Connect(); err != nil {
					log.Printf("Reconnection failed: %v", err)
					backoff *= 2
					if backoff > maxBackoff {
						backoff = maxBackoff
					}
				} else {
					backoff = time.Second // Reset backoff on successful connection
				}
			}
		}
	}()
}

// Stop stops the client and closes the connection
func (c *WSClient) Stop() {
	close(c.stop)
	if c.conn != nil {
		c.conn.Close()
	}
}

// IsConnected returns true if the client is currently connected
func (c *WSClient) IsConnected() bool {
	return c.conn != nil
}
