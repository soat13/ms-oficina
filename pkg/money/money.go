package money

type Money struct {
	Cents int64
}

func New(cents int64) (Money, error) {
	if cents < 0 {
		return Money{}, ErrMoneyNegative
	}

	return Money{Cents: cents}, nil
}
