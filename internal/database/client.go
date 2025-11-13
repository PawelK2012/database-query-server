package database

import (
	"context"
)

// Repository implements commond DB client methods
// This allow each DB SDK to be wrapped in a Repository ie. Postgress, Redis etc
type ClientInterface interface {
	ExecQuery(ctx context.Context, query string, params map[string]any) ([]map[string]interface{}, error)
	GetSchema(ctx context.Context, tables []string) ([]map[string]interface{}, error)
}
