package domain

// --- Common Enums ---

// SortOrder represents the order direction for sorting
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// --- Common Structs ---

// PaginationOptions represents pagination configuration that can be used across all entities
type PaginationOptions struct {
	Limit  int    `json:"limit" example:"10"`
	Offset int    `json:"offset" example:"0"`
	Cursor string `json:"cursor,omitempty" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
}
