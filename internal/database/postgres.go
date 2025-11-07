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
	//log.Println("connStr:", connStr)

	pg, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Postgress{
		Pg: pg,
	}, err
}

func (s *Postgress) ExecQuery(ctx context.Context, query string, params map[string]any) ([]map[string]interface{}, error) {
	// remove
	log.Println("----")
	log.Printf("ExecQuery %v", query)
	log.Println("----")
	log.Printf("ExecQuery params %v", params)
	log.Println("----")

	var rows *sql.Rows
	stmt, err := s.Pg.PrepareContext(ctx, query)
	if err != nil {
		// improve
		log.Fatal(err)
	}
	defer stmt.Close()

	if len(params) > 0 {
		var args []any
		for k := range params {
			fmt.Printf("key[%s] value[%s]\n", k, params[k])
			args = append(args, params[k])
		}
		rows, err = stmt.QueryContext(ctx, args...)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		rows, err = stmt.QueryContext(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var allMaps []map[string]interface{}
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
			log.Fatal(err)
		}
		resultMap := make(map[string]interface{})
		for i, val := range values {
			fmt.Printf("Adding key=%s val=%v\n", columns[i], val)
			resultMap[columns[i]] = val
		}
		allMaps = append(allMaps, resultMap)
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// remove
	log.Println("-----db q res")
	log.Printf("names %v ", allMaps)
	log.Println("-----db q res")
	return allMaps, nil
}

// The logic in Init() could be improved, but it's not a priority right now.
// I might revisit and refine this function later.
func (s *Postgress) Init(ctx context.Context) error {
	fmt.Printf("initialising Customer table with test values")
	query := `CREATE TABLE IF NOT EXISTS Customers (
		id SERIAL PRIMARY KEY,
		CustomerName VARCHAR(200),
		ContactName VARCHAR(250),
		Address VARCHAR(500),
		City VARCHAR(250),
		PostalCode VARCHAR(150),
		Country VARCHAR(250),
		created_at TIMESTAMP
	)`

	_, err := s.Pg.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	queryInsert := `INSERT INTO Customers (CustomerName, ContactName, Address, City, PostalCode, Country)
			VALUES
			('Cardinal', 'Tom B. Erichsen', 'Skagen 21', 'Stavanger', '4006', 'Norway'),
			('Greasy Burger', 'Per Olsen', 'Gateveien 15', 'Sandnes', '4306', 'Norway'),
			('Tasty Tee', 'Finn Egan', 'Streetroad 19B', 'Liverpool', 'L1 0AA', 'UK');`

	_, err = s.Pg.ExecContext(ctx, queryInsert)
	return err
}
