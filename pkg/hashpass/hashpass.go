package hashpass

import (
	"crypto/sha256"
	"crypto/subtle"

	"golang.org/x/crypto/bcrypt"
)

const DefaultCost = 12
const Size = 32

func CheckBcryptPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CreateBcryptHashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CreateBytesHashSHA256(token string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hasher.Sum(nil)
}

func CheckHashSHA256(hashedToken []byte, token string) bool {
	return subtle.ConstantTimeCompare(CreateBytesHashSHA256(token), hashedToken) == 1
}

func CreateBytesHashSHA256WithSalt(salt []byte, token string) []byte {
	hasher := sha256.New()
	hasher.Write(salt)
	hasher.Write([]byte(token))
	return hasher.Sum(nil)
}
