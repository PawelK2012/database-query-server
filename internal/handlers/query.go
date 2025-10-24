package handlers

import (
	"context"
	"fmt"
	"log"

	"exmple.com/database-query-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
)

func (qh *QueryHandler) ExecuteQuery(ctx context.Context, req mcp.CallToolRequest, args types.QueryRequest) (*types.QueryResponse, error) {
	// Input is already validated and bound to SearchRequest struct
	limit := args.Limit
	if limit <= 0 {
		limit = 10
	}

	// Placeholder implementation
	fmt.Printf("execute_query handler got query %v with format %v", args.Query, args.Format)

	r, err := qh.repository.Postgress.ExecQuery(ctx, args.Query, args.Parameters)
	if err != nil {
		log.Printf("execute_query %v failed %v", args.Query, err)
		return nil, err
	}
	log.Printf("execute_query response %+v", r)

	response := &types.QueryResponse{
		Query:    "some sql",
		Response: "fix me please", // fix this rubish
	}
	return response, nil
}
