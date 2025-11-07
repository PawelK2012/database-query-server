package database_test

import (
	"context"
	"regexp"
	"testing"

	"exmple.com/database-query-server/internal/database"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TODO All these test cases should be refactored into one single GO table tests
func TestExecQuery_Happy_Path(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	defer db.Close()

	pg := &database.Postgress{Pg: db}

	ctx := context.Background()
	query := "SELECT CustomerName, Address FROM customers WHERE Country =$1;"

	rows := sqlmock.NewRows([]string{"CustomerName", "Address"}).
		AddRow("Alice", "Some address 12")

	// query params
	params := make(map[string]any)
	params["1"] = "UK"

	expected := make(map[string]any)
	expected["CustomerName"] = "Alice"
	expected["Address"] = "Some address 12"

	mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs("UK").WillReturnRows(rows)

	result, err := pg.ExecQuery(ctx, query, params)
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	assert.Equal(t, 1, len(result), "test should return 1 row")
	for i := range result {
		assert.Equal(t, expected, result[i], "returned data should match expectation")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestExecQuery_Happy_Path_With_2_Args(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	defer db.Close()

	pg := &database.Postgress{Pg: db}

	ctx := context.Background()
	query := "SELECT CustomerName, Address FROM customers WHERE Country =$1 AND City LIKE $2"

	// expected rows to return
	rows := sqlmock.NewRows([]string{"CustomerName", "Address"}).
		AddRow("Alice", "Alice Address").
		AddRow("Bob", "Some address 12")

	// query params
	par := make(map[string]any)
	par["1"] = "UK"
	par["2"] = "L%"

	expected := []map[string]any{
		{"CustomerName": "Alice", "Address": "Alice Address"},
		{"CustomerName": "Bob", "Address": "Some address 12"},
	}

	mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs("UK", "L%").WillReturnRows(rows)

	result, err := pg.ExecQuery(ctx, query, par)
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}

	assert.Equal(t, 2, len(result), "test should return 2 rows")
	for i := range result {
		assert.Equal(t, expected[i], result[i], "returned data should match expectation")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestExecQuery_Select_All(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}

	defer db.Close()

	pg := &database.Postgress{Pg: db}

	ctx := context.Background()
	query := "select * from customers;"

	// expected rows to return
	rows := sqlmock.NewRows([]string{"CustomerName", "ContactName", "Address", "City", "PostalCode", "Country"}).
		AddRow("Alice", "Alice contact name", "Alice Address", "Liverpool", "xd12", "UK").
		AddRow("Bob", "Bob contact name", "Bob Address", "Leeds", "bbb23", "UK")

	expected := []map[string]any{
		{"CustomerName": "Alice", "ContactName": "Alice contact name", "Address": "Alice Address", "City": "Liverpool", "PostalCode": "xd12", "Country": "UK"},
		{"CustomerName": "Bob", "ContactName": "Bob contact name", "Address": "Bob Address", "City": "Leeds", "PostalCode": "bbb23", "Country": "UK"},
	}

	// query params
	params := make(map[string]any)

	mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs().WillReturnRows(rows)

	result, err := pg.ExecQuery(ctx, query, params)
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}

	assert.Equal(t, 2, len(result), "test should return 2 rows")
	for i := range result {
		assert.Equal(t, expected[i], result[i], "returned data should match expectation")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
