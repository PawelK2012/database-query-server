package database_test

import (
	"context"
	"log"
	"regexp"
	"testing"

	"exmple.com/database-query-server/internal/database"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestExecQuery_Happy_Path(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	defer db.Close()

	pg := &database.Postgress{Pg: db}

	ctx := context.Background()
	query := "SELECT id, name FROM users WHERE id = $1"

	// expected rows to return
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Alice").
		AddRow(2, "Bob")

	// query params
	par := make(map[string]any)
	par["1"] = "Alice"

	mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, name FROM users WHERE id = $1")).ExpectQuery().WithArgs(par["1"]).WillReturnRows(rows)

	result, err := pg.ExecQuery(ctx, query, par)
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	log.Printf("resultsssss %v", result)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestExecQuery_Success_Selecting_2_Row(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	defer db.Close()

	pg := &database.Postgress{Pg: db}

	ctx := context.Background()
	query := "SELECT id, name FROM users WHERE id = $1"

	// expected rows to return
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Bob")

	// query params
	par := make(map[string]any)
	par["1"] = "Alice"

	mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs(par["1"]).WillReturnRows(rows)

	result, err := pg.ExecQuery(ctx, query, par)
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	log.Printf("resultsssss %v", result)

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
	query := "select * from users;"

	// expected rows to return
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Alice").
		AddRow(2, "Bob")

	// query params
	par := make(map[string]any)
	par["1"] = "Alice"
	par["2"] = "Bob"

	mock.ExpectPrepare(regexp.QuoteMeta("select * from users;")).ExpectQuery().WithArgs().WillReturnRows(rows)

	result, err := pg.ExecQuery(ctx, query, par)
	if err != nil {
		t.Errorf("ExecQuery() failed: %v", err)
		return
	}
	log.Printf("resultsssss %v", result)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
