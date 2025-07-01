package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Tool represents an MCP tool
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
	Dangerous   bool        `json:"dangerous,omitempty"`
}

// MCPServer represents the MCP server
type MCPServer struct {
	clients *Clients
}

// handleListTools returns the list of available tools
func (s *MCPServer) handleListTools(w http.ResponseWriter, r *http.Request) {
	tools := []Tool{
		{
			Name:        "userservice_getuser",
			Description: "GetUser operation from UserService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "userservice_listusers",
			Description: "ListUsers operation from UserService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "userservice_createuser",
			Description: "CreateUser operation from UserService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   false,
		},
		{
			Name:        "userservice_deleteuser",
			Description: "DeleteUser operation from UserService service",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Dangerous:   true,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tools": tools,
	})
}

// handleRPC handles MCP RPC calls
func (s *MCPServer) handleRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Method string                 `json:"method"`
		Params map[string]interface{} `json:"params"`
		Auth   map[string]interface{} `json:"auth"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Authenticate request
	if err := s.authenticate(req.Auth); err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
		return
	}

	// Route to appropriate handler
	var result interface{}
	var err error

	switch req.Method {
	case "userservice_getuser":
		result, err = s.handleUserServiceGetUser(r.Context(), req.Params)
	case "userservice_listusers":
		result, err = s.handleUserServiceListUsers(r.Context(), req.Params)
	case "userservice_createuser":
		result, err = s.handleUserServiceCreateUser(r.Context(), req.Params)
	case "userservice_deleteuser":
		result, err = s.handleUserServiceDeleteUser(r.Context(), req.Params)
	default:
		http.Error(w, "Unknown method", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": result,
	})
}
