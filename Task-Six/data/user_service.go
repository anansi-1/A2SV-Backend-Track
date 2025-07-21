package service

import (
	"context"
	"errors"
	"task-six-authentication/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserCollection *mongo.Collection
}

func NewUserService(client *mongo.Client) *UserService {
	collection := client.Database("taskdb").Collection("users")
	return &UserService{UserCollection: collection}
}

func (s *UserService) Register(newUser models.User) error {

	filter := bson.D{{Key: "username", Value: newUser.Username}}
	var existing_user models.User
	err := s.UserCollection.FindOne(context.TODO(), filter).Decode(&existing_user)
	if err == nil {
		return errors.New("the username is already used")
	}

	hashedPassword, err2 := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err2 != nil {
		return err2
	}
	newUser.Password = string(hashedPassword)
	if s.UserCollection == nil {
		newUser.Role = "admin"
	} else {
		newUser.Role = "regular"
	}
	_, err3 := s.UserCollection.InsertOne(context.TODO(), newUser)
	if err3 != nil {
		return err3
	}
	return nil
}

func (s *UserService) Login(unauthenticatedUser models.User) (string, error) {
	jwt_secret := []byte("It has to be a secret")

	var existing_user models.User
	filter := bson.D{{Key: "username", Value: unauthenticatedUser.Username}}
	err := s.UserCollection.FindOne(context.TODO(), filter).Decode(&existing_user)
	if err != nil {
		return "", errors.New("no user is registerd with this user name")
	}

	if err2 := bcrypt.CompareHashAndPassword([]byte(existing_user.Password), []byte(unauthenticatedUser.Password)); err2 != nil {
		return "", errors.New("the error is while comparing the password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
	token.Claims = jwt.MapClaims{
		"_id":      existing_user.ID,
		"username": existing_user.Username,
		"role":     existing_user.Role,
	}

	jwt_token, err3 := token.SignedString(jwt_secret)
	if err3 != nil {
		return "", err3
	}
	return jwt_token, nil
}

func (s *UserService) Promote(regular_user_id primitive.ObjectID) error {
	var existing_user models.User
	filter := bson.D{{Key: "_id", Value: regular_user_id}}
	if err := s.UserCollection.FindOne(context.TODO(), filter).Decode(&existing_user); err != nil {
		return errors.New("user not found")
	}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "role", Value: "admin"},
	}}}
	if _, err := s.UserCollection.UpdateOne(context.TODO(), filter, update); err != nil {
		return err
	}
	return nil
}
