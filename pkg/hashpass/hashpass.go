package hashpass

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

const DefaultCost = 12
const Size = 32

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CreateHashPass(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CreateHashSHA256(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	sum := hasher.Sum(nil)
	return hex.EncodeToString(sum)
}

func CheckHashSHA256(hashedToken, token string) bool {
	return hashedToken == CreateHashSHA256(token)
}
