package data

import "github.com/sparrowsl/greenlight/internal/validator"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilters(val *validator.Validator, filters Filters) {
	// Check that the page and page_size parameters contain sensible values.
	val.Check(filters.Page > 0, "page", "must be greater than zero")
	val.Check(filters.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	val.Check(filters.PageSize > 0, "page_size", "must be greater than zero")
	val.Check(filters.PageSize <= 100, "page_size", "must be a maximum of 100")

	// Check that the sort parameter matches a value in the safelist.
	val.Check(validator.PermittedValue(filters.Sort, filters.SortSafelist...), "sort", "invalid sort value")
}
