package password

import (
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	Hash string
}

func New(plainPassword string) (Password, error) {
	if len(plainPassword) < 8 {
		return Password{}, ErrPasswordTooShort
	}
	if len(plainPassword) > 72 {
		return Password{}, ErrPasswordTooLong
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{Hash: string(hash)}, nil
}

func FromHash(hash string) (Password, error) {
	if !isValidHash(hash) {
		return Password{}, ErrInvalidPasswordHash
	}
	return Password{Hash: hash}, nil
}

func (p Password) Matches(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Hash), []byte(plainPassword))
	return err == nil
}

func isValidHash(hash string) bool {
	_, err := bcrypt.Cost([]byte(hash))
	return err == nil
}
