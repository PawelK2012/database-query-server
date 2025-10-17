package main

import (
	"context"
	"fmt"
	"log"

	"exmple.com/database-query-server/internal/handlers"
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
			//mcp.WithString("query", mcp.Description("DB query")),
			// mcp.WithInputSchema[QueryRequest](),
			// mcp.WithOutputSchema[QueryResponse](),
			mcp.WithString("query", mcp.Description("Search query")),
			mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Max(100)),
			mcp.WithNumber("offset", mcp.DefaultNumber(0), mcp.Min(0)),
		),
		mcp.NewStructuredToolHandler(handleDBQuery),
	)

	// Start StreamableHTTP server
	log.Println("**Starting StreamableHTTP server on :8080")
	httpServer := server.NewStreamableHTTPServer(s)
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

// Handler receives typed input and returns typed output
func handleDBQuery(ctx context.Context, req mcp.CallToolRequest, args QueryRequest) (QueryResponse, error) {
	// Input is already validated and bound to SearchRequest struct
	fmt.Printf("handler - db query")
	limit := args.Limit
	if limit <= 0 {
		limit = 10
	}

	qh := handlers.NewQueryHandler()
	// Perform search logic
	resp, err := qh.SearchUsersInDB(args.Query, args.Format)
	if err != nil {
		log.Fatal("DB quesry failed")
	}
	// Return structured response
	return QueryResponse{
		Query:    resp.Query,
		Response: resp.Response,
		//Format:    json,
	}, nil
}
