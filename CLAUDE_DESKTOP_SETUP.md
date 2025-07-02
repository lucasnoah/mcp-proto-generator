# Connecting Claude Desktop to Your MCP Server

## Prerequisites

1. Claude Desktop app installed
2. Your MCP server running (the one we just tested)
3. Network access between Claude Desktop and your MCP server

## Step 1: Start Your MCP Server

First, make sure your MCP server is running with the real services:

```bash
# Navigate to your generated MCP server
cd /home/lucas/rescue-titan-v2/mcp-proto-generator/mcp-server

# Start the server (disable auth for testing)
MCP_AUTH_ENABLED=false ./mcp-server
```

You should see:
```
2025/07/01 20:05:30 🚀 MCP server starting on port 3333
2025/07/01 20:05:30 📊 Serving 99 tools from 7 services
```

## Step 2: Configure Claude Desktop

### Option A: Local Development (Same Machine)

1. Open Claude Desktop
2. Go to Settings → Developer → MCP Servers
3. Click "Add Server" 
4. Configure as follows:

```json
{
  "name": "Rescue Titan Services",
  "url": "http://localhost:3333",
  "description": "Access to all Rescue Titan microservices"
}
```

### Option B: Remote Server

If your MCP server is on a different machine:

```json
{
  "name": "Rescue Titan Services",
  "url": "http://YOUR_SERVER_IP:3333",
  "description": "Access to all Rescue Titan microservices"
}
```

## Step 3: Enable Authentication (Production)

For production use, enable API key authentication:

1. Set environment variables when starting the server:

```bash
# Set API keys (comma-separated for multiple)
export MCP_API_KEYS="your-secret-key-1,your-secret-key-2"
export MCP_AUTH_ENABLED=true

# Start the server
./mcp-server
```

2. Update Claude Desktop configuration:

```json
{
  "name": "Rescue Titan Services",
  "url": "http://localhost:3333",
  "description": "Access to all Rescue Titan microservices",
  "auth": {
    "type": "api_key",
    "api_key": "your-secret-key-1"
  }
}
```

## Step 4: Test the Connection

In Claude Desktop, you should now be able to:

1. See "Rescue Titan Services" in the available tools
2. Use commands like:
   - "List all users in the database"
   - "Show me active projects"
   - "Trigger an ALBI sync"

## Step 5: Configure Service Endpoints

Make sure your services are actually running:

```bash
# Set service addresses if not using defaults
export DATABASE_SERVICE_ADDR=localhost:50051
export AUTH_SERVICE_ADDR=localhost:50053
export ALBI_SYNC_SERVICE_ADDR=localhost:50054
# ... etc for other services

# Start MCP server
./mcp-server
```

## Available Tools

Once connected, Claude will have access to all 99 tools including:

### Database Operations
- `databaseservice_listusers` - List all users
- `databaseservice_listprojects` - List all projects  
- `databaseservice_listcollections` - List collections
- `databaseservice_createuser` - Create new user
- etc.

### ALBI Sync Operations
- `albisyncservice_triggerfullsync` - Trigger full ALBI sync
- `albisyncservice_triggerprojectsync` - Sync specific project
- `albisyncservice_getsyncstatus` - Check sync status
- etc.

### Auth Operations
- `authservice_login` - User login
- `authservice_validatetoken` - Validate auth token
- `authservice_createapikey` - Create API key
- etc.

## Troubleshooting

### Connection Issues

1. Check server is running:
```bash
curl http://localhost:3333/health
# Should return: OK
```

2. List available tools:
```bash
curl http://localhost:3333/tools | jq '.tools[0:3]'
```

3. Test a tool call:
```bash
curl -X POST http://localhost:3333/rpc \
  -H "Content-Type: application/json" \
  -d '{"method": "databaseservice_listusers", "params": {}, "auth": {}}'
```

### Common Problems

**"Connection refused"**
- Ensure MCP server is running
- Check firewall settings
- Verify correct port (3333)

**"Authentication failed"**
- Check MCP_AUTH_ENABLED setting
- Verify API key matches

**"Method not found"**  
- Tool name might be lowercase with underscores
- Check available tools at /tools endpoint

**"Not implemented"**
- The underlying gRPC service needs that method implemented
- This is expected for some methods

## Advanced Configuration

### Custom Port
```bash
PORT=8080 ./mcp-server
```

### Multiple Environments

Create different configs for different environments:

```bash
# Development
MCP_AUTH_ENABLED=false ./mcp-server

# Staging  
MCP_AUTH_ENABLED=true \
MCP_API_KEYS="staging-key" \
DATABASE_SERVICE_ADDR=staging-db:50051 \
./mcp-server

# Production
MCP_AUTH_ENABLED=true \
MCP_API_KEYS="prod-key-1,prod-key-2" \
DATABASE_SERVICE_ADDR=prod-db:50051 \
AUTH_SERVICE_ADDR=prod-auth:50053 \
./mcp-server
```

### Monitoring

Enable request logging:
```bash
MCP_LOG_REQUESTS=true ./mcp-server
```

## Usage Examples

Once connected, you can ask Claude:

1. **"Show me all active projects"**
   - Claude will use `databaseservice_listprojects`

2. **"List users with admin role"**
   - Claude will use `databaseservice_listusers` with filtering

3. **"Trigger a sync for project ABC123"**
   - Claude will use `albisyncservice_triggerprojectsync`

4. **"What collections need attention?"**
   - Claude will use `databaseservice_listcollections` and analyze

5. **"Send a notification to user X"**
   - Claude will use `notificationservice_sendemail`

## Security Notes

1. **Never expose MCP server to public internet without auth**
2. **Use strong API keys in production**
3. **Consider TLS/HTTPS for remote connections**
4. **Implement rate limiting for production use**
5. **Monitor and log all access**

## Next Steps

1. Implement missing service methods that return "not implemented"
2. Add more sophisticated auth (OAuth2, JWT)
3. Create service-specific MCP servers for granular access control
4. Add request/response logging for debugging
5. Implement caching for frequently accessed data

That's it! Your Claude Desktop should now have full access to all your microservices through the MCP protocol.