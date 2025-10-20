package main

import (
	"context"
	"log"
	"os"

	"exmple.com/database-query-server/internal/handlers"
	"exmple.com/database-query-server/internal/repository"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type QueryRequest struct {
	Database   string         `json:"database"`
	Query      string         `json:"query"`
	Parameters map[string]any `json:"parameters,omitempty"`
	Format     string         `json:"format,omitempty"` // json, csv, table
	Limit      int            `json:"limit,omitempty"`
	Timeout    int            `json:"timeout,omitempty"`
}

type QueryResponse struct {
	Query    string `json:"query"`
	Response string `json:"response"`
	Format   string `json:"format,omitempty"` // json, csv, table
}

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
			mcp.WithInputSchema[QueryRequest](),
			mcp.WithOutputSchema[QueryResponse](),
		),
		mcp.NewStructuredToolHandler(qh.SearchUsersInDB),
	)

	// Start StreamableHTTP server
	log.Println("**Starting StreamableHTTP server on :8080")
	httpServer := server.NewStreamableHTTPServer(s)
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
