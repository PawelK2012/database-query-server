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

	item1 := make(map[string]interface{})
	item1["id"] = 22
	item1["CustomerName"] = "John"
	item1["ContactName"] = true
	item1["Address"] = "Some street in Dublin"
	item1["City"] = "Dublin"
	item1["PostalCode"] = "1dbld12"
	item1["Country"] = "Ireland"

	mtbl = append(mtbl, item, item1)

	pg, _ := database.NewPostgresClientMock(mtbl, false)

	repo := &repository.Repository{Postgress: pg}

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
	}
	expectedCSVOutput := types.QueryResponse{
		Query:    "SELECT * FROM customers",
		Response: "Address,City,ContactName,Country,CustomerName,PostalCode,id\nSome street in London,London,Bob mum,UK,Bob,1ld12,1\n",
	}

	tests := []struct {
		name       string
		repository *repository.Repository
		req        mcp.CallToolRequest
		args       types.QueryRequest
		want       *types.QueryResponse
		wantErr    bool
	}{
		{name: "Happy Flow", repository: repo, req: request, args: reqArgs, want: &expected, wantErr: false},
		// {name: "Sad Flow", repository: repo, req: request, args: reqArgs, want: &expected, wantErr: true},
		{name: "Happy Flow - to CSV export", repository: repo, req: request, args: reqArgsCVS, want: &expectedCSVOutput, wantErr: false},
		// {name: "Happy Flow - various types", repository: repo, req: request, args: reqArgsCVS, want: &expectedCSVOutput, wantErr: false},
		{name: "Fail - invalid format", repository: repo, req: request, args: reqArgsInvalidFormat, want: &expectedCSVOutput, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qh := handlers.NewQueryHandler(tt.repository)
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
