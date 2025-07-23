package infrastructure

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func CheckPassword(password1, password2 string) error {
	return bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
}
