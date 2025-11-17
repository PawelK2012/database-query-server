package handlers_test

import (
	"context"
	"testing"

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

	tests := []struct {
		name string
		// repository *repository.Repository
		req       mcp.CallToolRequest
		args      types.QueryRequest
		tableMock []map[string]interface{}
		want      *types.QueryResponse
		wantErr   bool
	}{
		{name: "Happy Flow", req: request, args: reqArgs, tableMock: mtbl, want: &expected, wantErr: false},
		// {name: "Sad Flow", repository: repo, req: request, args: reqArgs, want: &expected, wantErr: true},
		{name: "Happy Flow - to CSV export", req: request, args: reqArgsCVS, tableMock: mtbl, want: &expectedCSVOutput, wantErr: false},
		{name: "Happy Flow - to CSV export with different Postgres types", req: request, args: reqArgsCVS, tableMock: mtblWithDifferentTypes, want: &expectedCSVOutputWithDifferentTypes, wantErr: false},
		// {name: "Happy Flow - various types", repository: repo, req: request, args: reqArgsCVS, want: &expectedCSVOutput, wantErr: false},
		{name: "Fail - invalid format", req: request, args: reqArgsInvalidFormat, tableMock: mtbl, want: &expectedCSVOutput, wantErr: true},
		{name: "Happy Flow - export to HTML table", req: request, args: reqToTable, tableMock: mtbl, want: &expectedTableOutput, wantErr: false},
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
		{name: "Happy Flow", req: request, args: reqArgs, tableMock: mtbl, want: &expected, wantErr: false},
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
