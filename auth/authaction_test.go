package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/maximotejeda/auth/sqles"
)

func TestRegister(t *testing.T) {
	cases := map[string]struct {
		user   User
		expect error
	}{
		"success/user":                     {User{Username: "juanito", Email: "juan@gmail.com", password: "test", Rol: "user"}, nil},
		"fail/user/username_too_short":     {User{Username: "juan", Email: "juan@email.com", password: "test", Rol: "user"}, ErrInvalidUser},
		"fail/user/email_bad_structure":    {User{Username: "juan", Email: "juanemail.com", password: "test", Rol: "user"}, ErrInvalidEmail},
		"fail/user/username_bad_character": {User{Username: "jua>", Email: "juan@email.com", password: "test", Rol: "user"}, ErrInvalidUser},
	}
	ctx := context.Background()
	db, _ := sqles.Connect(ctx, sqles.DefaultDriver, ":memory:")
	for name, scenario := range cases {
		t.Run(name, func(t *testing.T) {

			ac := AuthAction{DB: db}
			err := ac.Register(ctx, scenario.user)
			if err != nil && scenario.expect != nil {
				if !errors.Is(err, scenario.expect) {
					t.Errorf("\n\twant %v\n\tgot: %v", scenario.expect, err)
				}
			}

		})
	}

}

func TestDelete(t *testing.T) {
	cases := map[string]struct {
		user   User
		expect error
	}{
		"success/user/delete":      {User{Username: "juanito", Email: "juan@gmail.com"}, nil},
		"fail/NotExistUser/delete": {User{Username: "juanito", Email: "juan@gmail.com"}, ErrUserNotExist},
		"fail/BadUserName/delete":  {User{Username: "juan", Email: "juan@gmail.com"}, ErrInvalidUser},
		"fail/BadEmail/delete":     {User{Username: "juanito", Email: "juan_gmail.com"}, ErrInvalidEmail},
	}
	ctx := context.Background()
	db, _ := sqles.Connect(ctx, sqles.DefaultDriver, ":memory:")
	ac := AuthAction{DB: db}
	_ = ac.Register(ctx, User{Username: "juanito", Email: "juan@gmail.com", password: "test", Rol: "user"})
	for name, scenario := range cases {
		t.Run(name, func(t *testing.T) {
			err := ac.Delete(ctx, &scenario.user)
			if err != nil && scenario.expect != nil {
				if !errors.Is(err, scenario.expect) {
					t.Errorf("\n\twant %v\n\tgot: %v", scenario.expect, err)
				}
			} else if err != nil && scenario.expect == nil {
				t.Errorf("expected <nil> got: %s", err.Error())
			}

		})
	}

}

func TestLogin(t *testing.T) {
	cases := map[string]struct {
		user   User
		expect error
	}{
		"success/user/username/login":       {User{Username: "juanito", password: "test"}, nil},
		"success/user/email/login":          {User{Email: "juan@gmail.com", password: "test"}, nil},
		"fail/user/username/login":          {User{Username: "juanito1", password: "test"}, ErrUserNotExist},
		"fail/user/Email/login":             {User{Email: "juanito1@example.com", password: "test"}, ErrUserNotExist},
		"fail/user/BadEmailString/login":    {User{Email: "juanito1_example.com", password: "test"}, ErrInvalidEmail},
		"fail/user/BadUsernameString/login": {User{Username: "juan", password: "test"}, ErrInvalidUser},
	}
	ctx := context.Background()
	db, _ := sqles.Connect(ctx, sqles.DefaultDriver, ":memory:")
	ac := AuthAction{DB: db}
	_ = ac.Register(ctx, User{Username: "juanito", Email: "juan@gmail.com", password: "test", Rol: "user"})
	for name, scenario := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := ac.Login(ctx, &scenario.user, scenario.user.password)
			if err != nil && scenario.expect != nil {
				if !errors.Is(err, scenario.expect) {
					t.Errorf("\n\twant %v\n\tgot: %v", scenario.expect, err)
				}
			} else if err != nil && scenario.expect == nil {
				t.Errorf("expected <nil> got: %s", err.Error())
			}

		})
	}

}
