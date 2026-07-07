package pagination

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type Params struct {
	Page     int
	PageSize int
}

func New(
	page int,
	pageSize int,
) Params {
	if page < 1 {
		page = DefaultPage
	}

	if pageSize < 1 {
		pageSize = DefaultPageSize
	}

	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return Params{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p Params) Limit() int {
	return p.PageSize
}

type Result[T any] struct {
	Items      []T
	Page       int
	PageSize   int
	TotalItems int64
	TotalPages int
}

func NewResult[T any](
	items []T,
	totalItems int64,
	params Params,
) Result[T] {
	if items == nil {
		items = make([]T, 0)
	}

	totalPages := 0

	if totalItems > 0 {
		totalPages = int(
			(totalItems + int64(params.PageSize) - 1) /
				int64(params.PageSize),
		)
	}

	return Result[T]{
		Items:      items,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
