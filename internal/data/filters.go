package data

import (
	"math"
	"strings"

	"github.com/sparrowsl/greenlight/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
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

func calculateMetadata(totalRecords int, page int, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func (filters Filters) limit() int {
	return filters.PageSize
}

func (filters Filters) offset() int {
	return (filters.Page - 1) * filters.PageSize
}

func (filters Filters) sortColumn() string {
	for _, safeValue := range filters.SortSafelist {
		if filters.Sort == safeValue {
			return strings.TrimPrefix(filters.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + filters.Sort)
}

func (filters Filters) sortDirection() string {
	if strings.HasPrefix(filters.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}
