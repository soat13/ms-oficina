package pagination

type Pagination struct {
	Limit  int
	Offset int
}

func New(limit, offset int) *Pagination {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return &Pagination{Limit: limit, Offset: offset}
}
