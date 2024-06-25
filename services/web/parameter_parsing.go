package web

import (
	"encoding/json"
	"fmt"
	"github.com/greenbone/opensight-golang-libraries/pkg/query"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/filter"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/paging"
	"github.com/greenbone/opensight-golang-libraries/pkg/query/sorting"
	"github.com/samber/lo"
	"net/http"
	"strings"
)

const DefaultLimit = 10

// PrepareResultSelector converts the common query parameters to a ResultSelector.
// Parameters are verified where possible and defaults are set.
func PrepareResultSelector(req *http.Request, filterOptions []filter.RequestOption, allowedSortFields []string, defaults query.ResultSelector) (resultSelector query.ResultSelector, err error) {
	resultSelector = query.ResultSelector{}
	err = json.NewDecoder(req.Body).Decode(&resultSelector)
	if err != nil {
		err = fmt.Errorf("failed to parse request body: %w", err)
		return resultSelector, err // TODO: declare error types
	}

	//apply defaults
	resultSelector = applyDefaults(resultSelector, defaults)

	err = validate(resultSelector, filterOptions, allowedSortFields)
	if err != nil {
		err = fmt.Errorf("failed to validate result: %w", err)
		return resultSelector, err
	}

	return resultSelector, nil
}

// ResultSelectorDefaults holds default result selectors
func ResultSelectorDefaults(sortingRequest *sorting.Request) query.ResultSelector {
	return query.ResultSelector{
		Paging: &paging.Request{
			PageIndex: 0,
			PageSize:  DefaultLimit,
		},
		Sorting: sortingRequest,
	}
}

func applyDefaults(resultSelector query.ResultSelector, defaults query.ResultSelector) query.ResultSelector {
	if resultSelector.Paging == nil {
		if resultSelector.Paging.PageIndex < 0 || resultSelector.Paging.PageSize < 1 {
			resultSelector.Paging = defaults.Paging
		}
	}
	if resultSelector.Sorting == nil {
		resultSelector.Sorting = defaults.Sorting
	} else {
		if resultSelector.Sorting.SortColumn == "" {
			resultSelector.Sorting.SortColumn = defaults.Sorting.SortColumn
			resultSelector.Sorting.SortDirection = defaults.Sorting.SortDirection
		}

		if resultSelector.Sorting.SortDirection == "" {
			resultSelector.Sorting.SortDirection = defaults.Sorting.SortDirection
		}
	}
	return resultSelector
}

func validate(resultSelector query.ResultSelector, filterOptions []filter.RequestOption, allowedSortFields []string) error {
	err := filter.ValidateFilter(resultSelector.Filter, filterOptions)
	if err != nil {
		return err
	}

	err = validateSorting(resultSelector.Sorting, allowedSortFields)
	if err != nil {
		return err
	}

	err = validatePaging(resultSelector.Paging)
	if err != nil {
		return err
	}

	return nil
}

func validateSorting(sortingRequest *sorting.Request, allowedSortFields []string) error {
	if sortingRequest == nil {
		return nil
	}

	err := sorting.ValidateSortingRequest(sortingRequest)
	if err != nil {
		return err
	}

	if !lo.Contains(allowedSortFields, sortingRequest.SortColumn) {
		return sorting.NewSortingError(fmt.Sprintf("%s is no valid sort column, possible values: %s",
			sortingRequest.SortColumn, strings.Join(allowedSortFields, ", ")))
	}

	return nil
}

func validatePaging(pagingRequest *paging.Request) error {
	if pagingRequest == nil {
		return nil
	}

	if pagingRequest.PageIndex < 0 {
		return paging.NewPagingError("%d is no valid page index, it must be >= 0", pagingRequest.PageIndex)
	}

	if pagingRequest.PageSize <= 0 {
		return paging.NewPagingError("%d is no valid page size, it must be > 0", pagingRequest.PageSize)
	}

	return nil
}
