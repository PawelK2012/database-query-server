package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"exmple.com/database-query-server/internal/utils"
	"exmple.com/database-query-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// GetSchema retrieves the schema information for the specified tables and formats response for MCP client
func (qh *QueryHandler) GetSchema(ctx context.Context, req mcp.CallToolRequest, args types.SchemaRequest) (*types.QueryResponse, error) {
	log.Printf("execute GetSchema for tables: %v deatiled: %v", args.Tables, args.Detailed)
	res, err := qh.repository.Postgress.GetSchema(ctx, args.Tables)
	if err != nil {
		return nil, fmt.Errorf("get_schema for table %v failed %v", args.Tables, err)
	}

	jsonRes, err := utils.DataToJson(res)
	if err != nil {
		return nil, fmt.Errorf("failed to encode get_schema response to JSON format %v", err)
	}

	response := &types.QueryResponse{
		Query:    "get_schema",
		Response: jsonRes,
		Format:   "json",
	}
	return response, nil
}

// GetStatus returns the database status, the number of connections, and the timestamp of the last DB access
func (qh *QueryHandler) GetStatus(ctx context.Context, req mcp.CallToolRequest, args types.ConnectionStatus) (*types.ConnectionStatusResp, error) {
	log.Printf("execute GetStatus for DB: %v ", args.Database)
	par := make(map[string]any)
	par["1"] = args.Database
	query := "SELECT numbackends FROM pg_stat_database WHERE datname = $1"

	resp, err := qh.repository.Postgress.ExecQuery(ctx, query, par)
	if err != nil {
		return nil, err
	}
	//casting to int64
	i64, ok := resp[0]["numbackends"].(int64)
	if !ok {
		return nil, fmt.Errorf("failed converting connection pool status")
	}

	parLastPing := make(map[string]any)
	parLastPing["1"] = args.Database

	queryLastPing := "SELECT state_change FROM pg_stat_activity WHERE datname= $1 ORDER BY state_change DESC"
	respLastPing, err := qh.repository.Postgress.ExecQuery(ctx, queryLastPing, parLastPing)
	if err != nil {
		return nil, err
	}

	//casting to time
	st, ok := respLastPing[0]["state_change"].(time.Time)
	if !ok {
		return nil, fmt.Errorf("failed converting last ping status")
	}

	statResp := &types.ConnectionStatusResp{
		Database:  args.Database,
		Connected: true, // if the DB isn’t up, this function will return an error before reaching this point
		PoolStats: int(i64),
		LastPing:  st.Format(time.RFC3339),
	}

	return statResp, nil
}

// ExecuteQuery executes a SQL query and returns the results in the specified format
func (qh *QueryHandler) ExecuteQuery(ctx context.Context, req mcp.CallToolRequest, args types.QueryRequest) (*types.QueryResponse, error) {
	log.Printf("execute_query handler got query %v with format %v", args.Query, args.Format)

	// Input is already validated and bound to SearchRequest struct
	limit := args.Limit
	if limit <= 0 {
		limit = 10
	}

	qResp, err := qh.repository.Postgress.ExecQuery(ctx, args.Query, args.Parameters)
	if err != nil {
		return nil, fmt.Errorf("execute_query %v failed %v", args.Query, err)
	}
	formattedResp, err := formatData(args.Format, qResp)
	if err != nil {
		return nil, err
	}

	response := &types.QueryResponse{
		Query:    args.Query,
		Response: formattedResp,
		Format:   args.Format,
	}
	return response, nil
}

func (qh *QueryHandler) ExecutePrepared(ctx context.Context, req mcp.CallToolRequest, args types.PreparedRequest) (*types.QueryResponse, error) {
	log.Printf("execute_prepared handler got query %v with format %v", args.StatementName, args.Format)
	qResp, err := qh.repository.Postgress.ExecPrepared(ctx, args.StatementName, args.Parameters)
	if err != nil {
		return nil, fmt.Errorf("execute_prepared %v failed %v", args.StatementName, err)
	}
	formattedResp, err := formatData(args.Format, qResp)
	if err != nil {
		return nil, err
	}
	response := &types.QueryResponse{
		Query:    args.StatementName,
		Response: formattedResp,
		Format:   args.Format,
	}
	return response, nil
}

func formatData(format string, data []map[string]interface{}) (string, error) {
	switch format {
	case "json":
		formattedResp, err := utils.DataToJson(data)
		if err != nil {
			return "", fmt.Errorf("failed to encode execute_query response to JSON format %v", err)
		}
		return formattedResp, nil
	case "csv":
		formattedResp, err := utils.DataToCSV(data)
		if err != nil {
			return "", fmt.Errorf("failed to encode execute_query response to CSV format %v", err)
		}
		return formattedResp, nil
	case "table":
		formattedResp, err := utils.DataToHTMLTable(data)
		if err != nil {
			return "", fmt.Errorf("failed to encode execute_query response to HTML table format %v", err)
		}
		return formattedResp, nil
	default:
		return "", fmt.Errorf("format %v not supported", format)
	}

}
