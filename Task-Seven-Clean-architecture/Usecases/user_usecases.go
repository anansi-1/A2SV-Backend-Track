package Usecases

import (
	"errors"
	"task-manager/Domain"
)

type UserUsecase struct {
	UserRepo       Domain.IUserRepository
	PasswordHasher Domain.IPasswordService
	JWTService     Domain.IJWTService
}

func NewUserUsecase(repo Domain.IUserRepository, hasher Domain.IPasswordService, jwt Domain.IJWTService) *UserUsecase {
	return &UserUsecase{
		UserRepo:       repo,
		PasswordHasher: hasher,
		JWTService:     jwt,
	}
}

func (u *UserUsecase) Register(user Domain.User) (Domain.User, error) {
	_, err := u.UserRepo.FindByEmail(user.Email)
	if err == nil {
		return Domain.User{}, errors.New("email already registered")
	}

	hashedPassword, err := u.PasswordHasher.Hash(user.Password)
	if err != nil {
		return Domain.User{}, errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	createdUser, err := u.UserRepo.Create(user)
	if err != nil {
		return Domain.User{}, err
	}
	return createdUser, nil
}

func (u *UserUsecase) Login(email, password string) (string, error) {
	user, err := u.UserRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !u.PasswordHasher.Compare(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := u.JWTService.GenerateToken(user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (u *UserUsecase) PromoteUser(id string) (Domain.User, error) {
	return u.UserRepo.Promote(id)
}
