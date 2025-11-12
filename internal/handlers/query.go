package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

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
		formattedResp, err = qh.dataToJson(qResp)
		if err != nil {
			//formattedResp = fmt.Sprintf("encoding execute_query failed %v", err)
			return nil, fmt.Errorf("encoding execute_query failed %v", err)
		}
	case "csv":
		fmt.Println("encoding to CSV")
		qh.dataToCSV(qResp)
	case "table":
		fmt.Println("encoding to table")
	}
	//add default to switch

	response := &types.QueryResponse{
		Query:    args.Query,
		Response: formattedResp,
	}
	return response, nil
}

func (qh *QueryHandler) dataToJson(d []map[string]interface{}) (string, error) {
	enco, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(enco), nil
}

func (qh *QueryHandler) dataToCSV(d []map[string]interface{}) ([][]string, error) {
	var records [][]string
	var headers []string
	var val []string

	for k, v := range d {
		// The system should capture headers only during the initial run
		if k == 0 {
			headers = qh.getCSVHeaders(v)
		}
		for _, y := range v {
			enco, err := json.Marshal(y)
			if err != nil {
				return nil, err
			}
			val = append(val, string(enco))
		}
	}
	//headers = append(headers, val...)

	records = append(records, headers)
	records = append(records, val)

	fmt.Println("0000")
	w := csv.NewWriter(os.Stdout)
	w.WriteAll(records) // calls Flush internally
	fmt.Println("0000")
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

	fmt.Println("11*********")
	fmt.Printf("== finalllllllll CSSSVVVV %v", records)
	fmt.Println("22*********")
	return records, nil
}

func (qh *QueryHandler) getCSVHeaders(d map[string]interface{}) []string {
	var headers []string
	for header := range d {
		headers = append(headers, header)
	}

	return headers
}
