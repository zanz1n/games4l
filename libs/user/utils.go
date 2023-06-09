package user

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	nanoid "github.com/matoous/go-nanoid/v2"
)

const UserIDSize = 18

var validate = validator.New()

const mailRegex = `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`

func emailIsValid(email string) bool {
	match, err := regexp.MatchString(mailRegex, email)

	if err != nil {
		match = false
	}

	return match
}

func GenerateID() string {
	id, err := nanoid.New(UserIDSize)

	if err != nil {
		id, err = nanoid.New(UserIDSize)

		if err != nil {
			panic(err)
		}
	}

	return id
}

func (u *User) ToEncodable() *UserJsonEncodable {
	return &UserJsonEncodable{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
		Role:     u.Role,
	}
}
