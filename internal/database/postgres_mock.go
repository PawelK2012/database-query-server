package database

import (
	"context"
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
	//TODO handle multiple rows
	db := c.mockSQLTable[0]
	var result []map[string]interface{}
	result = append(result, db)
	return result, nil
}
