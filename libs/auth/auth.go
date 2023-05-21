package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
)

type ByteEncoding string

const (
	ByteEncodingBase64 ByteEncoding = "BASE64"
	ByteEncodingHex    ByteEncoding = "HEX"
)

type AuthProvider struct {
	sigKey []byte
}

func NewAuthProvider(sigKey []byte) *AuthProvider {
	return &AuthProvider{
		sigKey: sigKey,
	}
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
