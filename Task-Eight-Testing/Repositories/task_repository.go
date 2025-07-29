package Repositories

import (
    "context"
    "errors"
    "os"
    "task-manager/Domain"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type ITaskRepository interface {
    GetAll() ([]Domain.Task, error)
    GetByID(id string) (Domain.Task, error)
    Create(task Domain.Task) (Domain.Task, error)
    Update(id string, task Domain.Task) (Domain.Task, error)
    Delete(id string) error
}

type taskRepository struct {
    taskCollection *mongo.Collection
}

func NewTaskRepository() ITaskRepository {
    uri := os.Getenv("MONGODB_URI")
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        panic(err)
    }

    taskCollection := client.Database("task_db").Collection("tasks")
    return &taskRepository{taskCollection: taskCollection}
}

func (r *taskRepository) GetAll() ([]Domain.Task, error) {
    cursor, err := r.taskCollection.Find(context.TODO(), bson.D{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    var tasks []Domain.Task
    for cursor.Next(context.TODO()) {
        var taskDoc bson.M
        err := cursor.Decode(&taskDoc)
        if err != nil {
            return nil, err
        }

        id, ok := taskDoc["_id"].(primitive.ObjectID)
        if !ok {
            return nil, errors.New("invalid ID type")
        }

        dueDate, ok := taskDoc["due_date"].(primitive.DateTime)
        if !ok {
            dueDate = primitive.DateTime(time.Now().Unix() * 1000)
        }

        task := Domain.Task{
            ID:          id.Hex(),
            Title:       taskDoc["title"].(string),
            Description: taskDoc["description"].(string),
            DueDate:     dueDate.Time(),
            Status:      taskDoc["status"].(string),
        }
        tasks = append(tasks, task)
    }

    return tasks, nil
}

func (r *taskRepository) GetByID(id string) (Domain.Task, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return Domain.Task{}, err
    }

    filter := bson.M{"_id": objectID}
    var taskDoc bson.M
    err = r.taskCollection.FindOne(context.TODO(), filter).Decode(&taskDoc)
    if err != nil {
        return Domain.Task{}, err
    }

    oid, ok := taskDoc["_id"].(primitive.ObjectID)
    if !ok {
        return Domain.Task{}, errors.New("invalid ID type")
    }

    dueDate, ok := taskDoc["due_date"].(primitive.DateTime)
    if !ok {
        dueDate = primitive.DateTime(time.Now().Unix() * 1000) 
    }

    task := Domain.Task{
        ID:          oid.Hex(),
        Title:       taskDoc["title"].(string),
        Description: taskDoc["description"].(string),
        DueDate:     dueDate.Time(),
        Status:      taskDoc["status"].(string),
    }
    return task, nil
}

func (r *taskRepository) Create(task Domain.Task) (Domain.Task, error) {
    newID := primitive.NewObjectID()

    doc := bson.M{
        "_id":         newID,
        "title":       task.Title,
        "description": task.Description,
        "due_date":    task.DueDate,
        "status":      task.Status,
    }

    _, err := r.taskCollection.InsertOne(context.TODO(), doc)
    if err != nil {
        return Domain.Task{}, err
    }

    task.ID = newID.Hex()
    return task, nil
}

func (r *taskRepository) Update(id string, task Domain.Task) (Domain.Task, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return Domain.Task{}, err
    }

    filter := bson.M{"_id": objectID}
    update := bson.M{
        "$set": bson.M{
            "title":       task.Title,
            "description": task.Description,
            "due_date":    task.DueDate,
            "status":      task.Status,
        },
    }

    _, err = r.taskCollection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        return Domain.Task{}, err
    }

    return r.GetByID(id)
}

func (r *taskRepository) Delete(id string) error {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    _, err = r.taskCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
    return err
}
