package qeetid

// BulkImportResult is the outcome of BulkCreate or BulkImport.
type BulkImportResult struct {
	Succeeded int                  `json:"succeeded"`
	Failed    int                  `json:"failed"`
	Errors    []BulkImportRowError `json:"errors,omitempty"`
}

// BulkImportRowError describes one failed row in a bulk import.
type BulkImportRowError struct {
	Line    int    `json:"line,omitempty"`
	Email   string `json:"email,omitempty"`
	Message string `json:"message"`
}
