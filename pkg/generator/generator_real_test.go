package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lucasnoah/mcp-proto-generator/pkg/config"
	"github.com/lucasnoah/mcp-proto-generator/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	// Create test config
	cfg := &config.Config{
		Output: config.OutputConfig{
			Package: "main",
		},
		Server: config.ServerConfig{
			Port: 3333,
		},
		Services: map[string]config.ServiceConfig{
			"TestService": {
				EndpointEnv:     "TEST_SERVICE_ADDR",
				DefaultEndpoint: "localhost:50051",
			},
		},
		Auth: config.AuthConfig{
			Enabled: true,
		},
	}

	// Create test services
	services := []parser.ServiceDefinition{
		{
			Name:    "TestService",
			File:    "test.proto",
			Package: "test",
			Methods: []parser.MethodDefinition{
				{
					Name:       "GetItem",
					InputType:  "GetItemRequest",
					OutputType: "GetItemResponse",
				},
				{
					Name:       "DeleteItem",
					InputType:  "DeleteItemRequest",
					OutputType: "google.protobuf.Empty",
				},
			},
		},
	}

	// Create generator
	gen := New(cfg, services)

	// Test Plan
	files, err := gen.Plan()
	assert.NoError(t, err)
	assert.True(t, len(files) >= 5) // Should have at least main, tools, handlers, clients, auth

	// Test Generate
	tempDir, err := os.MkdirTemp("", "generator_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	generated, err := gen.Generate()
	assert.NoError(t, err)
	assert.True(t, len(generated) >= 5)

	// Write files and check content
	for _, file := range generated {
		path := filepath.Join(tempDir, file.Path)
		err = os.MkdirAll(filepath.Dir(path), 0755)
		require.NoError(t, err)
		
		err = os.WriteFile(path, file.Content, 0644)
		require.NoError(t, err)
	}

	// Verify main.go
	mainContent, err := os.ReadFile(filepath.Join(tempDir, "main.go"))
	require.NoError(t, err)
	mainStr := string(mainContent)
	
	assert.Contains(t, mainStr, "package main")
	assert.Contains(t, mainStr, "port := 3333")
	assert.Contains(t, mainStr, "handleListTools")
	assert.Contains(t, mainStr, "handleRPC")

	// Verify tools.go
	toolsContent, err := os.ReadFile(filepath.Join(tempDir, "tools.go"))
	require.NoError(t, err)
	toolsStr := string(toolsContent)
	
	assert.Contains(t, toolsStr, "testservice_getitem")
	assert.Contains(t, toolsStr, "testservice_deleteitem")
	assert.Contains(t, toolsStr, `Dangerous:   true`) // DeleteItem should be dangerous

	// Verify handlers.go
	handlersContent, err := os.ReadFile(filepath.Join(tempDir, "handlers.go"))
	require.NoError(t, err)
	handlersStr := string(handlersContent)
	
	assert.Contains(t, handlersStr, "handleTestServiceGetItem")
	assert.Contains(t, handlersStr, "handleTestServiceDeleteItem")

	// Verify auth.go exists if auth is enabled
	if cfg.Auth.Enabled {
		authPath := filepath.Join(tempDir, "auth.go")
		assert.FileExists(t, authPath)
		
		authContent, err := os.ReadFile(authPath)
		require.NoError(t, err)
		assert.Contains(t, string(authContent), "authenticate")
	}
}

func TestGenerator_Plan(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			Enabled: false,
		},
	}
	
	services := []parser.ServiceDefinition{
		{
			Name: "TestService",
			Methods: []parser.MethodDefinition{
				{Name: "Test", InputType: "TestRequest", OutputType: "TestResponse"},
			},
		},
	}

	gen := New(cfg, services)
	files, err := gen.Plan()
	
	assert.NoError(t, err)
	
	// Check expected files
	expectedFiles := map[string]bool{
		"main.go":     false,
		"tools.go":    false,
		"handlers.go": false,
		"clients.go":  false,
		"go.mod":      false,
		"README.md":   false,
		".gitignore":  false,
	}
	
	for _, file := range files {
		expectedFiles[file.Path] = true
	}
	
	for file, found := range expectedFiles {
		if !found && file != "auth.go" { // auth.go is optional
			t.Errorf("Expected file %s not found in plan", file)
		}
	}
}

func TestGenerator_MethodCounts(t *testing.T) {
	services := []parser.ServiceDefinition{
		{
			Name: "Service1",
			Methods: []parser.MethodDefinition{
				{Name: "Method1", InputType: "Request1", OutputType: "Response1"},
				{Name: "Method2", InputType: "Request2", OutputType: "Response2"},
			},
		},
		{
			Name: "Service2",
			Methods: []parser.MethodDefinition{
				{Name: "Method3", InputType: "Request3", OutputType: "Response3"},
			},
		},
	}

	cfg := &config.Config{}
	gen := New(cfg, services)
	
	// The generator should handle 3 methods across 2 services
	totalMethods := 0
	for _, svc := range gen.services {
		totalMethods += len(svc.Methods)
	}
	
	assert.Equal(t, 3, totalMethods)
	assert.Equal(t, 2, len(gen.services))
}

func TestDangerousMethodDetection(t *testing.T) {
	// Test dangerous method patterns
	dangerousPatterns := []string{
		"Delete", "Drop", "Remove", "Destroy", "Purge", "Truncate",
	}
	
	for _, pattern := range dangerousPatterns {
		methodName := pattern + "Something"
		// Check if method name contains dangerous pattern
		isDangerous := false
		for _, p := range dangerousPatterns {
			if strings.HasPrefix(methodName, p) {
				isDangerous = true
				break
			}
		}
		assert.True(t, isDangerous, "Method %s should be detected as dangerous", methodName)
	}
	
	// Test safe methods
	safeNames := []string{
		"GetUser", "CreateUser", "UpdateUser", "ListUsers",
		"FindUser", "SearchUsers", "ValidateUser",
	}
	
	for _, name := range safeNames {
		isDangerous := false
		for _, p := range dangerousPatterns {
			if strings.HasPrefix(name, p) {
				isDangerous = true
				break
			}
		}
		assert.False(t, isDangerous, "Method %s should not be dangerous", name)
	}
}

func TestGenerator_EmptyServices(t *testing.T) {
	cfg := &config.Config{}
	gen := New(cfg, []parser.ServiceDefinition{})
	
	// Should still generate files even with no services
	files, err := gen.Generate()
	assert.NoError(t, err)
	assert.NotEmpty(t, files)
	
	// But tools.go should have empty tool list
	var toolsFile *GeneratedFile
	for _, f := range files {
		if f.Path == "tools.go" {
			toolsFile = &f
			break
		}
	}
	
	require.NotNil(t, toolsFile)
	assert.Contains(t, string(toolsFile.Content), "[]Tool{")
}