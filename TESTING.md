# MCP Proto Generator Testing Guide

## Overview

This document describes the comprehensive test suite for the MCP Proto Generator.

## Test Structure

```
test/
├── fixtures/           # Sample proto files for testing
│   ├── simple.proto   # Basic service with CRUD operations
│   └── complex.proto  # Complex service with streaming, nested types
├── smoke_test.go      # End-to-end smoke tests
└── integration_test.go # Full integration test suite
```

## Running Tests

### All Tests
```bash
make test              # Run all tests (unit + integration)
make test-all          # Run all tests with coverage and benchmarks
```

### Unit Tests Only
```bash
make test-unit         # Run unit tests (fast, no external deps)
go test -v -short ./pkg/...
```

### Integration Tests
```bash
make test-integration  # Run integration tests
go test -v ./test/...
```

### Coverage
```bash
make test-coverage     # Generate coverage report
# Opens coverage.html in browser
```

### Benchmarks
```bash
make test-benchmark    # Run performance benchmarks
go test -bench=. -benchmem ./test/...
```

### Specific Tests
```bash
make test-specific TEST=TestFunctionName
go test -v -run TestFunctionName ./...
```

## Test Categories

### 1. Parser Tests (`pkg/parser/*_test.go`)

Tests the proto file parsing functionality:

- **TestSimpleParser_Real**: Tests parsing real proto files
- **TestSimpleParser_NoGoPackage**: Validates go_package requirement
- **TestSimpleParser_MultipleServices**: Tests multiple services in one file
- **TestFindProtoFiles**: Tests proto file discovery

**Key test cases:**
- Valid proto parsing with services and methods
- Missing go_package option handling
- Multiple services in single file
- Recursive directory scanning
- Exclude pattern support

### 2. Generator Tests (`pkg/generator/*_test.go`)

Tests code generation from parsed protos:

- **TestGenerator_Generate**: Full generation test
- **TestGenerator_Plan**: Tests file planning without generation
- **TestGenerator_MethodCounts**: Validates method counting
- **TestDangerousMethodDetection**: Tests dangerous operation detection

**Key test cases:**
- Correct file generation (main.go, tools.go, etc.)
- Dangerous method flagging (Delete*, Drop*, etc.)
- Empty service handling
- Multi-service generation

### 3. Integration Tests (`test/integration_test.go`)

Full end-to-end tests:

- **TestEndToEndGeneration**: Complete flow from proto to running server
- **TestGeneratorWithRealProtos**: Tests with actual rescue-titan protos
- **TestConcurrentRequests**: Load testing the generated server
- **BenchmarkToolListing**: Performance benchmarks

**Key validations:**
- Generated server compiles and runs
- All tools are exposed correctly
- RPC endpoints work
- Concurrent request handling

### 4. Smoke Tests (`test/smoke_test.go`)

Quick validation tests:

- **TestMCPProtoGeneratorSmoke**: Basic functionality test
- **TestRealProtosGeneration**: Test with real proto files

## Test Fixtures

### simple.proto
Basic CRUD service for testing standard operations:
```protobuf
service SimpleService {
  rpc GetItem(GetItemRequest) returns (GetItemResponse) {}
  rpc CreateItem(CreateItemRequest) returns (CreateItemResponse) {}
  rpc DeleteItem(DeleteItemRequest) returns (google.protobuf.Empty) {}
}
```

### complex.proto
Advanced features testing:
- Streaming RPCs
- Nested message types
- Enums and oneofs
- Map fields
- Admin/dangerous operations

## Writing New Tests

### Unit Test Template
```go
func TestFeatureName(t *testing.T) {
    // Arrange
    tempDir, err := os.MkdirTemp("", "test")
    require.NoError(t, err)
    defer os.RemoveAll(tempDir)
    
    // Act
    result, err := FunctionUnderTest(input)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Integration Test Template
```go
func TestIntegrationScenario(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Full end-to-end test
}
```

## Continuous Integration

The test suite is designed for CI/CD:

```bash
# Quick CI tests
make test-ci      # Lint + unit tests

# Full validation
make test-all     # Everything including benchmarks
```

## Known Issues

1. **Parser Tests**: Simple regex parser has limitations with complex proto syntax
2. **Integration Tests**: Require building the binary first
3. **Real Proto Tests**: Skip if proto files not available

## Test Coverage Goals

- **Overall**: 80% minimum
- **Critical paths**: 100% (parser, generator core)
- **Integration**: Key user workflows covered

## Debugging Tests

```bash
# Verbose output
go test -v ./...

# Run single test with logs
go test -v -run TestName ./pkg/... 

# Debug with delve
dlv test ./pkg/parser -- -test.run TestName
```

## Performance Testing

Benchmarks measure:
- Proto parsing speed
- Code generation time
- Generated server performance
- Concurrent request handling

Run benchmarks:
```bash
make test-benchmark
```

Compare benchmarks:
```bash
go test -bench=. -benchmem ./test/... > new.txt
benchcmp old.txt new.txt
```