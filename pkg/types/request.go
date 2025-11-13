package types

type SchemaRequest struct {
	Database string   `json:"database"`
	Tables   []string `json:"tables,omitempty"`
	Detailed bool     `json:"detailed,omitempty"`
}
type QueryRequest struct {
	Database   string         `json:"database"`
	Query      string         `json:"query"`
	Parameters map[string]any `json:"parameters,omitempty"`
	Format     string         `json:"format,omitempty"` // json, csv, table
	Limit      int            `json:"limit,omitempty"`
	Timeout    int            `json:"timeout,omitempty"`
}

type QueryResponse struct {
	Query    string `json:"query"`
	Response string `json:"response"`
	Format   string `json:"format,omitempty"` // json, csv, table
}
