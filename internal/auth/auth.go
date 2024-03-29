package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"fmt"
	"time"

	"github.com/games4l/internal/logger"
	"github.com/games4l/internal/utils"
	"github.com/games4l/pkg/errors"
	"github.com/golang-jwt/jwt/v4"
)

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

func (ap *AuthProvider) AuthUser(payload string) (*JwtUserData, error) {
	var (
		valId       string
		valUsername string
		valRole     string
		err         error = nil
	)
	token, _ := jwt.Parse(payload, func(t *jwt.Token) (interface{}, error) {
		var (
			ok     bool
			claims jwt.MapClaims
		)

		if _, ok = t.Method.(*jwt.SigningMethodHMAC); !ok {
			headerAlg, ok := t.Header["alg"].(string)
			if !ok {
				headerAlg = "<IMPOSSIBLE TO DECODE>"
			}

			return nil, fmt.Errorf("invalid auth token signing method: %s", headerAlg)
		}

		if claims, ok = t.Claims.(jwt.MapClaims); !ok {
			logger.Error("Jwt token invalidation: is not jwt.MapClaims")
			err = errors.ErrInvalidJwtTokenFormat
			return nil, errors.ErrInvalidJwtTokenFormat
		}

		if valId, ok = claims["id"].(string); !ok {
			logger.Error("Jwt token invalidation: does not contain an id")
			err = errors.ErrInvalidJwtTokenFormat
			return nil, errors.ErrInvalidJwtTokenFormat
		}

		if valUsername, ok = claims["username"].(string); !ok {
			logger.Error("Jwt token invalidation: does not contain an username")
			err = errors.ErrInvalidJwtTokenFormat
			return nil, errors.ErrInvalidJwtTokenFormat
		}

		if valRole, ok = claims["role"].(string); !ok {
			logger.Error("Jwt token invalidation: does not contain a role")
			err = errors.ErrInvalidJwtTokenFormat
			return nil, errors.ErrInvalidJwtTokenFormat
		}

		if !utils.SliceContains(ValidUserRoles, UserRole(valRole)) {
			logger.Error("Jwt token invalidation: the role is not valid")
			err = errors.ErrInvalidJwtTokenFormat
			return nil, errors.ErrInvalidJwtTokenFormat
		}

		if t.Claims.(jwt.MapClaims)["exp"].(float64) < float64(time.Now().Unix()) {
			err = errors.ErrJwtTokenExpired
			return nil, errors.ErrJwtTokenExpired
		}

		return []byte(ap.jwtKey), nil
	})

	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, errors.ErrInvalidJwtTokenFormat
	}

	return &JwtUserData{
		ID:       valId,
		Role:     UserRole(valRole),
		Username: valUsername,
	}, nil
}

func (ap *AuthProvider) GenerateUserJwtToken(info JwtUserData, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":       info.ID,
		"username": info.Username,
		"role":     info.Role,
		"exp":      time.Now().Add(exp).Unix(),
	})

	tokenEnc, err := token.SignedString(ap.jwtKey)

	return tokenEnc, err
}

func (ap *AuthProvider) ValidateSignature(method ByteEncoding, body, givenBytes []byte) error {
	var err error

	digest := sha256.New()

	if _, err = digest.Write(body); err != nil {
		return errors.ErrMalformedOrTooBigBody
	}
	if _, err = digest.Write(ap.sigKey); err != nil {
		return errors.ErrMalformedOrTooBigBody
	}

	sum := digest.Sum([]byte{})

	var expected string

	switch method {
	case ByteEncodingHex:
		expected = hex.EncodeToString(sum)
	case ByteEncodingBase64:
		expected = base64.RawStdEncoding.EncodeToString(sum)
	default:
		return errors.ErrInvalidAuthSignatureEncodingMethod
	}

	if expected != string(givenBytes) {
		return errors.ErrInvalidAuthSignature
	}

	return nil
}

func (ap *AuthProvider) AuthenticateUserHeader(header string) (*JwtUserData, error) {
	info, err := ExtractAuthHeaderInfo(header)
	if err != nil {
		return nil, err
	}

	if info.Method != AuthMethodBearer {
		return nil, errors.ErrInvalidAuthStrategy
	}

	return ap.AuthUser(info.Payload)
}

func (ap *AuthProvider) AuthenticateAdminHeader(header string, body []byte) error {
	info, err := ExtractAuthHeaderInfo(header)
	if err != nil {
		return err
	}

	switch info.Method {
	case AuthMethodBearer:
		var user *JwtUserData

		if user, err = ap.AuthUser(info.Payload); err != nil {
			return err
		}

		if user.Role != UserRoleAdmin {
			return errors.ErrRouteRequiresAdminAuth
		}
	case AuthMethodSignature:
		if err = ap.ValidateSignature(info.Encoding, body, utils.S2B(info.Payload)); err != nil {
			return err
		}
	default:
		return errors.ErrInvalidAuthStrategy
	}

	return nil
}
