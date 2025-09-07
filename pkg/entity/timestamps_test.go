package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTimestamps(t *testing.T) {
	now := time.Now()
	timestamps := NewTimestamps(now, now)

	assert.Equal(t, now, timestamps.CreatedAt)
	assert.Equal(t, now, timestamps.UpdatedAt)
}

func TestTimestamps_Touch(t *testing.T) {
	now := time.Now()
	timestamps := NewTimestamps(now, now)

	time.Sleep(1 * time.Second)

	timestamps.Touch()

	assert.True(t, timestamps.UpdatedAt.After(now))
}
