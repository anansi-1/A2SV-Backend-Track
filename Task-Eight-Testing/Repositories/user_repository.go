package Repositories

import (
    "context"
    "errors"
    "os"
    "task-manager/Domain"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)



type userRepository struct {
    userCollection *mongo.Collection
}

func NewUserRepository() Domain.IUserRepository {
    uri := os.Getenv("MONGODB_URI")
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        panic(err)
    }

    userCollection := client.Database("task_db").Collection("user")
    return &userRepository{userCollection: userCollection}
}


func (r *userRepository) FindByEmail(email string) (Domain.User, error) {
    var result bson.M
    err := r.userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&result)
    if err != nil {
        return Domain.User{}, err
    }

    id, ok := result["_id"].(primitive.ObjectID)
    if !ok {
        return Domain.User{}, errors.New("invalid ID type")
    }

    return Domain.User{
        ID:       id.Hex(),
        Email:    result["email"].(string),
        Password: result["password"].(string),
        Role:     result["role"].(string),
    }, nil
}

func (r *userRepository) Create(user Domain.User) (Domain.User, error) {
    _, err := r.FindByEmail(user.Email)
    if err == nil {
        return Domain.User{}, errors.New("email already registered")
    }

    count, err := r.userCollection.CountDocuments(context.TODO(), bson.D{})
    if err != nil {
        return Domain.User{}, errors.New("failed to check user count")
    }

    if count == 0 {
        user.Role = "admin"
    } else {
        user.Role = "user"
    }

    userObjectID := primitive.NewObjectID()
    user.ID = userObjectID.Hex()

    doc := bson.M{
        "_id":      userObjectID,
        "email":    user.Email,
        "password": user.Password,
        "role":     user.Role,
    }

    _, err = r.userCollection.InsertOne(context.TODO(), doc)
    if err != nil {
        return Domain.User{}, errors.New("failed to insert user")
    }

    return user, nil
}

func (r *userRepository) Promote(id string) (Domain.User, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return Domain.User{}, err
    }

    filter := bson.M{"_id": objectID}
    update := bson.M{"$set": bson.M{"role": "admin"}}

    _, err = r.userCollection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        return Domain.User{}, err
    }

    var updated bson.M
    err = r.userCollection.FindOne(context.TODO(), filter).Decode(&updated)
    if err != nil {
        return Domain.User{}, err
    }

    oid, ok := updated["_id"].(primitive.ObjectID)
    if !ok {
        return Domain.User{}, errors.New("invalid ID type after update")
    }

    return Domain.User{
        ID:       oid.Hex(),
        Email:    updated["email"].(string),
        Password: updated["password"].(string),
        Role:     updated["role"].(string),
    }, nil
}
