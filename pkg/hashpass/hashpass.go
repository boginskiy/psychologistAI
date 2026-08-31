package hashpass

import "golang.org/x/crypto/bcrypt"

const DefaultCost = 12

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
