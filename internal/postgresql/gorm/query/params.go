package query

// * SortOptions mendefinisikan parameter untuk sorting.
type SortOptions struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

// * PaginationOptions mendefinisikan parameter untuk pagination.
type PaginationOptions struct {
	Limit  int
	Offset int
	Cursor string
}

// * Params adalah container untuk semua parameter query.
// * Filters menggunakan `any` agar setiap repository bisa mendefinisikan struct filternya sendiri.
type Params struct {
	SearchQuery *string
	Filters     any
	Sort        *SortOptions
	Pagination  *PaginationOptions
}
