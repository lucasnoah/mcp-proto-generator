# Simple Example

This example demonstrates basic MCP server generation from a simple user service proto.

## Usage

```bash
# Generate MCP server
../../bin/mcp-proto-gen generate

# Run the server
cd mcp-server
go run .
```

## Testing

```bash
# List available tools
curl http://localhost:3333/tools

# Call a tool
curl -X POST http://localhost:3333/rpc \
  -H "Content-Type: application/json" \
  -d '{
    "method": "userservice_getuser",
    "params": {"id": "123"},
    "auth": {"api_key": "test-key"}
  }'
```