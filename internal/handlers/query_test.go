package handlers_test

import (
	"context"
	"testing"

	"exmple.com/database-query-server/internal/database"
	"exmple.com/database-query-server/internal/handlers"
	"exmple.com/database-query-server/internal/repository"
	"exmple.com/database-query-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestQueryHandler_ExecuteQuery(t *testing.T) {

	//mocked SQL table
	var mtbl []map[string]interface{}
	item := make(map[string]interface{})
	item["row1"] = "sss"
	mtbl = append(mtbl, item)

	pg, _ := database.NewPostgresClientMock(mtbl, false)

	repo := &repository.Repository{Postgress: pg}

	reqArgs := types.QueryRequest{
		Database: "postgres",
		Query:    "SELECT * FROM customers",
		Format:   "json",
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
		Response: "xxxx",
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		repository *repository.Repository
		// Named input parameters for target function.
		req     mcp.CallToolRequest
		args    types.QueryRequest
		want    *types.QueryResponse
		wantErr bool
	}{
		{name: "Happy Flow", repository: repo, req: request, args: reqArgs, want: &expected, wantErr: false},
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
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("ExecuteQuery() SSS = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestPrepareJsonResp(t *testing.T) {
// 	tests := []struct {
// 		name string // description of this test case
// 		// Named input parameters for target function.
// 		d       []map[string]interface{}
// 		want    string
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, gotErr := handlers.prepareJsonResp(tt.d)
// 			if gotErr != nil {
// 				if !tt.wantErr {
// 					t.Errorf("PrepareJsonResp() failed: %v", gotErr)
// 				}
// 				return
// 			}
// 			if tt.wantErr {
// 				t.Fatal("PrepareJsonResp() succeeded unexpectedly")
// 			}
// 			// TODO: update the condition below to compare got with tt.want.
// 			if true {
// 				t.Errorf("PrepareJsonResp() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
