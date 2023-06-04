package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
	"github.com/golang-jwt/jwt/v4"
)

type ByteEncoding string

const (
	ByteEncodingBase64 ByteEncoding = "BASE64"
	ByteEncodingHex    ByteEncoding = "HEX"
)

type UserRole string

const (
	UserRolePacient UserRole = "PACIENT"
	UserRoleAdmin   UserRole = "ADMIN"
	UserRoleClient  UserRole = "CLIENT"
)

var ValidUserRoles = []UserRole{UserRoleAdmin, UserRoleClient, UserRolePacient}

type JwtUserData struct {
	Role     UserRole `json:"role"`
	Username string   `json:"username"`
	ID       string   `json:"id"`
}

type AuthProvider struct {
	sigKey []byte
	jwtKey []byte
}

func NewAuthProvider(sigKey []byte, jwtKey []byte) *AuthProvider {
	return &AuthProvider{
		sigKey: sigKey,
		jwtKey: jwtKey,
	}
}

func (ap *AuthProvider) AuthUser(payload string) (*JwtUserData, utils.StatusCodeErr) {
	var (
		valId       string
		valUsername string
		valRole     UserRole
	)
	token, err := jwt.Parse(payload, func(t *jwt.Token) (interface{}, error) {
		var (
			ok     bool
			claims jwt.MapClaims
		)

		formatErr := errors.New("invalid auth token metadata format, please sign-out ang log-in again")

		if _, ok = t.Method.(*jwt.SigningMethodHMAC); !ok {
			headerAlg, ok := t.Header["alg"].(string)
			if !ok {
				headerAlg = "<IMPOSSIBLE TO DECODE>"
			}

			return nil, fmt.Errorf("invalid auth token signing method: %s", headerAlg)
		}

		if claims, ok = t.Claims.(jwt.MapClaims); !ok {
			return nil, formatErr
		}

		if valId, ok = claims["id"].(string); !ok {
			return nil, formatErr
		}

		if valUsername, ok = claims["username"].(string); !ok {
			return nil, formatErr
		}

		if valRole, ok = claims["role"].(UserRole); !ok {
			return nil, formatErr
		}

		if !utils.SliceContains(ValidUserRoles, valRole) {
			return nil, formatErr
		}

		if t.Claims.(jwt.MapClaims)["exp"].(float64) < float64(time.Now().Unix()) {
			return nil, errors.New("token is expired")
		}

		return []byte(ap.jwtKey), nil
	})

	if err != nil || !token.Valid {
		return nil, utils.NewStatusCodeErr(err.Error(), httpcodes.StatusUnauthorized)
	}

	return &JwtUserData{
		ID:       valId,
		Role:     valRole,
		Username: valUsername,
	}, nil
}

func (ap *AuthProvider) GenerateUserJwtToken(info JwtUserData, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":       info.ID,
		"username": info.Username,
		"role":     info.Role,
		"exp": time.Now().Add(exp).Unix(),
	})

	tokenEnc, err := token.SignedString(ap.jwtKey)

	return tokenEnc, err
}

func (ap *AuthProvider) ValidateSignature(method ByteEncoding, body, givenBytes []byte) utils.StatusCodeErr {
	digest := sha256.New()

	_, err := digest.Write(body)

	if err != nil {
		return utils.NewStatusCodeErr("failed to hash request body", httpcodes.StatusBadRequest)
	}

	sum := digest.Sum(ap.sigKey)

	var expected string

	if method == ByteEncodingHex {
		expected = hex.EncodeToString(sum)
	} else if method == ByteEncodingBase64 {
		expected = base64.RawStdEncoding.EncodeToString(sum)
	} else {
		return utils.NewStatusCodeErr("invalid encoding method", httpcodes.StatusBadRequest)
	}

	if expected != string(givenBytes) {
		return utils.NewStatusCodeErr("signatures do not match", httpcodes.StatusUnauthorized)
	}

	return nil
}
