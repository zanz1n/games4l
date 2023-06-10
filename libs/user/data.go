package user

import (
	"context"
	"time"

	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID       primitive.ObjectID `bson:"_id,omitempty" validate:"required"`      // Primary key
	Username string             `bson:"username,omitempty" validate:"required"` // Index
	Email    string             `bson:"email,omitempty" validate:"required"`    // Index
	Password string             `bson:"password,omitempty" validate:"required"`
	Role     auth.UserRole      `bson:"role,omitempty" validate:"required"`
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

func NewUserService(client *mongo.Client, ap *auth.AuthProvider, cfg *Config) *UserService {
	db := client.Database(cfg.MongoDbName)

	col := db.Collection("users")

	return &UserService{
		client: client,
		cfg:    cfg,
		col:    col,
		ap:     ap,
	}
}

func (s *UserService) SignInUser(ctx context.Context, credential string, passwd string) (string, utils.StatusCodeErr) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwd))
	if err != nil {
		return "", utils.NewStatusCodeErr(
			"user do not exist or password do not match",
			httpcodes.StatusUnauthorized,
		)
	}

	tokenPayload, err := s.ap.GenerateUserJwtToken(auth.JwtUserData{
		ID:       user.ID.Hex(),
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

func (s *UserService) FindByID(ctx context.Context, hexID string) (*User, utils.StatusCodeErr) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(hexID)

	if err != nil {
		return nil, utils.NewStatusCodeErr(
			"invalid object id format",
			httpcodes.StatusBadRequest,
		)
	}

	user := User{}

	err = s.col.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&user)

	if err != nil {
		return nil, utils.NewStatusCodeErr(
			"user could not be found",
			httpcodes.StatusNotFound,
		)
	}

	return &user, nil
}

func (s *UserService) CreateUser(ctx context.Context, role auth.UserRole, data *CreateUserData) (*User, utils.StatusCodeErr) {
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

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	passwdEnc, err := bcrypt.GenerateFromPassword([]byte(data.Password), s.cfg.BcryptSaltLength)

	if err != nil {
		return nil, utils.NewStatusCodeErr(
			"failed to hash the password",
			httpcodes.StatusInternalServerError,
		)
	}

	oid := primitive.NewObjectID()

	user := User{
		ID:       oid,
		Username: data.Username,
		Email:    data.Email,
		Password: string(passwdEnc),
		Role:     role,
	}

	if _, err = s.col.InsertOne(ctx, user); err != nil {
		return nil, utils.NewStatusCodeErr(
			"the user already exists",
			httpcodes.StatusConflict,
		)
	}

	return &user, nil
}
