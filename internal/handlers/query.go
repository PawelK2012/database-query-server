package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"exmple.com/database-query-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
)

func (qh *QueryHandler) ExecuteQuery(ctx context.Context, req mcp.CallToolRequest, args types.QueryRequest) (*types.QueryResponse, error) {
	// Input is already validated and bound to SearchRequest struct
	limit := args.Limit
	if limit <= 0 {
		limit = 10
	}

	fmt.Printf("execute_query handler got query %v with format %v", args.Query, args.Format)

	r, err := qh.repository.Postgress.ExecQuery(ctx, args.Query, args.Parameters)
	if err != nil {
		return nil, fmt.Errorf("execute_query %v failed %v", args.Query, err)
	}

	var formattedResp string
	switch args.Format {
	case "json":
		formattedResp, err = prepareJsonResp(r)
		if err != nil {
			//formattedResp = fmt.Sprintf("encoding execute_query failed %v", err)
			return nil, fmt.Errorf("encoding execute_query failed %v", err)
		}
	case "csv":
		fmt.Println("encoding to CSV")
	case "table":
		fmt.Println("encoding to table")
	}

	response := &types.QueryResponse{
		Query:    args.Query,
		Response: formattedResp,
	}
	return response, nil
}

func prepareJsonResp(d []map[string]interface{}) (string, error) {
	enco, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(enco), nil
}
