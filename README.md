# database-query-server
Sample MCP Server - Go (database-query-server)

### Execute a simple SELECT query
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "execute_query",
    "arguments": {
      "database": "primary",
      "query": "SELECT id, name, email FROM users WHERE active = $1",
      "parameters": {"1": true},
      "format": "json",
      "limit": 100
    }
  }
}
```

### Get table schema

```json
echo '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call", 
  "params": {
    "name": "get_schema",
    "arguments": {
      "database": "primary",
      "tables": ["users", "orders"],
      "detailed": true
    }
  }
}' 
```

### Example Usage

```json
# Basic JSON-RPC request
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "MCP-Protocol-Version: 2024-11-05" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
    "name": "execute_query",
    "arguments": {
      "database": "primary",
      "query": "SELECT id, name, email FROM users WHERE active = $1",
      "parameters": {"1": true},
      "format": "json",
      "limit": 100
    }
  }
  }'
```

  ### Initialize

  ```json
  curl -v -X POST http://localhost:8080/mcp \
     -H "Content-Type: application/json" \
     -d '{
       "jsonrpc": "2.0",
       "id": 1,
       "method": "initialize",
       "params": {
         "protocolVersion": "2025-03-26",
         "capabilities": {
           "tools": {},
           "resources": {},
           "prompts": {}
         },
         "clientInfo": {
           "name": "curl-client",
           "version": "1.0"
         }
       }
     }'
```

