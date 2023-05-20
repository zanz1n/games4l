package providers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"
)

type ByteEncoding string

const (
	ByteEncodingBase64 ByteEncoding = "BASE64"
	ByteEncodingHex    ByteEncoding = "HEX"
)

type AuthProvider struct {
	sigKey []byte
}

func NewAuthProvider() *AuthProvider {
	return &AuthProvider{
		sigKey: []byte(GetConfig().WebhookSig),
	}
}

func (ap *AuthProvider) ValidateSignature(method ByteEncoding, body, givenBytes []byte) StatusCodeErr {
	digest := sha256.New()

	_, err := digest.Write(body)

	if err != nil {
		return NewStatusCodeErr("failed to hash request body", fiber.StatusBadRequest)
	}

	sum := digest.Sum(ap.sigKey)

	var expected string

	if method == ByteEncodingHex {
		expected = hex.EncodeToString(sum)
	} else if method == ByteEncodingBase64 {
		expected = base64.RawStdEncoding.EncodeToString(sum)
	} else {
		return NewStatusCodeErr("invalid encoding method", fiber.StatusBadRequest)
	}

	if expected != string(givenBytes) {
		return NewStatusCodeErr("signatures do not match", fiber.StatusUnauthorized)
	}

	return nil
}
