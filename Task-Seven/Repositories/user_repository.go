package repositories

import (
	"context"
	"errors"
	domain "task-seven/Domain"
	infrastructure "task-seven/Infrastructure"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	database   mongo.Database
	collection string
}

func NewUserRepository(db mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

func (ur *userRepository) Register(c context.Context, user *domain.User) error {
	collection := ur.database.Collection(ur.collection)
	filter := bson.D{{Key: "username", Value: user.Username}}
	var existing_user domain.User
	err := collection.FindOne(c, filter).Decode(&existing_user)
	if err == nil {
		return errors.New("the username is already used")
	}

	hashedPassword, err2 := infrastructure.HashPassword(user.Password)
	if err2 != nil {
		return err2
	}
	user.Password = string(hashedPassword)
	if collection == nil {
		user.Role = "admin"
	} else {
		user.Role = "regular"
	}
	_, err3 := collection.InsertOne(c, user)
	if err3 != nil {
		return err3
	}
	return nil
}

func (ur *userRepository) Login(c context.Context, user *domain.User) (string, error) {
	collection := ur.database.Collection(ur.collection)
	var existing_user domain.User
	filter := bson.D{{Key: "username", Value: user.Username}}
	err := collection.FindOne(c, filter).Decode(&existing_user)
	if err != nil {
		return "", errors.New("no user is registerd with this user name")
	}

	if err2 := infrastructure.CheckPassword(existing_user.Password, user.Password); err2 != nil {
		return "", errors.New("the error is while comparing the password")
	}

	jwt_token, err3 := infrastructure.GenerateToken(existing_user)
	if err3 != nil {
		return "", err3
	}
	return jwt_token, nil
}

func (ur *userRepository) Promote(c context.Context, userID primitive.ObjectID) error {
	collection := ur.database.Collection(ur.collection)
	var existing_user domain.User
	filter := bson.D{{Key: "_id", Value: userID}}
	if err := collection.FindOne(c, filter).Decode(&existing_user); err != nil {
		return errors.New("user not found")
	}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "role", Value: "admin"},
	}}}
	if _, err := collection.UpdateOne(c, filter, update); err != nil {
		return err
	}
	return nil
}
