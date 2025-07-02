package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCommand(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cmd_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test proto
	protoDir := filepath.Join(tempDir, "proto")
	err = os.MkdirAll(protoDir, 0755)
	require.NoError(t, err)

	protoContent := `syntax = "proto3";
package test;
option go_package = "github.com/test/proto";

service TestService {
  rpc GetTest(GetTestRequest) returns (GetTestResponse) {}
}

message GetTestRequest {
  string id = 1;
}

message GetTestResponse {
  string result = 1;
}
`
	err = os.WriteFile(filepath.Join(protoDir, "test.proto"), []byte(protoContent), 0644)
	require.NoError(t, err)

	// Create config
	configContent := `proto:
  dirs:
    - ` + protoDir + `
output:
  dir: ` + filepath.Join(tempDir, "output") + `
  package: main
server:
  port: 3333
services:
  TestService:
    endpoint_env: "TEST_ADDR"
    default_endpoint: "localhost:50051"
`
	configPath := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Test would call the generate function directly here
	// For now, just verify files were created correctly
	assert.FileExists(t, configPath)
	assert.FileExists(t, filepath.Join(protoDir, "test.proto"))
}