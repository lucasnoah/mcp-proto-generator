package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleParser_Real(t *testing.T) {
	// Create temp directory with test proto
	tempDir, err := os.MkdirTemp("", "parser_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Write a test proto file
	protoContent := `
syntax = "proto3";
package test;
option go_package = "github.com/test/proto";

service TestService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {}
  rpc DeleteUser(DeleteUserRequest) returns (google.protobuf.Empty) {}
}

message GetUserRequest {
  string id = 1;
}

message GetUserResponse {
  User user = 1;
}

message User {
  string id = 1;
  string name = 2;
}
`
	protoPath := filepath.Join(tempDir, "test.proto")
	err = os.WriteFile(protoPath, []byte(protoContent), 0644)
	require.NoError(t, err)

	// Parse it
	parser := NewSimple([]string{})
	services, err := parser.ParseDirectoriesSimple([]string{tempDir}, nil)
	
	assert.NoError(t, err)
	assert.Len(t, services, 1)
	
	// Check service
	service := services[0]
	assert.Equal(t, "TestService", service.Name)
	assert.Equal(t, "test.proto", filepath.Base(service.File))
	assert.Len(t, service.Methods, 2)
	
	// Check methods
	methodNames := make(map[string]bool)
	for _, method := range service.Methods {
		methodNames[method.Name] = true
	}
	assert.True(t, methodNames["GetUser"])
	assert.True(t, methodNames["DeleteUser"])
}

func TestSimpleParser_NoGoPackage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "parser_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Proto without go_package option
	protoContent := `
syntax = "proto3";
package test;

service TestService {
  rpc Test(TestRequest) returns (TestResponse) {}
}
`
	protoPath := filepath.Join(tempDir, "test.proto")
	err = os.WriteFile(protoPath, []byte(protoContent), 0644)
	require.NoError(t, err)

	parser := NewSimple([]string{})
	_, err = parser.ParseDirectoriesSimple([]string{tempDir}, nil)
	
	// Should error on missing go_package
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "go_package")
}

func TestSimpleParser_MultipleServices(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "parser_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Proto with multiple services
	protoContent := `
syntax = "proto3";
package test;
option go_package = "github.com/test/proto";

service ServiceOne {
  rpc MethodOne(Request) returns (Response) {}
}

service ServiceTwo {
  rpc MethodTwo(Request) returns (Response) {}
  rpc MethodThree(Request) returns (Response) {}
}
`
	protoPath := filepath.Join(tempDir, "multi.proto")
	err = os.WriteFile(protoPath, []byte(protoContent), 0644)
	require.NoError(t, err)

	parser := NewSimple([]string{})
	services, err := parser.ParseDirectoriesSimple([]string{tempDir}, nil)
	
	assert.NoError(t, err)
	assert.Len(t, services, 2)
	
	// Find services by name
	var svc1, svc2 *ServiceDefinition
	for i := range services {
		if services[i].Name == "ServiceOne" {
			svc1 = &services[i]
		} else if services[i].Name == "ServiceTwo" {
			svc2 = &services[i]
		}
	}
	
	require.NotNil(t, svc1)
	require.NotNil(t, svc2)
	assert.Len(t, svc1.Methods, 1)
	assert.Len(t, svc2.Methods, 2)
}

func TestFindProtoFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "parser_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create proto files in different subdirs
	subDir := filepath.Join(tempDir, "sub")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	// Create files
	files := []string{
		filepath.Join(tempDir, "test1.proto"),
		filepath.Join(tempDir, "test2.proto"),
		filepath.Join(subDir, "test3.proto"),
		filepath.Join(tempDir, "not_proto.txt"),
	}

	for _, f := range files {
		err = os.WriteFile(f, []byte("test"), 0644)
		require.NoError(t, err)
	}

	// Find proto files
	found, err := findProtoFiles(tempDir, nil)
	assert.NoError(t, err)
	assert.Len(t, found, 3) // Should find 3 proto files

	// Test with exclude patterns
	found, err = findProtoFiles(tempDir, []string{"**/sub/*"})
	assert.NoError(t, err)
	assert.Len(t, found, 2) // Should exclude the sub directory
}