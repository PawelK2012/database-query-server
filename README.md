# database-query-server
Sample MCP Server - Go (database-query-server)

# How to run project

1. Simple start with `docker compose up`

Related [documentation](https://github.com/docker/awesome-compose/tree/master/postgresql-pgadmin)

## Configuration

### .env
Before deploying this setup, you need to configure the following values in the [.env](.env) file.
- POSTGRES_USER
- POSTGRES_PW
- POSTGRES_DB (can be default value)
- PGADMIN_DEFAULT_EMAIL
- PGADMIN_DEFAULT_PASSWORD

## Deploy with docker compose
When deploying this setup, the pgAdmin web interface will be available at port 5050 (e.g. http://localhost:5050).  

``` shell
$ docker compose up
Starting postgres ... done
Starting pgadmin ... done
```

## Add postgres database to pgAdmin
After logging in with your credentials of the .env file, you can add your database to pgAdmin. 
1. Right-click "Servers" in the top-left corner and select "Create" -> "Server..."
2. Name your connection
3. Change to the "Connection" tab and add the connection details:
- Hostname: "postgres" (this would normally be your IP address of the postgres database - however, docker can resolve this container ip by its name)
- Port: "5432"
- Maintenance Database: $POSTGRES_DB (see .env)
- Username: $POSTGRES_USER (see .env)
- Password: $POSTGRES_PW (see .env)
  
## Expected result

Check containers are running:
```
$ docker ps
CONTAINER ID   IMAGE                           COMMAND                  CREATED             STATUS                 PORTS                                                                                  NAMES
849c5f48f784   postgres:latest                 "docker-entrypoint.s…"   9 minutes ago       Up 9 minutes           0.0.0.0:5432->5432/tcp, :::5432->5432/tcp                                              postgres
d3cde3b455ee   dpage/pgadmin4:latest           "/entrypoint.sh"         9 minutes ago       Up 9 minutes           443/tcp, 0.0.0.0:5050->80/tcp, :::5050->80/tcp                                         pgadmin
```

Stop the containers with
``` shell
$ docker compose down
# To delete all data run:
$ docker compose down -v
```

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

```curl
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

