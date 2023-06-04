package authdba

import (
	"context"
	"time"

	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	MongoDbName      string
	BcryptSaltLength int
	JwtExpiryTime    time.Duration
}

type UserService struct {
	client *mongo.Client
	col    *mongo.Collection
	cfg    *Config
	ap     *auth.AuthProvider
}

type User struct {
	ID       string        `bson:"_id,omitempty" validate:"required"`      // Primary key
	Username string        `bson:"username,omitempty" validate:"required"` // Index
	Email    string        `bson:"email,omitempty" validate:"required"`    // Index
	Password string        `bson:"password,omitempty" validate:"required"`
	Role     auth.UserRole `bson:"role,omitempty" validate:"required"`
}

type CreateUserData struct {
	Username string `json:"username,omitempty" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type UserJsonEncodable struct {
	ID       string        `json:"id"`
	Email    string        `json:"email"`
	Username string        `json:"username"`
	Role     auth.UserRole `json:"role"`
}

func NewUserService(client *mongo.Client, cfg *Config) *UserService {
	db := client.Database(cfg.MongoDbName)

	col := db.Collection("users")

	return &UserService{
		client: client,
		cfg:    cfg,
		col:    col,
	}
}

func (s *UserService) SignInUser(parentCtx context.Context, credential string, passwd string) (string, utils.StatusCodeErr) {
	ctx, cancel := context.WithTimeout(parentCtx, 15*time.Second)
	defer cancel()

	var (
		err  error
		user = User{}
	)

	if emailIsValid(credential) {
		err = s.col.FindOne(ctx, bson.D{{Key: "email", Value: credential}}).Decode(&user)
	} else {
		err = s.col.FindOne(ctx, bson.D{{Key: "username", Value: credential}}).Decode(&user)
	}

	if err != nil {
		return "", utils.NewStatusCodeErr(
			"user do not exist or password do not match",
			httpcodes.StatusUnauthorized,
		)
	}

	tokenPayload, err := s.ap.GenerateUserJwtToken(auth.JwtUserData{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}, s.cfg.JwtExpiryTime)

	if err != nil {
		return "", utils.NewStatusCodeErr(
			"something went wrong when trying to generate the auth token",
			httpcodes.StatusInternalServerError,
		)
	}

	return tokenPayload, nil
}

func (s *UserService) CreateUser(parentCtx context.Context, role auth.UserRole, data *CreateUserData) (*User, utils.StatusCodeErr) {
	if !utils.SliceContains(auth.ValidUserRoles, role) {
		return nil, utils.NewStatusCodeErr(
			"invalid user role enum type",
			httpcodes.StatusBadRequest,
		)
	}

	if !emailIsValid(data.Email) {
		return nil, utils.NewStatusCodeErr(
			"the provided email address is not valid",
			httpcodes.StatusBadRequest,
		)
	}

	if err := validate.Struct(*data); err != nil {
		return nil, utils.NewStatusCodeErr(
			"provided user creation data is not valid",
			httpcodes.StatusBadRequest,
		)
	}

	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	passwdEnc, err := bcrypt.GenerateFromPassword([]byte(data.Password), s.cfg.BcryptSaltLength)

	if err != nil {
		return nil, utils.NewStatusCodeErr(
			"failed to hash the password",
			httpcodes.StatusInternalServerError,
		)
	}

	user := User{
		ID:       GenerateID(),
		Username: data.Username,
		Email:    data.Email,
		Password: string(passwdEnc),
		Role:     role,
	}

	if _, err = s.col.InsertOne(ctx, user); err != nil {
		return nil, utils.NewStatusCodeErr(
			"failed to insert data into database",
			httpcodes.StatusInternalServerError,
		)
	}

	return &user, nil
}
