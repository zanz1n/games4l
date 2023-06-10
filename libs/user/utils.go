package user

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

const MailRegex = `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`

func emailIsValid(email string) bool {
	match, err := regexp.MatchString(MailRegex, email)

	if err != nil {
		match = false
	}

	return match
}

func (u *User) ToEncodable() *UserJsonEncodable {
	return &UserJsonEncodable{
		ID:       u.ID.Hex(),
		Email:    u.Email,
		Username: u.Username,
		Role:     u.Role,
	}
}
