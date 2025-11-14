package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	DB_USER     = os.Getenv("POSTGRES_USER")
	DB_PASSWORD = os.Getenv("POSTGRES_PW")
	DB_NAME     = os.Getenv("POSTGRES_DB")
)

type Postgress struct {
	Pg *sql.DB
}

func NewPostgressClient() (ClientInterface, error) {
	// remove host=postgresql when running localy
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", DB_USER, DB_NAME, DB_PASSWORD)

	pg, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Postgress{
		Pg: pg,
	}, err
}

// ExecQuery execute SQL queries with prepared statements
func (s *Postgress) ExecQuery(ctx context.Context, query string, params map[string]any) ([]map[string]interface{}, error) {
	// remove
	log.Printf("ExecQuery query: %v \n", query)
	log.Printf("ExecQuery params: %v \n", params)

	var allMaps []map[string]interface{}
	var rows *sql.Rows
	stmt, err := s.Pg.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// decide if system should use QueryContext with params
	if len(params) > 0 {
		var args []any

		for k := range params {
			fmt.Printf("key[%s] value[%s]\n", k, params[k])
			args = append(args, params[k])
		}
		rows, err = stmt.QueryContext(ctx, args...)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = stmt.QueryContext(ctx)
		if err != nil {
			return nil, err
		}
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if len(columns) == 0 {
		// Default response for operations that don't return rows
		// ie. CREATE SCHEMA schema_name OR INSERT INTO schema_name.table etc
		defaultResp := make(map[string]interface{})
		defaultResp["message"] = "success"
		allMaps = append(allMaps, defaultResp)

	} else {
		// The system handles dynamic queries, so the results are scanned into a slice of pointers to interface{} variables.
		for rows.Next() {
			values := make([]interface{}, len(columns))
			pointers := make([]interface{}, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			if err := rows.Scan(pointers...); err != nil {
				// Check for a scan error.
				// Query rows will be closed with defer.
				return nil, err
			}
			resultMap := make(map[string]interface{})
			for i, val := range values {
				resultMap[columns[i]] = val
			}
			allMaps = append(allMaps, resultMap)
		}
		// If the database is being written to ensure to check for Close
		// errors that may be returned from the driver. The query may
		// encounter an auto-commit error and be forced to rollback changes.
		err = rows.Close()
		if err != nil {
			return nil, err
		}

		// Rows.Err will report the last error encountered by Rows.Scan.
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	// remove
	log.Printf("ExecQuery response: %v \n ", allMaps)
	return allMaps, nil
}

// ExecQuery execute SQL queries with prepared statements
func (s *Postgress) GetSchema(ctx context.Context, tables []string) ([]map[string]interface{}, error) {
	log.Printf("GetSchema params: %v \n", tables)

	q := `select column_name, data_type, character_maximum_length from INFORMATION_SCHEMA.COLUMNS where table_name =$1;`

	var allMaps []map[string]interface{}
	var rows *sql.Rows
	stmt, err := s.Pg.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var args []any

	for k := range tables {
		args = append(args, tables[k])
	}
	rows, err = stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			return nil, err
		}
		resultMap := make(map[string]interface{})
		for i, val := range values {
			resultMap[columns[i]] = val
		}
		allMaps = append(allMaps, resultMap)
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	err = rows.Close()
	if err != nil {
		return nil, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allMaps, nil
}
