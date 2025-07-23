package usecases

import (
	"context"
	domain "task-seven/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userUsecase struct {
	userReposistory domain.UserRepository
	contextTimeout  time.Duration
}

func NewUserUsecase(userRepo domain.UserRepository, ct time.Duration) domain.UserRepository {
	return &userUsecase{
		userReposistory: userRepo,
		contextTimeout:  ct,
	}
}

func (uu *userUsecase) Register(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userReposistory.Register(ctx, user)
}

func (uu *userUsecase) Login(c context.Context, user *domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userReposistory.Login(ctx, user)

}

func (uu *userUsecase) Promote(c context.Context, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userReposistory.Promote(ctx, userID)
}
