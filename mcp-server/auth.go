package main

import (
	"fmt"
	"os"
	"strings"
)

// authenticate validates the API key from the request
func (s *MCPServer) authenticate(auth map[string]interface{}) error {
	// Check if auth is disabled
	if os.Getenv("MCP_AUTH_ENABLED") == "false" {
		return nil
	}

	if auth == nil {
		return fmt.Errorf("missing auth")
	}

	apiKey, ok := auth["api_key"].(string)
	if !ok || apiKey == "" {
		return fmt.Errorf("missing api_key")
	}

	// Get valid API keys from environment
	validKeys := os.Getenv("MCP_API_KEYS")
	if validKeys == "" {
		return fmt.Errorf("no API keys configured")
	}

	// Check if provided key is valid
	for _, key := range strings.Split(validKeys, ",") {
		if strings.TrimSpace(key) == apiKey {
			return nil
		}
	}

	return fmt.Errorf("invalid api key")
}
