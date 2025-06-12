package auth

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	RoleAdmin    string = "admin"
	RoleEmployee        = "employee"
)

type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Password  Password  `json:"-"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
}

type Password struct {
	hash []byte
}

func (p *Password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	p.hash = hash
	return nil
}

func (p *Password) Matches(plaintext string) bool {
	if len(p.hash) == 0 {
		return false
	}
	return bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext)) == nil
}

type Session struct {
	Token     string    `json:"token"`
	UserId    int       `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
