package generator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/lucasnoah/mcp-proto-generator/pkg/config"
	"github.com/lucasnoah/mcp-proto-generator/pkg/parser"
)

// Generator generates MCP server code from parsed proto definitions
type Generator struct {
	config   *config.Config
	services []parser.ServiceDefinition
}

// GeneratedFile represents a generated file
type GeneratedFile struct {
	Path    string
	Content []byte
}

// New creates a new generator
func New(cfg *config.Config, services []parser.ServiceDefinition) *Generator {
	return &Generator{
		config:   cfg,
		services: services,
	}
}

// Plan returns what files would be generated without actually generating them
func (g *Generator) Plan() ([]GeneratedFile, error) {
	var files []GeneratedFile

	// Main server file
	files = append(files, GeneratedFile{Path: "main.go"})
	
	// Tools definition file
	files = append(files, GeneratedFile{Path: "tools.go"})
	
	// Handlers file
	files = append(files, GeneratedFile{Path: "handlers.go"})
	
	// gRPC clients file
	files = append(files, GeneratedFile{Path: "clients.go"})
	
	// Auth middleware file
	if g.config.Auth.Enabled {
		files = append(files, GeneratedFile{Path: "auth.go"})
	}
	
	// Go module files
	files = append(files, GeneratedFile{Path: "go.mod"})
	files = append(files, GeneratedFile{Path: "go.sum"})
	
	// Dockerfile
	files = append(files, GeneratedFile{Path: "Dockerfile"})
	
	// Makefile
	files = append(files, GeneratedFile{Path: "Makefile"})

	return files, nil
}

// Generate generates all MCP server files
func (g *Generator) Generate() ([]GeneratedFile, error) {
	var files []GeneratedFile

	// Generate main.go
	content, err := g.generateMain()
	if err != nil {
		return nil, fmt.Errorf("failed to generate main.go: %w", err)
	}
	files = append(files, GeneratedFile{Path: "main.go", Content: content})

	// Generate tools.go
	content, err = g.generateTools()
	if err != nil {
		return nil, fmt.Errorf("failed to generate tools.go: %w", err)
	}
	files = append(files, GeneratedFile{Path: "tools.go", Content: content})

	// Generate handlers.go
	content, err = g.generateHandlers()
	if err != nil {
		return nil, fmt.Errorf("failed to generate handlers.go: %w", err)
	}
	files = append(files, GeneratedFile{Path: "handlers.go", Content: content})

	// Generate clients.go
	content, err = g.generateClients()
	if err != nil {
		return nil, fmt.Errorf("failed to generate clients.go: %w", err)
	}
	files = append(files, GeneratedFile{Path: "clients.go", Content: content})

	// Generate auth.go if enabled
	if g.config.Auth.Enabled {
		content, err = g.generateAuth()
		if err != nil {
			return nil, fmt.Errorf("failed to generate auth.go: %w", err)
		}
		files = append(files, GeneratedFile{Path: "auth.go", Content: content})
	}

	// Generate go.mod
	content, err = g.generateGoMod()
	if err != nil {
		return nil, fmt.Errorf("failed to generate go.mod: %w", err)
	}
	files = append(files, GeneratedFile{Path: "go.mod", Content: content})

	// Generate Dockerfile
	content, err = g.generateDockerfile()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Dockerfile: %w", err)
	}
	files = append(files, GeneratedFile{Path: "Dockerfile", Content: content})

	// Generate Makefile
	content, err = g.generateMakefile()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Makefile: %w", err)
	}
	files = append(files, GeneratedFile{Path: "Makefile", Content: content})

	return files, nil
}

// generateMain generates the main.go file
func (g *Generator) generateMain() ([]byte, error) {
	tmpl := `package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// Get port from environment or use default
	port := {{.Port}}
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	// Initialize gRPC clients
	clients, err := initClients()
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	defer closeClients(clients)

	// Create MCP server
	server := &MCPServer{
		clients: clients,
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	
	// MCP endpoints
	mux.HandleFunc("/tools", server.handleListTools)
	mux.HandleFunc("/rpc", server.handleRPC)
	
{{if .HealthEnabled}}	// Health check endpoint  
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
{{end}}
{{if .MetricsEnabled}}	// Metrics endpoint
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("# Metrics endpoint placeholder"))
	})
{{end}}

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 MCP server starting on port %d", port)
		log.Printf("📊 Serving %d tools from %d services", {{.TotalTools}}, {{.TotalServices}})
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("✅ Server stopped")
}
`

	// Calculate totals
	totalTools := 0
	for _, svc := range g.services {
		for _, method := range svc.Methods {
			if !parser.IsExcludedMethod(method.Name, g.config.Generate.Methods.ExcludePatterns) {
				totalTools++
			}
		}
	}

	data := map[string]interface{}{
		"Port":           g.config.Server.Port,
		"HealthEnabled":  g.config.Server.HealthEnabled,
		"MetricsEnabled": g.config.Server.MetricsEnabled,
		"TotalTools":     totalTools,
		"TotalServices":  len(g.services),
	}

	return g.renderTemplate(tmpl, data)
}

