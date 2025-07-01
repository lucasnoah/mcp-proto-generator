# mcp-proto-generator

Generate MCP (Model Context Protocol) servers from gRPC proto files with zero manual code.

## Quick Start

```bash
# Install
go install github.com/lucasnoah/mcp-proto-generator/cmd/mcp-proto-gen@latest

# Generate MCP server from your protos
mcp-proto-gen generate --proto-dir ./proto

# Run the generated server
cd mcp-server && go run .
```

## Features

- 🚀 **Zero Manual Code** - 100% generated from proto definitions
- 🔑 **API Key Auth** - Built-in authentication with environment config
- 🛠️ **Smart Conventions** - Automatic service discovery and naming
- 📝 **Type Safe** - Full proto type mapping to JSON Schema
- 🐳 **Docker Ready** - Generated Dockerfile and deployment configs
- ⚡ **Fast Generation** - Handles large proto projects in seconds

## How It Works

1. **Parse** - Reads all your .proto files and imports
2. **Extract** - Finds services, methods, and message types
3. **Generate** - Creates MCP server with tools for each RPC method
4. **Deploy** - Single binary with zero dependencies

## Configuration

Create `mcp-gen.yaml` in your project:

```yaml
proto:
  dirs: ["./proto"]
  
generate:
  output_dir: "./mcp-server"
  module: "github.com/yourorg/mcp-server"
  
auth:
  enabled: true
  env_var: "MCP_API_KEYS"
  
server:
  port: 3333
```

## Examples

### Basic Usage
```bash
# Generate from current directory
mcp-proto-gen generate

# Specify proto directory
mcp-proto-gen generate --proto-dir ./api/proto

# Custom config
mcp-proto-gen generate --config mcp-gen.yaml
```

### Integration Example
```yaml
# .github/workflows/mcp.yml
- name: Generate MCP Server
  run: |
    mcp-proto-gen generate
    cd mcp-server && go test ./...
```

## Proto to MCP Mapping

Your proto services become MCP tools:

```proto
service UserService {
  rpc GetUser(GetUserRequest) returns (User);
  rpc DeleteUser(DeleteUserRequest) returns (Empty);
}
```

Generates MCP tools:
- `user_service_get_user` - Get user information
- `user_service_delete_user` - Delete user (DANGEROUS)

## Requirements

- Go 1.21+
- protoc compiler
- Your project's proto files

## Contributing

1. Fork the repository
2. Create feature branch
3. Add tests
4. Submit pull request

## License

MIT License - see [LICENSE](LICENSE)