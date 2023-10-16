// actions from auth are made here
package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/maximotejeda/auth/sqles"
	"golang.org/x/crypto/bcrypt"
)

type AuthAction struct {
	// TODO db conction from sqles
	DB *sqles.DB
}

// Login
// user must be populated either with username or email
// as we can omit or username or email we dont validate the complete user just what we need
// compare credentials provided with a db version
func (a *AuthAction) Login(ctx context.Context, user *User, password string) (*User, error) {
	if user.Username == "" && user.Email == "" {
		return nil, fmt.Errorf("login information required")
	}
	if user.Username != "" {
		err := validateUsername(user.Username)
		if err != nil {
			return nil, err
		}
	}
	if user.Email != "" {
		err := validateEmail(user.Email)
		if err != nil {
			return nil, err
		}
	}
	const query = `
              SELECT username, password, email, rol
                 FROM users
                      WHERE username = ? OR email = ?;`
	result := User{}
	err := a.DB.QueryRowContext(ctx, query, user.Username, user.Email).Scan(&result.Username, &result.password, &result.Email, &result.Rol)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotExist
		}
		return nil, fmt.Errorf("login querying database: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(result.password), []byte(password)); err != nil {
		return nil, fmt.Errorf("comparing passwords: %w", err)
	}
	result.password = "" //  as returning pointer no pass provided
	// TODO pending populate user info
	return &result, nil
}

// Register
// Add new row to database with user information
// Validation over username or email is made to sanitize input
func (a *AuthAction) Register(ctx context.Context, user User) error {
	if err := ValidateUser(&user); err != nil {
		return err
	}
	const query = `
              INSERT INTO users (
                  username, password, email, rol, created_at, edited_at
              ) VALUES (
              ?, ?, ?, ?, datetime(), datetime() 
              );`
	password, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}
	_, err = a.DB.ExecContext(ctx, query, user.Username, string(password), user.Email, user.Rol)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}
	// TODO pending populate user info
	return nil
}

func (a *AuthAction) Delete(ctx context.Context, user *User) error {
	if err := ValidateUser(user); err != nil {
		return err
	}
	const query = `
              DELETE FROM users
                  WHERE username = ? AND email = ?`
	res, err := a.DB.ExecContext(ctx, query, user.Username, user.Email)
	if err != nil {
		return fmt.Errorf("deleting record: %w", err)
	}
	if nm, _ := res.RowsAffected(); nm == 0 {
		return fmt.Errorf("no record to delete: %w", ErrUserNotExist)
	}
	return nil

}

// Validate
// user information and that  exists on DB
func (a *AuthAction) Validate(ctx context.Context, user User) {}

// Refresh
// user token if a valid time has passed
// the time can not be superior to the issued time in the token
// example if the token is valid for an hour and expired and time elapsed
// is less than an hour the refresh can be made
func (a *AuthAction) Refresh(ctx context.Context, user User) {}
