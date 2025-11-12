package handlers

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"

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
	qResp, err := qh.repository.Postgress.ExecQuery(ctx, args.Query, args.Parameters)
	if err != nil {
		return nil, fmt.Errorf("execute_query %v failed %v", args.Query, err)
	}
	var formattedResp string
	switch args.Format {
	case "json":
		fmt.Println("encoding to JSON")
		formattedResp, err = dataToJson(qResp)
		if err != nil {
			return nil, fmt.Errorf("failed to encode execute_query response to JSON format %v", err)
		}
	case "csv":
		fmt.Println("encoding to CSV")
		formattedResp, err = dataToCSV(qResp)
		if err != nil {
			return nil, fmt.Errorf("failed to encode execute_query response to CSV format %v", err)
		}
	case "table":
		fmt.Println("encoding to table")
	default:
		return nil, fmt.Errorf("format %v not supported", args.Format)
	}

	response := &types.QueryResponse{
		Query:    args.Query,
		Response: formattedResp,
	}
	return response, nil
}

func dataToJson(data []map[string]interface{}) (string, error) {
	enco, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(enco), nil
}

func dataToCSV(data []map[string]interface{}) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("no data to convert")
	}

	// Extract headers from the first map
	headers := make([]string, 0, len(data[0]))
	for key := range data[0] {
		headers = append(headers, key)
	}
	// Map iteration order is intentionally randomized, so we use sorting for consistency
	// See https://go.dev/blog/maps#iteration-order
	sort.Strings(headers)

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write headers
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write headers: %w", err)
	}

	// Write rows
	for _, row := range data {
		record := make([]string, len(headers))
		for i, header := range headers {
			if val, ok := row[header]; ok && val != nil {
				switch v := val.(type) {
				case string:
					record[i] = v
				case fmt.Stringer:
					record[i] = v.String()
				case int, int8, int16, int32, int64:
					record[i] = fmt.Sprintf("%d", v)
				case float32, float64:
					record[i] = strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64)
				case bool:
					record[i] = strconv.FormatBool(v)
				default:
					record[i] = fmt.Sprintf("%v", v)
				}
			}
		}
		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("csv write error: %w", err)
	}

	return buf.String(), nil
}
