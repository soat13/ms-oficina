package domain

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	Value     string
	ExpiresAt time.Time
}

type Claims struct {
	UserID uuid.UUID
	Email  string
	Roles  []string
	Exp    time.Time
}
