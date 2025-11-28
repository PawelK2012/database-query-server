package handlers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"exmple.com/database-query-server/internal/database"
	"exmple.com/database-query-server/internal/handlers"
	"exmple.com/database-query-server/internal/repository"
	"exmple.com/database-query-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func TestQueryHandler_ExecuteQuery(t *testing.T) {

	//mocked table rows
	var mtbl []map[string]interface{}
	item := make(map[string]interface{})
	item["id"] = "1"
	item["CustomerName"] = "Bob"
	item["ContactName"] = "Bob mum"
	item["Address"] = "Some street in London"
	item["City"] = "London"
	item["PostalCode"] = "1ld12"
	item["Country"] = "UK"
	mtbl = append(mtbl, item)

	var mtblWithDifferentTypes []map[string]interface{}
	item1 := make(map[string]interface{})
	item1["id"] = 22
	item1["price"] = 123.78
	item1["ContactName"] = true
	item1["Address"] = nil
	item1["City"] = "Dublin"
	item1["PostalCode"] = "1dbld12"
	item1["Country"] = "Ireland"
	mtblWithDifferentTypes = append(mtblWithDifferentTypes, item1)

	reqArgs := types.QueryRequest{
		Database: "postgres",
		Query:    "SELECT * FROM customers",
		Format:   "json",
	}

	reqArgsCVS := types.QueryRequest{
		Database: "postgres",
		Query:    "SELECT * FROM customers",
		Format:   "csv",
	}

	reqArgsInvalidFormat := types.QueryRequest{
		Database: "postgres",
		Query:    "SELECT * FROM customers",
		Format:   "invalid",
	}

	reqToTable := types.QueryRequest{
		Database: "postgres",
		Query:    "SELECT * FROM customers",
		Format:   "table",
	}

	reqInvalidQuery := types.QueryRequest{
		Database: "postgres",
		Query:    "SELECT_INJ SELECT * FROM customers",
		Format:   "table",
	}

	args := make(map[string]interface{})
	args["databse"] = "sql"
	args["query"] = "SELECT * FROM customers"
	args["format"] = "json"

	par := make(map[string]any)
	par["1"] = "UK"

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "execute_query",
			Arguments: args,
		},
	}

	expected := types.QueryResponse{
		Query:    "SELECT * FROM customers",
		Response: `[{"Address":"Some street in London","City":"London","ContactName":"Bob mum","Country":"UK","CustomerName":"Bob","PostalCode":"1ld12","id":"1"}]`,
		Format:   "json",
	}
	expectedCSVOutput := types.QueryResponse{
		Query:    "SELECT * FROM customers",
		Response: "Address,City,ContactName,Country,CustomerName,PostalCode,id\nSome street in London,London,Bob mum,UK,Bob,1ld12,1\n",
		Format:   "csv",
	}

	expectedCSVOutputWithDifferentTypes := types.QueryResponse{
		Query:    "SELECT * FROM customers",
		Response: "Address,City,ContactName,Country,PostalCode,id,price\n,Dublin,true,Ireland,1dbld12,22,123.78\n",
		Format:   "csv",
	}

	expectedTableOutput := types.QueryResponse{
		Query:    "SELECT * FROM customers",
		Response: "<table><thead><tr><th>Address</th><th>City</th><th>ContactName</th><th>Country</th><th>CustomerName</th><th>PostalCode</th><th>id</th></tr></thead><tbody><tr><td>Some street in London</td><td>London</td><td>Bob mum</td><td>UK</td><td>Bob</td><td>1ld12</td><td>1</td></tr></tbody></table>",
		Format:   "table",
	}

	expectedInvalidQueryErr := types.QueryResponse{
		Query:    "SELECT * FROM customers",
		Response: "---",
		Format:   "table",
	}

	tests := []struct {
		name      string
		req       mcp.CallToolRequest
		args      types.QueryRequest
		tableMock []map[string]interface{}
		want      *types.QueryResponse
		wantErr   bool
	}{
		{name: "Happy Flow execute_query", req: request, args: reqArgs, tableMock: mtbl, want: &expected, wantErr: false},
		{name: "Happy Flow execute_query - to CSV export", req: request, args: reqArgsCVS, tableMock: mtbl, want: &expectedCSVOutput, wantErr: false},
		{name: "Happy Flow execute_query - to CSV export with different Postgres types", req: request, args: reqArgsCVS, tableMock: mtblWithDifferentTypes, want: &expectedCSVOutputWithDifferentTypes, wantErr: false},
		{name: "Fail execute_query - invalid format", req: request, args: reqArgsInvalidFormat, tableMock: mtbl, want: &expectedCSVOutput, wantErr: true},
		{name: "Happy Flow execute_query - export to HTML table", req: request, args: reqToTable, tableMock: mtbl, want: &expectedTableOutput, wantErr: false},
		{name: "Fail execute_query - query must start with SELECT statement", req: request, args: reqInvalidQuery, tableMock: mtbl, want: &expectedInvalidQueryErr, wantErr: true},

		// Add execute_query test
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg, _ := database.NewPostgresClientMock(tt.tableMock, false)
			repo := &repository.Repository{Postgress: pg}
			qh := handlers.NewQueryHandler(repo)
			got, gotErr := qh.ExecuteQuery(context.Background(), tt.req, tt.args)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExecuteQuery() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ExecuteQuery() succeeded unexpectedly")
			}
			if true {
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestQueryHandler_GetSchema(t *testing.T) {
	//mocked table rows
	var mtbl []map[string]interface{}
	row2 := make(map[string]interface{})
	row2["column_name"] = "customername"
	row2["data_type"] = "character varying"
	row2["character_maximum_length"] = int64(200)

	mtbl = append(mtbl, row2)

	var tbls []string
	tbls = append(tbls, "customers")
	reqArgs := types.SchemaRequest{
		Database: "postgres",
		Tables:   tbls,
		Detailed: true,
	}
	args := make(map[string]interface{})
	args["database"] = "postgres"
	args["tables"] = tbls
	args["detailed"] = true

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "get_schema",
			Arguments: args,
		},
	}

	expected := types.QueryResponse{
		Query:    "get_schema",
		Response: `[{"character_maximum_length":200,"column_name":"customername","data_type":"character varying"}]`,
		Format:   "json",
	}
	tests := []struct {
		name       string
		repository *repository.Repository
		req        mcp.CallToolRequest
		args       types.SchemaRequest
		tableMock  []map[string]interface{}
		want       *types.QueryResponse
		wantErr    bool
	}{
		{name: "Happy Flow - GetSchema", req: request, args: reqArgs, tableMock: mtbl, want: &expected, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg, _ := database.NewPostgresClientMock(tt.tableMock, false)
			repo := &repository.Repository{Postgress: pg}
			qh := handlers.NewQueryHandler(repo)
			got, gotErr := qh.GetSchema(context.Background(), tt.req, tt.args)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetSchema() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetSchema() succeeded unexpectedly")
			}
			if true {
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestQueryHandler_GetStatus(t *testing.T) {
	timNow := time.Now()
	var mtbl []map[string]interface{}
	row2 := make(map[string]interface{})
	row2["numbackends"] = int64(20)
	row2["state_change"] = timNow
	row2["character_maximum_length"] = int64(200)
	mtbl = append(mtbl, row2)

	reqArgs := types.ConnectionStatus{
		Database: "mcp-db",
	}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "get_connection_status",
			Arguments: reqArgs,
		},
	}

	expected := types.ConnectionStatusResp{
		Database:  "mcp-db",
		Connected: true,
		PoolStats: 20,
		LastPing:  timNow.Format(time.RFC3339),
	}

	reqArgsEmpty := types.ConnectionStatus{
		Database: "",
	}

	requestFail := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "get_connection_status",
			Arguments: reqArgsEmpty,
		},
	}

	var mtbl1 []map[string]interface{}
	rowErr := make(map[string]interface{})
	rowErr["numbackends"] = 200
	rowErr["state_change"] = timNow
	rowErr["character_maximum_length"] = int64(200)
	mtbl1 = append(mtbl1, rowErr)

	expectedErr := fmt.Errorf("failed converting connection pool status")

	tests := []struct {
		name       string
		repository *repository.Repository
		tableMock  []map[string]interface{}
		req        mcp.CallToolRequest
		args       types.ConnectionStatus
		want       *types.ConnectionStatusResp
		wantErr    bool
	}{
		{name: "Happy Flow - GetStatus", req: request, args: reqArgs, tableMock: mtbl, want: &expected, wantErr: false},
		{name: "Sad Flow - GetStatus connection pool status error", req: requestFail, args: reqArgsEmpty, tableMock: mtbl1, want: &expected, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg, _ := database.NewPostgresClientMock(tt.tableMock, false)
			repo := &repository.Repository{Postgress: pg}
			qh := handlers.NewQueryHandler(repo)
			got, gotErr := qh.GetStatus(context.Background(), tt.req, tt.args)
			if gotErr != nil {
				if !tt.wantErr {
					assert.EqualValues(t, expectedErr, gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetStatus() succeeded unexpectedly")
			}
			if true {
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestQueryHandler_ExecutePrepared(t *testing.T) {
	query := "SELECT CustomerName, Address FROM customers WHERE Country =$1 AND City LIKE $2"
	var mtblMock []map[string]interface{}
	item := make(map[string]interface{})
	item["id"] = "1"
	item["CustomerName"] = "Bob"
	item["ContactName"] = "Bob mum"
	item["Address"] = "Some street in London"
	item["City"] = "London"
	item["PostalCode"] = "1ld12"
	item["Country"] = "UK"

	item1 := make(map[string]interface{})
	item1["id"] = "2"
	item1["CustomerName"] = "Joe"
	item1["ContactName"] = "Joe Mc Cormack"
	item1["Address"] = "Vikings st"
	item1["City"] = "Liverpool"
	item1["PostalCode"] = "1liv122"
	item1["Country"] = "UK"

	item2 := make(map[string]interface{})
	item2["id"] = "3"
	item2["CustomerName"] = "Anna"
	item2["ContactName"] = "Anna Valerina"
	item2["Address"] = "Near me st"
	item2["City"] = "Dublin"
	item2["PostalCode"] = "Dub02"
	item2["Country"] = "IE"

	mtblMock = append(mtblMock, item, item1, item2)

	par := []any{
		"UK", "L%",
	}

	args := make(map[string]interface{})
	args["databse"] = "mcp-db"
	args["StatementName"] = "SELECT CustomerName, Address FROM customers WHERE Country =$1 AND City LIKE $2"
	args["format"] = "json"
	args["parameters"] = par

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "execute_prepared",
			Arguments: args,
		},
	}

	reqArgs := types.PreparedRequest{
		Database:      "mcp-db",
		StatementName: query,
		Parameters:    par,
		Format:        "json",
	}

	reqArgsCSV := types.PreparedRequest{
		Database:      "mcp-db",
		StatementName: query,
		Parameters:    par,
		Format:        "csv",
	}

	reqArgsTable := types.PreparedRequest{
		Database:      "mcp-db",
		StatementName: query,
		Parameters:    par,
		Format:        "table",
	}

	expected := types.QueryResponse{
		Query:    query,
		Response: `[{"Address":"Some street in London","City":"London","ContactName":"Bob mum","Country":"UK","CustomerName":"Bob","PostalCode":"1ld12","id":"1"}]`,
		Format:   "json",
	}

	expectedCSV := types.QueryResponse{
		Query:    query,
		Response: "Address,City,ContactName,Country,CustomerName,PostalCode,id\nSome street in London,London,Bob mum,UK,Bob,1ld12,1\n",
		Format:   "csv",
	}

	expectedTable := types.QueryResponse{
		Query:    query,
		Response: "<table><thead><tr><th>Address</th><th>City</th><th>ContactName</th><th>Country</th><th>CustomerName</th><th>PostalCode</th><th>id</th></tr></thead><tbody><tr><td>Some street in London</td><td>London</td><td>Bob mum</td><td>UK</td><td>Bob</td><td>1ld12</td><td>1</td></tr></tbody></table>",
		Format:   "table",
	}

	tests := []struct {
		name       string
		repository *repository.Repository
		req        mcp.CallToolRequest
		args       types.PreparedRequest
		tableMock  []map[string]interface{}
		want       *types.QueryResponse
		wantErr    bool
	}{
		{name: "Happy Flow execute_query - expect JSON format", req: request, args: reqArgs, tableMock: mtblMock, want: &expected, wantErr: false},
		{name: "Happy Flow execute_query - expect CSV format", req: request, args: reqArgsCSV, tableMock: mtblMock, want: &expectedCSV, wantErr: false},
		{name: "Happy Flow execute_query - expect HTML table format", req: request, args: reqArgsTable, tableMock: mtblMock, want: &expectedTable, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg, _ := database.NewPostgresClientMock(tt.tableMock, false)
			repo := &repository.Repository{Postgress: pg}
			qh := handlers.NewQueryHandler(repo)
			got, gotErr := qh.ExecutePrepared(context.Background(), tt.req, tt.args)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExecutePrepared() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ExecutePrepared() succeeded unexpectedly")
			}
			if true {
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}
