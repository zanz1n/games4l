package handler

import (
	"github.com/games4l/internal/auth"
)

var (
	applicationJsonHeader = map[string]string{
		"Content-Type": "application/json",
	}
	ap *auth.AuthProvider
)
