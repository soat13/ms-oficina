package pagination

func LimitAndOffset(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
