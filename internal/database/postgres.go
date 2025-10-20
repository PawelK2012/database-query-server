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
	DB_USER     = os.Getenv("POSTGRES_MCP_DEMO_USER")
	DB_PASSWORD = os.Getenv("POSTGRES_PASS_MCP_DEMO")
	DB_NAME     = os.Getenv("POSTGRES_MCP_DB_NAME")
)

type Postgress struct {
	pg *sql.DB
}

func NewPostgressClient() (ClientInterface, error) {
	// remove host=postgresql when running localy
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", DB_USER, DB_NAME, DB_PASSWORD)
	log.Println("connStr:", connStr)

	pg, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Postgress{
		pg: pg,
	}, err
}

func (s *Postgress) Init(ctx context.Context) error {
	fmt.Printf("Calling Postgres init func")
	query := `CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		author VARCHAR(200),
		title VARCHAR(250),
		description VARCHAR(5000),
		tags VARCHAR(250),
		created_at TIMESTAMP
	)`
	_, err := s.pg.ExecContext(ctx, query)
	return err
}
