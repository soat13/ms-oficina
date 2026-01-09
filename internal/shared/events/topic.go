package events

type (
	Event interface {
		Topic() string
	}
)
