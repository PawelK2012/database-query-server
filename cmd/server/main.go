package main

import (
	"context"
	"log"
	"os"

	"exmple.com/database-query-server/internal/handlers"
	"exmple.com/database-query-server/internal/repository"
	"exmple.com/database-query-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Instantiate repository to open the Redis connection and create the Cloudant client
	repository, err := repository.New()
	if err != nil {
		log.Fatal(context.Background(), "Can not instantiate repository package, error: %v", err)
		os.Exit(1)
	}

	qh := handlers.NewQueryHandler(repository)

	s := server.NewMCPServer("**StreamableHTTP API Server", "1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),  // Enable MCP protocol logging
		server.WithRecovery(), // Recover from panics in handlers
	)

	s.AddTool(
		mcp.NewTool("execute_query",
			mcp.WithDescription("Demonstrate secure database operations via MCP"),
			mcp.WithTitleAnnotation("Execute DB Query"),
			mcp.WithString("query", mcp.Description("DB query")),
			mcp.WithInputSchema[types.QueryRequest](),
			mcp.WithOutputSchema[types.QueryResponse](),
		),
		mcp.NewStructuredToolHandler(qh.ExecuteQuery),
	)

	s.AddTool(
		mcp.NewTool("get_schema",
			mcp.WithDescription("Retrieve database schema information"),
			mcp.WithTitleAnnotation("Execute get_schema operations"),
			mcp.WithString("schema", mcp.Description("DB schema")),
			mcp.WithInputSchema[types.SchemaRequest](),
			mcp.WithOutputSchema[types.QueryResponse](),
		),
		mcp.NewStructuredToolHandler(qh.GetSchema),
	)

	// Start StreamableHTTP server
	log.Println("**Starting StreamableHTTP server on :8080")
	httpServer := server.NewStreamableHTTPServer(s)
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
