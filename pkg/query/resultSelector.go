package query

import (
	"github.com/leetatech/leeta_backend/pkg/query/filter"
	"github.com/leetatech/leeta_backend/pkg/query/paging"
)

// ResultSelector is a type that represents the selection criteria for querying data. It contains a filter, sorting, and paging information.
// Filter is a pointer to a filter.Request struct that specifies the filtering criteria for the query.
// Sorting is a pointer to a sorting.Request struct that specifies the sorting order for the query.
// Paging is a pointer to a paging.Request struct that specifies the paging configuration for the query.
type ResultSelector struct {
	Filter *filter.Request `json:"filter" binding:"omitempty"`
	Paging *paging.Request `json:"paging" binding:"omitempty"`
}
