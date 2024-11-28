package bcrypt_tools

import (
	"golang.org/x/crypto/bcrypt"
)

func HashToken(data string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	return string(bytes), err
}

func CompareTokenHash(token, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
}
