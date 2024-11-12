package change

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

func ToBase64(input string) string {

	utf8Bytes := []byte(input)

	encodedStr := base64.StdEncoding.EncodeToString(utf8Bytes)
	return encodedStr
}

func ToUnicode(encodedStr string) (string, error) {

	decodedBytes, err := base64.StdEncoding.DecodeString(encodedStr)
	if err != nil {
		return "", err
	}

	decodedStr := string(decodedBytes)
	return decodedStr, nil
}

func GetPublicKey(input_base64 string) (string, error) {
	var err error
	var privateKey, publicKey []byte
	var encoding *base64.Encoding

	encoding = base64.RawURLEncoding

	if len(input_base64) > 0 {
		privateKey, err = encoding.DecodeString(input_base64)
		if err != nil {
			return "", err
		}
		if len(privateKey) != curve25519.ScalarSize {
			return "", fmt.Errorf("Invalid length of private key.")
		}
	}

	if privateKey == nil {
		return "", fmt.Errorf("Invalid length of private key.")
	}

	// Modify random bytes using algorithm described at:
	// https://cr.yp.to/ecdh.html.
	privateKey[0] &= 248
	privateKey[31] &= 127 | 64

	if publicKey, err = curve25519.X25519(privateKey, curve25519.Basepoint); err != nil {
		return "", err

	}

	return encoding.EncodeToString(publicKey), nil

}
