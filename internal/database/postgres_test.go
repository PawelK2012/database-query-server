package database_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"exmple.com/database-query-server/internal/database"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestExecQuery_Happy_Path_Select_All(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}

	defer db.Close()

	pg := &database.Postgress{Pg: db}

	ctx := context.Background()
	query := "SELECT * FROM users;"

	// expected rows to return
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Alice").
		AddRow(2, "Bob")

	// query params
	par := make(map[string]any)
	par["1"] = "Alice"
	par["2"] = "Bob"

	var expected []map[string]interface{}
	row := make(map[string]interface{})
	row["id"] = int64(1)
	row["name"] = "Alice"

	row2 := make(map[string]interface{})
	row2["id"] = int64(2)
	row2["name"] = "Bob"

	expected = append(expected, row, row2)

	mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs().WillReturnRows(rows)

	result, err := pg.ExecQuery(ctx, query, par)
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.EqualValues(t, expected, result)
}

func TestExecQuery_Sad_Path(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	defer db.Close()

	pg := &database.Postgress{Pg: db}

	ctx := context.Background()
	query := "SELECT * FROM users"

	expected := fmt.Errorf("some error")

	// query params
	par := make(map[string]any)
	par["1"] = "Alice"

	mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs().WillReturnError(fmt.Errorf("some error"))

	_, err = pg.ExecQuery(ctx, query, par)
	if err != nil {
		assert.EqualValues(t, expected, err)
		return
	}
}

func TestExecQuery_Happy_Path_Query_Will_Not_Return_Rows(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	defer db.Close()

	pg := &database.Postgress{Pg: db}

	ctx := context.Background()
	query := "CREATE SCHEMA new_schema"

	rows := sqlmock.NewRows([]string{})

	var expected []map[string]interface{}
	row := make(map[string]interface{})
	row["message"] = "success"

	expected = append(expected, row)

	// query params
	par := make(map[string]any)
	mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WillReturnRows(rows)

	result, err := pg.ExecQuery(ctx, query, par)
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.EqualValues(t, expected, result)
}

func TestExecPrepared_Happy_Path(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	defer db.Close()

	pg := &database.Postgress{Pg: db}

	ctx := context.Background()
	query := "SELECT id, name FROM users WHERE id=$1"

	resultMock := sqlmock.NewResult(int64(1), int64(1))

	var expected []map[string]interface{}
	row := make(map[string]interface{})
	row["message"] = "success"
	row["rowsAffected"] = int64(1)
	expected = append(expected, row)

	params := []any{
		1,
	}

	// mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs(params[0]).WillReturnRows(rows)
	mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectExec().WithArgs(params[0]).WillReturnResult(resultMock)

	result, err := pg.ExecPrepared(ctx, query, params)
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.EqualValues(t, expected, result)
}

func TestGetSchema_Happy_Path(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	defer db.Close()

	pg := &database.Postgress{Pg: db}

	ctx := context.Background()
	query := `select column_name, data_type, character_maximum_length from INFORMATION_SCHEMA.COLUMNS where table_name =$1;`

	// expected rows to return
	rows := sqlmock.NewRows([]string{"column_name", "data_type", "character_maximum_length"}).
		AddRow("id", "integer", "null").
		AddRow("customername", "character varying", 200).
		AddRow("contactname", "character varying", 250)

	var expected []map[string]interface{}
	row := make(map[string]interface{})
	row["column_name"] = "id"
	row["data_type"] = "integer"
	row["character_maximum_length"] = "null"

	row2 := make(map[string]interface{})
	row2["column_name"] = "customername"
	row2["data_type"] = "character varying"
	row2["character_maximum_length"] = int64(200)

	row3 := make(map[string]interface{})
	row3["column_name"] = "contactname"
	row3["data_type"] = "character varying"
	row3["character_maximum_length"] = int64(250)
	expected = append(expected, row, row2, row3)

	var tables []string
	tbl := append(tables, "customers")

	mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs(tbl[0]).WillReturnRows(rows)

	result, err := pg.GetSchema(ctx, tbl)
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.EqualValues(t, expected, result)
}
