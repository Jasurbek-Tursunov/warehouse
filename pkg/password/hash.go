package password

import (
	"golang.org/x/crypto/bcrypt"
)

func HashingPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func AssertPassword(password, target string) bool {
	return bcrypt.CompareHashAndPassword([]byte(target), []byte(password)) == nil
}
