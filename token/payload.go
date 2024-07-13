package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken    = errors.New("token has expired")
	ErrInvalidIssuedAt = errors.New("issued token time is invalid")
	ErrInvalidID       = errors.New("token ID is invalid")
	ErrInvalidUsername = errors.New("username is invalid")
	ErrInvalidToken    = errors.New("token is invalid")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	// jwt.StandardClaims
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        id,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	//
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	if payload.IssuedAt.After(time.Now()) {
		return ErrInvalidIssuedAt
	}

	if payload.ID == uuid.Nil {
		return ErrInvalidID
	}

	if payload.Username == "" {
		return ErrInvalidUsername
	}

	return nil
}
