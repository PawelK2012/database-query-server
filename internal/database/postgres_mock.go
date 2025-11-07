package database

import (
	"context"
	"fmt"
)

type PostgresClientMock struct {
	mockSQLTable    []map[string]interface{}
	simulateFailure bool
}

func NewPostgresClientMock(mockSQLTable []map[string]interface{}, simulateFailure bool) (ClientInterface, error) {
	return &PostgresClientMock{
		mockSQLTable:    mockSQLTable,
		simulateFailure: simulateFailure,
	}, nil
}

func (c *PostgresClientMock) ExecQuery(ctx context.Context, query string, params map[string]any) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	if c.simulateFailure {
		return nil, fmt.Errorf("error executing postgres query")
	}
	return result, nil
}
