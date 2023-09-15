package src

import (
	"github.com/games4l/internal/auth"
	"github.com/games4l/internal/user"
	"github.com/go-playground/validator/v10"
)

var (
	applicationJsonHeader = map[string]string{
		"Content-Type": "application/json",
	}
	validate = validator.New()
	ap       *auth.AuthProvider
	dba      *user.UserService
)

type JSON map[string]interface{}
