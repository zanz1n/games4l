package auth

import (
	"strings"

	"github.com/games4l/pkg/errors"
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

type AuthMethod string

const (
	AuthMethodBearer    AuthMethod = "Bearer"
	AuthMethodSignature AuthMethod = "Signature"
)

var ValidUserRoles = []UserRole{UserRoleAdmin, UserRoleClient, UserRolePacient}

type JwtUserData struct {
	Role     UserRole `json:"role"`
	Username string   `json:"username"`
	ID       string   `json:"id"`
}

type AuthInfo struct {
	Method   AuthMethod
	Encoding ByteEncoding
	Payload  string
}

func ExtractAuthHeaderInfo(header string) (*AuthInfo, error) {
	headerS := strings.Split(header, " ")
	if len(headerS) < 2 {
		return nil, errors.ErrRouteRequiresAdminAuth
	}

	info := AuthInfo{
		Method: AuthMethod(headerS[0]),
	}

	if info.Method == AuthMethodSignature {
		info.Encoding = ByteEncoding(headerS[1])

		switch info.Encoding {
		case ByteEncodingBase64:
		case ByteEncodingHex:
		default:
			return nil, errors.ErrInvalidAuthSignatureEncodingMethod
		}

		info.Payload = headerS[2]
	} else if info.Method == AuthMethodBearer {
		info.Encoding = ""

		info.Payload = headerS[1]
	} else {
		return nil, errors.ErrInvalidAuthStrategy
	}

	return &info, nil
}
