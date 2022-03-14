package key

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateHS256Key() (string, error) {
	key := make([]byte, 256/8) // 256 bits
	return generateHSKey(key)
}

func GenerateHS384Key() (string, error) {
	key := make([]byte, 384/8) // 384 bits
	return generateHSKey(key)
}

func GenerateHS512Key() (string, error) {
	key := make([]byte, 512/8) // 512 bits
	return generateHSKey(key)
}

func generateHSKey(key []byte) (string, error) {
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	ks := base64.RawURLEncoding.EncodeToString(key) // default no padding
	return fmt.Sprintf("%s", ks), nil
}

func ToBinaryRunes(s string) string {
	var buffer bytes.Buffer
	for _, runeValue := range s {
		fmt.Fprintf(&buffer, "%b", runeValue)
	}
	return fmt.Sprintf("%s", buffer.Bytes())
}