// generateTools generates the tools.go file
func (g *Generator) generateTools() ([]byte, error) {
	tmpl := `package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Tool represents an MCP tool
type Tool struct {
	Name        string      ` + "`" + `json:"name"` + "`" + `
	Description string      ` + "`" + `json:"description"` + "`" + `
	InputSchema interface{} ` + "`" + `json:"inputSchema"` + "`" + `
	Dangerous   bool        ` + "`" + `json:"dangerous,omitempty"` + "`" + `
}

// MCPServer represents the MCP server
type MCPServer struct {
	clients *Clients
}

// handleListTools returns the list of available tools
func (s *MCPServer) handleListTools(w http.ResponseWriter, r *http.Request) {
	tools := []Tool{
{{range .Tools}}		{
			Name:        "{{.Name}}",
			Description: "{{.Description}}",
			InputSchema: {{.InputSchema}},
			Dangerous:   {{.Dangerous}},
		},
{{end}}	}

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
		Method string                 ` + "`" + `json:"method"` + "`" + `
		Params map[string]interface{} ` + "`" + `json:"params"` + "`" + `
{{if .AuthEnabled}}		Auth   map[string]interface{} ` + "`" + `json:"auth"` + "`" + `
{{end}}	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

{{if .AuthEnabled}}	// Authenticate request
	if err := s.authenticate(req.Auth); err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
		return
	}
{{end}}
	// Route to appropriate handler
	var result interface{}
	var err error

	switch req.Method {
{{range .Tools}}	case "{{.Name}}":
		result, err = s.{{.HandlerName}}(r.Context(), req.Params)
{{end}}	default:
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
`

	// Generate tool definitions
	var tools []map[string]interface{}
	
	for _, svc := range g.services {
		for _, method := range svc.Methods {
			if parser.IsExcludedMethod(method.Name, g.config.Generate.Methods.ExcludePatterns) {
				continue
			}

			toolName := g.methodToToolName(svc.Name, method.Name)
			description := g.generateDescription(svc.Name, method.Name, method.Description)
			dangerous := parser.IsDangerousMethod(method.Name, g.config.Generate.Methods.DangerousPatterns)
			handlerName := g.methodToHandlerName(svc.Name, method.Name)
			inputSchema := g.generateInputSchema(method.InputType)

			tools = append(tools, map[string]interface{}{
				"Name":        toolName,
				"Description": description,
				"InputSchema": inputSchema,
				"Dangerous":   dangerous,
				"HandlerName": handlerName,
			})
		}
	}

	data := map[string]interface{}{
		"Tools":       tools,
		"AuthEnabled": g.config.Auth.Enabled,
	}

	return g.renderTemplate(tmpl, data)
}

// generateHandlers generates the handlers.go file  
func (g *Generator) generateHandlers() ([]byte, error) {
	tmpl := `package main

import (
	"context"
)

{{range .Handlers}}
// {{.HandlerName}} handles {{.ServiceName}}.{{.MethodName}}
func (s *MCPServer) {{.HandlerName}}(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// TODO: Convert params to actual proto request type and call gRPC service
	// For now, return mock response for testing
	return map[string]interface{}{
		"message": "{{.ServiceName}}.{{.MethodName}} called successfully",
		"params": params,
	}, nil
}
{{end}}
`

	// Generate handlers
	var handlers []map[string]interface{}
	
	for _, svc := range g.services {
		for _, method := range svc.Methods {
			if parser.IsExcludedMethod(method.Name, g.config.Generate.Methods.ExcludePatterns) {
				continue
			}

			handlers = append(handlers, map[string]interface{}{
				"HandlerName": g.methodToHandlerName(svc.Name, method.Name),
				"ServiceName": svc.Name,
				"MethodName":  method.Name,
			})
		}
	}

	data := map[string]interface{}{
		"Handlers": handlers,
	}

	return g.renderTemplate(tmpl, data)
}

