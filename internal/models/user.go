package models

import (
	"context"
	"errors"
	"time"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user (admin or customer)
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Name      string    `bun:"name,notnull"`
	Email     string    `bun:"email,unique,notnull"`
	Password  string    `bun:"password,notnull"` // hashed password
	Role      string    `bun:"role,notnull"`     // "admin" or "customer"
	CreatedAt time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,default:current_timestamp,nullzero"`
}

// BeforeInsert hook to set CreatedAt
func (u *User) BeforeInsert() {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}
}

// BeforeUpdate hook to set UpdatedAt
func (u *User) BeforeUpdate() {
	u.UpdatedAt = time.Now()
}

// GetUserByEmail fetches user by email from bun DB
func GetUserByEmail(ctx context.Context, db *bun.DB, email string) (*User, error) {
	var user User
	err := db.NewSelect().
		Model(&user).
		Where("email = ?", email).
		Scan(ctx)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GetUserByID fetches user by ID from bun DB
func GetUserByID(ctx context.Context, db *bun.DB, id int64) (*User, error) {
	var user User
	err := db.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// CheckPassword compares plain password with hashed password
func CheckPassword(plain, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}

// HashPassword helper
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}
