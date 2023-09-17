package src

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/games4l/internal/auth"
	"github.com/games4l/internal/sqli"
	"github.com/go-playground/validator/v10"
)

var (
	applicationJsonHeader = map[string]string{
		"Content-Type": "application/json",
	}
	validate = validator.New()
	ap       *auth.AuthProvider
	dba      sqli.Querier

	s3c      *s3.Client
	s3Bucket string
)

type JSON map[string]interface{}
