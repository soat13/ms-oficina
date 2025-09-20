package domain

import (
    "errors"
    "strings"
    "time"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

var (
    ErrInvalidEmail    = errors.New("invalid user email")
    ErrInvalidName     = errors.New("invalid user name")
    ErrInvalidPassword = errors.New("invalid user password")
)

type User struct {
    ID           uuid.UUID
    Name         string
    Email        string
    PasswordHash string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func NewUser(id uuid.UUID, name, email, rawPassword string, now time.Time) (User, error) {
    name = strings.TrimSpace(name)
    email = strings.ToLower(strings.TrimSpace(email))
    if len(name) < 3 {
        return User{}, ErrInvalidName
    }
    if !strings.Contains(email, "@") || len(email) < 6 {
        return User{}, ErrInvalidEmail
    }
    if len(rawPassword) < 6 {
        return User{}, ErrInvalidPassword
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
    if err != nil {
        return User{}, err
    }
    return User{
        ID:           id,
        Name:         name,
        Email:        email,
        PasswordHash: string(hash),
        CreatedAt:    now,
        UpdatedAt:    now,
    }, nil
}

func (u User) VerifyPassword(raw string) bool {
    return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(raw)) == nil
}

