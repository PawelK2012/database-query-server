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

	assert.Equal(t, 1, 1)
}

// func Test_ExecPrepared_InsertUser(t *testing.T) {
// 	params := []any{
// 		1,
// 		"Joe",
// 		"Blogs",
// 		"joe@example.com",
// 		"password",
// 		true,
// 		time.Now(),
// 		time.Now(),
// 	}

// 	query := "INSERT INTO public.usersTest (id,firt_name,last_name, email, password, is_admin, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id"
// 	resp, err := testRepo.ExecPrepared(context.Background(), query, params)

// 	var expected []map[string]interface{}
// 	row := make(map[string]interface{})
// 	row["message"] = "success"
// 	row["rowsAffected"] = int64(1)
// 	expected = append(expected, row)

// 	if err != nil {
// 		t.Errorf("InsertUser test failed %s", err)
// 	}

// 	fmt.Printf("response %v \n", resp)

// 	assert.Equal(t, expected, resp)
// }

func TestPostgress_ExecPrepared(t *testing.T) {
	insertParams := []any{
		1,
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
		1,
	}

	insertStmt := "INSERT INTO public.usersTest (id,firt_name,last_name, email, password, is_admin, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id"
	updateStmt := "UPDATE public.usersTest SET firt_name = $1, email = $2, updated_at = $3 WHERE id = $4"
	updateStmtErr := "UPDATE public.usersTest SET firt_name = $1, email = $2, updated_at = $3, WHERE id = $4"

	var expected []map[string]interface{}
	row := make(map[string]interface{})
	row["message"] = "success"
	row["rowsAffected"] = int64(1)
	expected = append(expected, row)

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
