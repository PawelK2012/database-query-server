package handlers

import (
	"context"
	"fmt"

	"exmple.com/database-query-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
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

func (qh *QueryHandler) SearchUsersInDB(ctx context.Context, req mcp.CallToolRequest, args QueryRequest) (*types.QueryResponse, error) {
	// Input is already validated and bound to SearchRequest struct
	fmt.Printf("queru SearchUsersInDB")
	limit := args.Limit
	if limit <= 0 {
		limit = 10
	}

	// Placeholder implementation
	fmt.Printf("handler got query %v with format %v", args.Query, args.Format)
	r := "handler response"
	qh.repository.Postgress.Init(ctx)
	// qh.repo.Init(ctx)

	response := &types.QueryResponse{
		Query:    "some sql",
		Response: r,
	}
	return response, nil
}
