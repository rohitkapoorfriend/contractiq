package contract

// Filter defines the criteria for querying contracts.
type Filter struct {
	Status   *Status
	PartyID  *string
	OwnerID  *string
	Search   *string
	Page     int
	PageSize int
}

// DefaultFilter returns a filter with sensible defaults.
func DefaultFilter() Filter {
	return Filter{
		Page:     1,
		PageSize: 20,
	}
}

// Offset returns the SQL offset for pagination.
func (f Filter) Offset() int {
	return (f.Page - 1) * f.PageSize
}

// PageResult wraps a page of results with total count for pagination.
type PageResult struct {
	Items      []*Contract
	TotalCount int
	Page       int
	PageSize   int
}

// TotalPages returns the total number of pages.
func (p PageResult) TotalPages() int {
	if p.PageSize == 0 {
		return 0
	}
	pages := p.TotalCount / p.PageSize
	if p.TotalCount%p.PageSize > 0 {
		pages++
	}
	return pages
}
