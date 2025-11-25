# database-query-server

[![Go](https://github.com/PawelK2012/database-query-server/actions/workflows/go.yml/badge.svg)](https://github.com/PawelK2012/database-query-server/actions/workflows/go.yml)

[![Go Coverage](https://github.com/PawelK2012/database-query-server/wiki/coverage.svg)](https://raw.githack.com/wiki/PawelK2012/database-query-server/coverage.html)

Sample MCP Server - Go (database-query-server) - built according to these [specifications](https://github.com/IBM/mcp-context-forge/issues/897) 


# How to run project

1. Simple start with `docker compose up` - Related [documentation](https://github.com/docker/awesome-compose/tree/master/postgresql-pgadmin) 
2. Start MPC server with command `make run`
3. You can use the Go [MCP-client](https://github.com/PawelK2012/mcp-client), which includes `tools/call` examples 

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
## Usage Examples

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
### Execute prepared statements safely
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "execute_prepared",
    "arguments": {
      "database": "primary",
      "StatementName": "SELECT id, name, email FROM users WHERE active = $1",
      "parameters": {"1": true},
      "format": "json"
    }
  }
}
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
### Get ConnectionStatus

```json
echo '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call", 
  "params": {
    "name": "get_connection_status",
    "arguments": {
      "database": "primary"//name of the DB you want to retrieve stats for
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

# Populate the Database with Test Data

## Postgress

You can easily populate your PostgreSQL database with test data by calling this MCP server using the following SQL queries:

1. Create `Customers` table

```
CREATE TABLE IF NOT EXISTS Customers (
		id SERIAL PRIMARY KEY,
		CustomerName VARCHAR(200),
		ContactName VARCHAR(250),
		Address VARCHAR(500),
		City VARCHAR(250),
		PostalCode VARCHAR(150),
		Country VARCHAR(250),
		created_at TIMESTAMP
	)
```
2. Update table with data

```
INSERT INTO Customers (CustomerName, ContactName, Address, City, PostalCode, Country)
			VALUES
			('Cardinal', 'Tom B. Erichsen', 'Skagen 21', 'Stavanger', '4006', 'Norway'),
			('Greasy Burger', 'Per Olsen', 'Gateveien 15', 'Sandnes', '4306', 'Norway'),
			('Tasty Tee', 'Finn Egan', 'Streetroad 19B', 'Liverpool', 'L1 0AA', 'UK');
```
3. Query DB

```
SELECT * FROM customers
```

or  (don't forget to provide required Parameters)

```
SELECT CustomerName, Address FROM customers WHERE Country =$1 AND City LIKE $2
```

Example

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "execute_query",
    "arguments": {
      "database": "primary",
      "query": "SELECT CustomerName, Address FROM customers WHERE Country =$1 AND City LIKE $2",
      "parameters": {"1": "UK", "2": "L%"},
      "format": "json",
      "limit": 100
    }
  }
}
```