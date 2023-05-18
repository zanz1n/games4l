package providers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
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

func (ap *AuthProvider) ValidateSignature(method ByteEncoding, body, givenBytes []byte) (err error) {
	digest := sha256.New()

	_, err = digest.Write(body)

	if err != nil {
		return
	}

	sum := digest.Sum(ap.sigKey)

	var expected string

	if method == ByteEncodingHex {
		expected = hex.EncodeToString(sum)
	} else if method == ByteEncodingBase64 {
		expected = base64.RawStdEncoding.EncodeToString(sum)
	} else {
		err = errors.New("invalid encoding method")
		return
	}

	if expected != string(givenBytes) {
		err = errors.New("signatures do not match")
	}

	return
}
