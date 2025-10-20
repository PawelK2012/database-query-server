package database

import "context"

// Repository implements commond DB client methods
// This allow each DB SDK to be wrapped in a Repository ie. Postgress, Redis etc
type ClientInterface interface {
	Init(ctx context.Context) error
}
