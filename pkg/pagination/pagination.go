package pagination

import "strconv"

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

func AtoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}
