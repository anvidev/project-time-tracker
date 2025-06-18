package users

import (
	"fmt"

	"github.com/anvidev/project-time-tracker/internal/types"
	"golang.org/x/crypto/bcrypt"
)

const (
	RoleAdmin    string = "admin"
	RoleEmployee        = "employee"
)

type User struct {
	Id        int64    `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Role      string   `json:"role"`
	Password  Password `json:"-"`
	IsActive  bool     `json:"isActive"`
	CreatedAt string   `json:"createdAt"` // yyyy-MM-dd (time.DateOnly)
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

type Hours struct {
	UserId    int64          `json:"userId"`
	Monday    types.Duration `json:"monday"`
	Tuesday   types.Duration `json:"tuesday"`
	Wednesday types.Duration `json:"wednesday"`
	Thursday  types.Duration `json:"thursday"`
	Friday    types.Duration `json:"friday"`
	Saturday  types.Duration `json:"saturday"`
	Sunday    types.Duration `json:"sunday"`
}

type RegisterUserInput struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
