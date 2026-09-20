package jwtservice

import "github.com/google/uuid"

// TokenUser
type TokenUser struct {
	ID   uuid.UUID
	Name string
	Role []string
}

func NewTokenUser(id uuid.UUID, name string, role []string) *TokenUser {
	return &TokenUser{
		ID:   id,
		Name: name,
		Role: role,
	}
}
