package handlers

import "exmple.com/database-query-server/internal/repository"

type QueryHandler struct {
	repository *repository.Repository
}

// NewSystemCollector creates a new system metrics collector
func NewQueryHandler(repository *repository.Repository) *QueryHandler {
	return &QueryHandler{
		repository: repository,
	}
}
