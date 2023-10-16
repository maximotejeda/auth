package auth

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const (
	MaxUserNameLen = 40
	MaxEmailLen    = 120
)

var (
	ErrInvalidUser     error = errors.New("invalid user string")
	ErrInvalidEmail    error = errors.New("invalid email string")
	ErrUserNotExist    error = errors.New("user not found")
	ErrUserExist       error = errors.New("user alredy on db")
	ErrInvalidPassword error = errors.New("password is not valid")
)

// User
// all params on user will be lower case to avoid normailzation
type User struct {
	ID       int
	Username string
	password string
	Email    string
	Rol      string
}

// ValidateUser
// Validate and modify user to have fields in lowercases
func ValidateUser(user *User) error {
	// check username length
	// check password length
	// check email format is correct
	// check Roll is a string list
	if len(user.Username) > MaxUserNameLen || len(user.Email) > MaxEmailLen {
		return fmt.Errorf("invalid length user or email")
	}
	if err := validateEmail(user.Email); err != nil {
		return err
	}
	if err := validateUsername(user.Username); err != nil {
		return err
	}
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	user.Rol = strings.ToLower(user.Rol)
	return nil
}

// validateEmail
// validates email strings with accepted params
// los unicos caracteres aceptados aqui son ?!%&._$- y \w and \d
// can be used to get distinct parts of an email
func validateEmail(email string) error {
	reCheck := regexp.MustCompile("(?P<name>[a-zA-Z0-9?!%&._$-]{1,20})@(?P<domain>[a-zA-Z0-9._-]{1,20})\\.(?P<sufix>[a-zA-Z0-9]{1,5})$")
	if !reCheck.MatchString(email) {
		return fmt.Errorf("email string %s: %w", email, ErrInvalidEmail)

	}
	// ?P<name> is a group submatch name we can retrieve them with
	defer func(results map[string]string) map[string]string {
		match := reCheck.FindStringSubmatch(email)
		for i, m := range reCheck.SubexpNames() {
			if i > 0 && i <= len(match) {
				results[m] = match[i]
			}
		}
		return results
	}(make(map[string]string))
	return nil
}

// validateUsername
// username must be a string with len between 6 and 40 chars
// with or without availabel special characters "_.?$@%!&-"
func validateUsername(username string) error {
	reCheck := regexp.MustCompile("[a-zA-Z0-9_.?$@%!&-]{6,40}")
	ok := reCheck.MatchString(username)
	if !ok {
		return fmt.Errorf("user string %s: %w", username, ErrInvalidUser)
	}
	return nil
}

// validatePassword
// Validate password to be as expected
// got the idea to check for special char without regexp
// https://stackoverflow.com/questions/55769838/ensure-specific-characters-are-in-a-string-regardless-of-position-using-a-rege
// if password contains # or | is invalid so return early
func validatePassword(pwd string) error {
	//check length
	if len(pwd) < 6 {
		return fmt.Errorf("password length %d :%w", len(pwd), ErrInvalidPassword)
	}
	var hasNumber, hasUpper, hasLower, hasSpecial bool
	for _, c := range pwd {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case c == '#' || c == '|':
			return fmt.Errorf("invalid characters %c: %w", c, ErrInvalidPassword)
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	if hasNumber && hasUpper && hasLower && hasSpecial {
		return nil
	}
	return fmt.Errorf("passwd: %w", ErrInvalidPassword)
}
