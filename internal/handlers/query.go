package handlers

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"exmple.com/database-query-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// GetSchema retrieves the schema information for the specified tables and formats response for MCP client
func (qh *QueryHandler) GetSchema(ctx context.Context, req mcp.CallToolRequest, args types.SchemaRequest) (*types.QueryResponse, error) {
	fmt.Printf("execute GetSchema for tables: %v deatiled: %v", args.Tables, args.Detailed)
	res, err := qh.repository.Postgress.GetSchema(ctx, args.Tables)
	if err != nil {
		return nil, fmt.Errorf("get_schema for table %v failed %v", args.Tables, err)
	}

	jsonRes, err := dataToJson(res)
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

// ExecuteQuery executes a SQL query and returns the results in the specified format
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
		formattedResp, err = dataToJson(qResp)
		if err != nil {
			return nil, fmt.Errorf("failed to encode execute_query response to JSON format %v", err)
		}
	case "csv":
		formattedResp, err = dataToCSV(qResp)
		if err != nil {
			return nil, fmt.Errorf("failed to encode execute_query response to CSV format %v", err)
		}
	case "table":
		formattedResp, err = dataHTMLTable(qResp)
		if err != nil {
			return nil, fmt.Errorf("failed to encode execute_query response to HTML table format %v", err)
		}
	default:
		return nil, fmt.Errorf("format %v not supported", args.Format)
	}

	response := &types.QueryResponse{
		Query:    args.Query,
		Response: formattedResp,
		Format:   args.Format,
	}
	return response, nil
}

func (qh *QueryHandler) ExecutePrepared(ctx context.Context, req mcp.CallToolRequest, args types.PreparedRequest) (*types.QueryResponse, error) {
	// Input is already validated and bound to SearchRequest struct
	// limit := args.Limit
	// if limit <= 0 {
	// 	limit = 10
	// }

	fmt.Printf("execute_prepared handler got query %v with format %v", args.StatementName, args.Format)
	qResp, err := qh.repository.Postgress.ExecPrepared(ctx, args.StatementName, args.Parameters)
	if err != nil {
		return nil, fmt.Errorf("execute_prepared %v failed %v", args.StatementName, err)
	}
	var formattedResp string
	switch args.Format {
	case "json":
		formattedResp, err = dataToJson(qResp)
		if err != nil {
			return nil, fmt.Errorf("failed to encode execute_query response to JSON format %v", err)
		}
	case "csv":
		formattedResp, err = dataToCSV(qResp)
		if err != nil {
			return nil, fmt.Errorf("failed to encode execute_query response to CSV format %v", err)
		}
	case "table":
		formattedResp, err = dataHTMLTable(qResp)
		if err != nil {
			return nil, fmt.Errorf("failed to encode execute_query response to HTML table format %v", err)
		}
	default:
		return nil, fmt.Errorf("format %v not supported", args.Format)
	}

	response := &types.QueryResponse{
		Query:    args.StatementName,
		Response: formattedResp,
		Format:   args.Format,
	}
	return response, nil
}

// dataToJson converts a slice of maps containing data into a JSON string
func dataToJson(data []map[string]interface{}) (string, error) {
	enco, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(enco), nil
}

// dataToCSV converts a slice of maps into a CSV string
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

// dataHTMLTable converts a slice of maps into an HTML table string.
// - rows: each map represents one table row (key -> cell value).
// - Columns are the union of all map keys, sorted alphabetically for deterministic output.
func dataHTMLTable(rows []map[string]interface{}) (string, error) {
	var b strings.Builder

	// start table
	b.WriteString("<table>")

	// no rows - return error
	if len(rows) == 0 {
		return "", fmt.Errorf("no data to convert")
	}

	// collect all column names
	colSet := make(map[string]struct{})
	for _, r := range rows {
		for k := range r {
			colSet[k] = struct{}{}
		}
	}

	cols := make([]string, 0, len(colSet))
	for k := range colSet {
		cols = append(cols, k)
	}
	sort.Strings(cols) // deterministic order

	// header
	b.WriteString("<thead><tr>")
	for _, c := range cols {
		b.WriteString("<th>")
		b.WriteString(html.EscapeString(c))
		b.WriteString("</th>")
	}
	b.WriteString("</tr></thead>")

	// body
	b.WriteString("<tbody>")
	for _, r := range rows {
		b.WriteString("<tr>")
		for _, c := range cols {
			b.WriteString("<td>")
			val, ok := r[c]
			if ok && val != nil {
				// convert value to string and escape HTML
				b.WriteString(html.EscapeString(fmt.Sprint(val)))
			}
			// if missing or nil -> empty cell
			b.WriteString("</td>")
		}
		b.WriteString("</tr>")
	}
	b.WriteString("</tbody>")

	// end table
	b.WriteString("</table>")

	return b.String(), nil
}
