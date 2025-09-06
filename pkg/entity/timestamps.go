package entity

import "time"

type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTimeStamps(createdAt, updatedAt time.Time) Timestamps {
	return Timestamps{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (t *Timestamps) Touch() {
	t.UpdatedAt = time.Now()
}
