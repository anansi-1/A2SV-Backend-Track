package Infrastructure

import (
	"golang.org/x/crypto/bcrypt"
	"task-manager/Domain"
)

type PasswordService struct{}

func NewPasswordService() Domain.IPasswordService {
	return &PasswordService{}
}

func (p *PasswordService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (p *PasswordService) Compare(plain, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
