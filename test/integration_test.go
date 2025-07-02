package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEndToEndGeneration tests the complete flow from proto files to running MCP server
func TestEndToEndGeneration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "mcp_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test proto file
	protoDir := filepath.Join(tempDir, "proto")
	err = os.MkdirAll(protoDir, 0755)
	require.NoError(t, err)

	protoContent := `
syntax = "proto3";
package test.services;
option go_package = "github.com/test/proto/services";

service TestService {
  rpc GetItem(GetItemRequest) returns (GetItemResponse) {}
  rpc CreateItem(CreateItemRequest) returns (CreateItemResponse) {}
  rpc DeleteItem(DeleteItemRequest) returns (google.protobuf.Empty) {}
}

message GetItemRequest {
  string id = 1;
}

message GetItemResponse {
  Item item = 1;
}

message CreateItemRequest {
  Item item = 1;
}

message CreateItemResponse {
  Item item = 1;
}

message DeleteItemRequest {
  string id = 1;
}

message Item {
  string id = 1;
  string name = 2;
  string description = 3;
}
`
	err = os.WriteFile(filepath.Join(protoDir, "test.proto"), []byte(protoContent), 0644)
	require.NoError(t, err)

	// Create config file
	configContent := `
proto:
  dirs:
    - ./proto
output:
  dir: ./generated
  package: main
server:
  port: 3334
services:
  TestService:
    endpoint_env: "TEST_SERVICE_ADDR"
    default_endpoint: "localhost:50052"
`
	configPath := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Run the generator
	cmd := exec.Command("go", "run", "../cmd/mcp-proto-generator/main.go", "generate", "-c", configPath)
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Generator failed: %s", string(output))

	// Verify generated files
	generatedDir := filepath.Join(tempDir, "generated")
	assert.DirExists(t, generatedDir)

	expectedFiles := []string{
		"main.go",
		"tools.go",
		"handlers.go",
		"clients.go",
		"auth.go",
		"go.mod",
		"README.md",
	}

	for _, file := range expectedFiles {
		assert.FileExists(t, filepath.Join(generatedDir, file))
	}

	// Build the generated server
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = generatedDir
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "go mod tidy failed: %s", string(output))

	cmd = exec.Command("go", "build", "-o", "mcp-server")
	cmd.Dir = generatedDir
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	// Start the server
	serverCmd := exec.Command("./mcp-server")
	serverCmd.Dir = generatedDir
	serverCmd.Env = append(os.Environ(), "MCP_AUTH_ENABLED=false")

	err = serverCmd.Start()
	require.NoError(t, err)
	defer serverCmd.Process.Kill()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Test the server endpoints
	t.Run("ListTools", func(t *testing.T) {
		resp, err := http.Get("http://localhost:3334/tools")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result struct {
			Tools []struct {
				Name        string      `json:"name"`
				Description string      `json:"description"`
				Dangerous   bool        `json:"dangerous"`
				InputSchema interface{} `json:"inputSchema"`
			} `json:"tools"`
		}

		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Len(t, result.Tools, 3)

		// Verify tool names
		toolNames := make(map[string]bool)
		for _, tool := range result.Tools {
			toolNames[tool.Name] = true
		}

		assert.True(t, toolNames["testservice_getitem"])
		assert.True(t, toolNames["testservice_createitem"])
		assert.True(t, toolNames["testservice_deleteitem"])

		// Verify dangerous flag
		for _, tool := range result.Tools {
			if tool.Name == "testservice_deleteitem" {
				assert.True(t, tool.Dangerous)
			} else {
				assert.False(t, tool.Dangerous)
			}
		}
	})

	t.Run("CallTool", func(t *testing.T) {
		payload := map[string]interface{}{
			"method": "testservice_getitem",
			"params": map[string]interface{}{
				"id": "test-123",
			},
			"auth": map[string]interface{}{},
		}

		jsonData, err := json.Marshal(payload)
		require.NoError(t, err)

		resp, err := http.Post("http://localhost:3334/rpc", "application/json", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get a response (even if it's a mock)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		// Should have a result field
		assert.Contains(t, result, "result")
	})

	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := http.Get("http://localhost:3334/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "OK", string(body))
	})
}

// TestGeneratorWithRealProtos tests against the actual rescue-titan proto files
func TestGeneratorWithRealProtos(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check if we have access to the real proto files
	protoPath := "../../proto/main"
	if _, err := os.Stat(protoPath); os.IsNotExist(err) {
		t.Skip("Real proto files not found")
	}

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "mcp_real_proto_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create config for real protos
	configContent := fmt.Sprintf(`
proto:
  dirs:
    - %s
output:
  dir: ./generated
  package: main
server:
  port: 3335
services:
  DatabaseService:
    endpoint_env: "DATABASE_SERVICE_ADDR"
    default_endpoint: "localhost:50051"
  AuthService:
    endpoint_env: "AUTH_SERVICE_ADDR"
    default_endpoint: "localhost:50053"
`, protoPath)

	configPath := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Run generator
	cmd := exec.Command("go", "run", "../cmd/mcp-proto-generator/main.go", "generate", "-c", configPath)
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Generator failed: %s", string(output))

	// Verify generated files contain expected methods
	toolsPath := filepath.Join(tempDir, "generated", "tools.go")
	toolsContent, err := os.ReadFile(toolsPath)
	require.NoError(t, err)

	// Check for some expected methods
	expectedMethods := []string{
		"databaseservice_listprojects",
		"databaseservice_createuser",
		"authservice_login",
		"authservice_validatetoken",
	}

	for _, method := range expectedMethods {
		assert.Contains(t, string(toolsContent), method)
	}
}

// TestConcurrentRequests tests the MCP server under concurrent load
func TestConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This assumes a server is already running from TestEndToEndGeneration
	// In a real test suite, you'd manage server lifecycle properly

	concurrency := 10
	requestsPerClient := 5

	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(clientID int) {
			defer func() { done <- true }()

			client := &http.Client{
				Timeout: 5 * time.Second,
			}

			for j := 0; j < requestsPerClient; j++ {
				// Test listing tools
				resp, err := client.Get("http://localhost:3334/tools")
				if err != nil {
					t.Errorf("Client %d request %d failed: %v", clientID, j, err)
					continue
				}
				resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					t.Errorf("Client %d request %d got status %d", clientID, j, resp.StatusCode)
				}
			}
		}(i)
	}

	// Wait for all clients to complete
	for i := 0; i < concurrency; i++ {
		<-done
	}
}

// BenchmarkToolListing benchmarks the tool listing endpoint
func BenchmarkToolListing(b *testing.B) {
	// This assumes a server is running
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.Get("http://localhost:3334/tools")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// BenchmarkRPCCall benchmarks RPC calls
func BenchmarkRPCCall(b *testing.B) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	payload := map[string]interface{}{
		"method": "testservice_getitem",
		"params": map[string]interface{}{
			"id": "bench-123",
		},
		"auth": map[string]interface{}{},
	}

	jsonData, _ := json.Marshal(payload)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.Post("http://localhost:3334/rpc", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}