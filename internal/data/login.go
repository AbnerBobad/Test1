package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/AbnerBobad/final_project/internal/validator"
	"github.com/lib/pq"
)

type User struct {
	UID            int64     `json:"user_id"`
	UName          string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"password_hash"`
	URole          string    `json:"role"`
	Active         bool      `json:"activated"`
	CreatedAt      time.Time `json:"created_at"`
}

// ErrDuplicateEmail is returned when a user tries to register with an email that already exists.
var ErrDuplicateEmail = errors.New("duplicate email")

type UserModel struct {
	DB *sql.DB
}

// validate user
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(validator.NotBlank(user.UName), "username", "Fullname must be provided")
	v.Check(validator.MaxLengths(user.UName, 50), "username", "Fullname must be less than 50 characters")
	v.Check(validator.MinLength(user.UName, 5), "username", "Fullname must be at least 5 characters")

	v.Check(validator.NotBlank(user.Email), "email", "Email must be provided")
	v.Check(validator.MaxLength(user.Email, 50), "email", "Email must be less than 50 Characters")
	v.Check(validator.IsValidEmail(user.Email), "email", "Email address is not valid")

	v.Check(validator.NotBlank(string(user.HashedPassword)), "password_hash", "Password must be provided")
	v.Check(validator.MaxLength(string(user.HashedPassword), 50), "password_hash", "Password must be less than 50 characters")
	v.Check(validator.MinLength(string(user.HashedPassword), 8), "password_hash", "Password must be at least 8 characters")

}

func (m *UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash, role, activated, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING user_id, created_at;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(
		ctx, query,
		user.UName,
		user.Email,
		user.HashedPassword,
		user.URole,
		user.Active,
	).Scan(&user.UID, &user.CreatedAt)
	if err != nil {
		// Detect duplicate email via pq error code 23505
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}
