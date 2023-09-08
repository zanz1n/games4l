package user

import (
	"context"
	"time"

	"github.com/games4l/internal/auth"
	"github.com/games4l/internal/errors"
	"github.com/games4l/internal/utils"
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

func (s *UserService) SignInUser(ctx context.Context, credential string, passwd string) (string, errors.StatusCodeErr) {
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
		return "", errors.DefaultErrorList.UserUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwd))
	if err != nil {
		return "", errors.DefaultErrorList.UserUnauthorized
	}

	tokenPayload, err := s.ap.GenerateUserJwtToken(auth.JwtUserData{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Role:     user.Role,
	}, s.cfg.JwtExpiryTime)

	if err != nil {
		return "", errors.DefaultErrorList.InternalServerError
	}

	return tokenPayload, nil
}

func (s *UserService) FindByID(ctx context.Context, hexID string) (*User, errors.StatusCodeErr) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(hexID)

	if err != nil {
		return nil, errors.DefaultErrorList.InvalidObjectID
	}

	user := User{}

	err = s.col.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&user)

	if err != nil {
		return nil, errors.DefaultErrorList.EntityNotFound
	}

	return &user, nil
}

func (s *UserService) CreateUser(ctx context.Context, role auth.UserRole, data *CreateUserData) (*User, errors.StatusCodeErr) {
	if !utils.SliceContains(auth.ValidUserRoles, role) {
		return nil, errors.DefaultErrorList.InvalidRequestEntity
	}

	if !emailIsValid(data.Email) {
		return nil, errors.DefaultErrorList.InvalidEmail
	}

	if err := validate.Struct(*data); err != nil {
		return nil, errors.DefaultErrorList.InvalidRequestEntity
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	passwdEnc, err := bcrypt.GenerateFromPassword([]byte(data.Password), s.cfg.BcryptSaltLength)

	if err != nil {
		return nil, errors.DefaultErrorList.InternalServerError
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
		return nil, errors.DefaultErrorList.EntityAlreadyExists
	}

	return &user, nil
}