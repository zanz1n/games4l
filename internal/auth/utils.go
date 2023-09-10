package auth

import (
	"strings"

	"github.com/games4l/internal/errors"
)

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