// generateClients generates the clients.go file
func (g *Generator) generateClients() ([]byte, error) {
	tmpl := `package main

import (
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Clients holds all gRPC clients
type Clients struct {
{{range .Services}}	{{.Name}} interface{} // TODO: Replace with actual client type
{{end}}}

// initClients initializes all gRPC clients
func initClients() (*Clients, error) {
	clients := &Clients{}

{{range .Services}}	// Initialize {{.Name}} client
	{{.Name}}Addr := os.Getenv("{{.EnvVar}}")
	if {{.Name}}Addr == "" {
		{{.Name}}Addr = "{{.DefaultAddr}}"
	}
	
	{{.Name}}Conn, err := grpc.Dial({{.Name}}Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to {{.Name}}: %w", err)
	}
	
	// TODO: Create actual client
	clients.{{.Name}} = {{.Name}}Conn

{{end}}	return clients, nil
}

// closeClients closes all gRPC connections
func closeClients(clients *Clients) {
	// TODO: Close all connections
}
`

	// Generate service client info
	var services []map[string]interface{}
	
	for _, svc := range g.services {
		envVar, defaultAddr := g.config.GetServiceEndpoint(svc.Name)
		
		services = append(services, map[string]interface{}{
			"Name":        svc.Name,
			"EnvVar":      envVar,
			"DefaultAddr": defaultAddr,
		})
	}

	data := map[string]interface{}{
		"Services": services,
	}

	return g.renderTemplate(tmpl, data)
}

// generateAuth generates the auth.go file
func (g *Generator) generateAuth() ([]byte, error) {
	tmpl := `package main

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
	validKeys := os.Getenv("{{.EnvVar}}")
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
`

	data := map[string]interface{}{
		"EnvVar": g.config.Auth.EnvVar,
	}

	return g.renderTemplate(tmpl, data)
}

// generateGoMod generates the go.mod file
func (g *Generator) generateGoMod() ([]byte, error) {
	tmpl := `module {{.Module}}

go 1.21

require (
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.33.0
)
`

	data := map[string]interface{}{
		"Module": g.config.Generate.Module,
	}

	return g.renderTemplate(tmpl, data)
}

// generateDockerfile generates the Dockerfile
func (g *Generator) generateDockerfile() ([]byte, error) {
	tmpl := `FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o mcp-server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/mcp-server .

EXPOSE {{.Port}}

CMD ["./mcp-server"]
`

	data := map[string]interface{}{
		"Port": g.config.Server.Port,
	}

	return g.renderTemplate(tmpl, data)
}

// generateMakefile generates the Makefile
func (g *Generator) generateMakefile() ([]byte, error) {
	tmpl := `.PHONY: build run test clean docker

build:
	@go build -o mcp-server .

run: build
	@./mcp-server

test:
	@go test ./...

clean:
	@rm -f mcp-server

docker:
	@docker build -t mcp-server .

docker-run: docker
	@docker run -p {{.Port}}:{{.Port}} \
		-e MCP_API_KEYS="test-key" \
		mcp-server
`

	data := map[string]interface{}{
		"Port": g.config.Server.Port,
	}

	return g.renderTemplate(tmpl, data)
}

// Helper functions

func (g *Generator) methodToToolName(serviceName, methodName string) string {
	svcName := strings.ToLower(serviceName)
	methodName = strings.ToLower(methodName)
	return fmt.Sprintf("%s_%s", svcName, methodName)
}

func (g *Generator) methodToHandlerName(serviceName, methodName string) string {
	return fmt.Sprintf("handle%s%s", serviceName, methodName)
}

func (g *Generator) generateDescription(serviceName, methodName, existingDesc string) string {
	if existingDesc != "" {
		return existingDesc
	}
	
	// Generate basic description from method name
	return fmt.Sprintf("%s operation from %s service", methodName, serviceName)
}

func (g *Generator) generateInputSchema(msgDef parser.MessageDefinition) string {
	// Simplified schema generation - real implementation would be more complete
	return `map[string]interface{}{"type": "object", "properties": map[string]interface{}{}}`
}

func (g *Generator) renderTemplate(tmplStr string, data interface{}) ([]byte, error) {
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}