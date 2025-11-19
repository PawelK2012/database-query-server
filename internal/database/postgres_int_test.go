package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB
var testRepo ClientInterface

// var testRepo repository.Repository

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	//clean up
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	// populate DB with empty tables
	err = createTables()
	if err != nil {
		log.Fatalf("error creating DB tables %s", err)
	}
	testRepo = &Postgress{Pg: db}
	// run tests
	m.Run()

}

func createTables() error {
	sqlTbl, err := os.ReadFile("./testData/users.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(sqlTbl))
	if err != nil {
		return err
	}
	return nil
}

func Test_pingDB(t *testing.T) {
	err := db.Ping()
	if err != nil {
		t.Error("can't ping DB")
	}
}

func TestPostgress_ExecPrepared(t *testing.T) {
	insertParams := []any{
		3,
		"Joe",
		"Blogs",
		"joe@example.com",
		"password",
		true,
		time.Now(),
		time.Now(),
	}

	updateParams := []any{
		"Joe - UPDATED",
		"joe@example.com - UPDATED",
		time.Now(),
		3,
	}

	insertStmt := "INSERT INTO public.usersTest (id,firt_name,last_name, email, password, is_admin, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id"
	insertMult := "INSERT INTO public.usersTest (id,firt_name,last_name, email, password, is_admin, created_at, updated_at) VALUES (4, 'MrPawel', 'My surname', 'some@email.com', 'pass1', false, '2004-10-19 10:23:54', '2004-10-19 10:23:54'), (5, 'Mr. X', 'xXx', 'x@email.com', 'passx', false, '2004-10-19 10:23:54', '2004-10-19 10:23:54');"
	updateStmt := "UPDATE public.usersTest SET firt_name = $1, email = $2, updated_at = $3 WHERE id = $4"
	updateStmtErr := "UPDATE public.usersTest SET firt_name = $1, email = $2, updated_at = $3, WHERE id = $4"

	var expected []map[string]interface{}
	row := make(map[string]interface{})
	row["message"] = "success"
	row["rowsAffected"] = int64(1)
	expected = append(expected, row)

	var expectedMultInsert []map[string]interface{}
	row1 := make(map[string]interface{})
	row1["message"] = "success"
	row1["rowsAffected"] = int64(2)
	expectedMultInsert = append(expectedMultInsert, row1)

	multipleRowsParams := []any{}

	tests := []struct {
		name      string
		statement string
		params    []any
		want      []map[string]interface{}
		wantErr   bool
	}{
		// TODO add more tests
		{name: "Happy Flow execute_prepared - INSERT INTO public.usersTest", statement: insertStmt, params: insertParams, want: expected, wantErr: false},
		{name: "Happy Flow execute_prepared - UPDATE public.usersTest", statement: updateStmt, params: updateParams, want: expected, wantErr: false},
		{name: "Happy Flow execute_prepared - INSERT multiple rows", statement: insertMult, params: multipleRowsParams, want: expectedMultInsert, wantErr: false},
		{name: "Sad Flow execute_prepared - UPDATE public.usersTest failed", statement: updateStmtErr, params: updateParams, want: expected, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := testRepo.ExecPrepared(context.Background(), tt.statement, tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExecPrepared() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ExecPrepared() succeeded unexpectedly")
			}

			if true {
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestPostgress_ExecQuery(t *testing.T) {
	paramsEmpty := map[string]any{}
	params := map[string]any{}
	params["id"] = 1
	query := "SELECT id, firt_name, last_name FROM public.usersTest;"
	querySelectById := "SELECT id, firt_name, last_name FROM public.usersTest WHERE id = $1"
	queryErr := "SELECT id, firt_name, last_name FROM public.UNDEFINED WHERE id = $1"

	var expectedSelectById []map[string]interface{}
	row := make(map[string]interface{})
	row["firt_name"] = "TestUser1"
	row["id"] = int64(1)
	row["last_name"] = "Test surname"
	expectedSelectById = append(expectedSelectById, row)

	var expectedSelectAll []map[string]interface{}
	row1 := make(map[string]interface{})
	row1["firt_name"] = "TestUser1"
	row1["id"] = int64(1)
	row1["last_name"] = "Test surname"

	row2 := make(map[string]interface{})
	row2["firt_name"] = "TestUser2"
	row2["id"] = int64(2)
	row2["last_name"] = "Surname 2"

	row3 := make(map[string]interface{})
	row3["firt_name"] = "Joe - UPDATED"
	row3["id"] = int64(3)
	row3["last_name"] = "Blogs"

	row4 := make(map[string]interface{})
	row4["firt_name"] = "MrPawel"
	row4["id"] = int64(4)
	row4["last_name"] = "My surname"

	row5 := make(map[string]interface{})
	row5["firt_name"] = "Mr. X"
	row5["id"] = int64(5)
	row5["last_name"] = "xXx"
	expectedSelectAll = append(expectedSelectAll, row1, row2, row3, row4, row5)

	tests := []struct {
		name    string
		query   string
		params  map[string]any
		want    []map[string]interface{}
		wantErr bool
	}{
		// TODO keep an eye on SELECT * FROM test, as it could cause issues
		// This test also depends on the two users inserted in the ExecPrepared test case.
		{name: "Happy Flow execute_query - SELECT * FROM public.usersTest", query: query, params: paramsEmpty, want: expectedSelectAll, wantErr: false},
		{name: "Happy Flow execute_query - SELECT by id", query: querySelectById, params: params, want: expectedSelectById, wantErr: false},
		{name: "SAD Flow execute_query - SELECT by id", query: queryErr, params: params, want: expectedSelectById, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := testRepo.ExecQuery(context.Background(), tt.query, tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExecQuery() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ExecQuery() succeeded unexpectedly")
			}
			if true {
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestPostgress_GetSchema(t *testing.T) {
	tables := []string{
		"userstest",
	}

	tablesDoesntExist := []string{
		"userstestXXXX",
	}

	var expectedEmpty []map[string]interface{}

	var expected []map[string]interface{}
	row := make(map[string]interface{})
	row["column_name"] = "id"
	row["data_type"] = "integer"
	row["character_maximum_length"] = interface{}(nil)

	row1 := make(map[string]interface{})
	row1["column_name"] = "firt_name"
	row1["data_type"] = "character varying"
	row1["character_maximum_length"] = int64(255)

	row2 := make(map[string]interface{})
	row2["column_name"] = "last_name"
	row2["data_type"] = "character varying"
	row2["character_maximum_length"] = int64(255)

	row3 := make(map[string]interface{})
	row3["column_name"] = "email"
	row3["data_type"] = "character varying"
	row3["character_maximum_length"] = int64(255)

	row4 := make(map[string]interface{})
	row4["column_name"] = "password"
	row4["data_type"] = "character varying"
	row4["character_maximum_length"] = int64(60)

	row5 := make(map[string]interface{})
	row5["column_name"] = "is_admin"
	row5["data_type"] = "boolean"
	row5["character_maximum_length"] = interface{}(nil)

	row6 := make(map[string]interface{})
	row6["column_name"] = "created_at"
	row6["data_type"] = "timestamp without time zone"
	row6["character_maximum_length"] = interface{}(nil)

	row7 := make(map[string]interface{})
	row7["column_name"] = "updated_at"
	row7["data_type"] = "timestamp without time zone"
	row7["character_maximum_length"] = interface{}(nil)
	expected = append(expected, row, row1, row2, row3, row4, row5, row6, row7)

	tests := []struct {
		name    string
		tables  []string
		want    []map[string]interface{}
		wantErr bool
	}{
		{name: "Happy Flow get_schema - GET userstest SCHEMA", tables: tables, want: expected, wantErr: false},
		{name: "Happy Flow get_schema - GET schema for table that doesn't exist", tables: tablesDoesntExist, want: expectedEmpty, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, gotErr := testRepo.GetSchema(context.Background(), tt.tables)
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
