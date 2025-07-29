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

type taskRepository struct {
    taskCollection *mongo.Collection
}

//real constructor 
func NewTaskRepository() Domain.ITaskRepository {
    uri := os.Getenv("MONGODB_URI")
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        panic(err)
    }
    collection := client.Database("task_db").Collection("tasks")
    return &taskRepository{taskCollection: collection}
}

//test constructor inject scollection for memongo
func NewTaskRepositoryWithCollection(collection *mongo.Collection) Domain.ITaskRepository {
    return &taskRepository{taskCollection: collection}
}

func (r *taskRepository) Create(task Domain.Task) (Domain.Task, error) {
   
    doc := bson.M{
        "title":       task.Title,
        "description": task.Description,
        "due_date":    task.DueDate,
        "status":      task.Status,
    }

    res, err := r.taskCollection.InsertOne(context.TODO(), doc)
    if err != nil {
        return Domain.Task{}, err
    }

    oid, ok := res.InsertedID.(primitive.ObjectID)
    if !ok {
        return Domain.Task{}, errors.New("failed to convert inserted ID to ObjectID")
    }

    task.ID = oid.Hex()
    return task, nil
}

func (r *taskRepository) GetByID(id string) (Domain.Task, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return Domain.Task{}, err
    }

    filter := bson.M{"_id": objectID}
    var doc bson.M
    err = r.taskCollection.FindOne(context.TODO(), filter).Decode(&doc)
    if err != nil {
        return Domain.Task{}, err
    }

    task := Domain.Task{
        ID:          id,
        Title:       doc["title"].(string),
        Description: doc["description"].(string),
        Status:      doc["status"].(string),
    }

    if dueDate, ok := doc["due_date"].(primitive.DateTime); ok {
        task.DueDate = dueDate.Time()
    } else if dueDate, ok := doc["due_date"].(time.Time); ok {
        task.DueDate = dueDate
    }

    return task, nil
}

func (r *taskRepository) GetAll() ([]Domain.Task, error) {
    cursor, err := r.taskCollection.Find(context.TODO(), bson.D{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    var tasks []Domain.Task
    for cursor.Next(context.TODO()) {
        var doc bson.M
        err = cursor.Decode(&doc)
        if err != nil {
            return nil, err
        }

        id := doc["_id"].(primitive.ObjectID).Hex()
        task := Domain.Task{
            ID:          id,
            Title:       doc["title"].(string),
            Description: doc["description"].(string),
            Status:      doc["status"].(string),
        }

        if dueDate, ok := doc["due_date"].(primitive.DateTime); ok {
            task.DueDate = dueDate.Time()
        } else if dueDate, ok := doc["due_date"].(time.Time); ok {
            task.DueDate = dueDate
        }

        tasks = append(tasks, task)
    }

    return tasks, nil
}

func (r *taskRepository) Update(id string, updatedTask Domain.Task) (Domain.Task, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return Domain.Task{}, err
    }

    filter := bson.M{"_id": objectID}
    update := bson.M{
        "$set": bson.M{
            "title":       updatedTask.Title,
            "description": updatedTask.Description,
            "due_date":    updatedTask.DueDate,
            "status":      updatedTask.Status,
        },
    }

    res, err := r.taskCollection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        return Domain.Task{}, err
    }

    if res.MatchedCount == 0 {
        return Domain.Task{}, mongo.ErrNoDocuments
    }

    updatedTask.ID = id
    return updatedTask, nil
}

func (r *taskRepository) Delete(id string) error {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    filter := bson.M{"_id": objectID}
    res, err := r.taskCollection.DeleteOne(context.TODO(), filter)
    if err != nil {
        return err
    }
    if res.DeletedCount == 0 {
        return mongo.ErrNoDocuments
    }
    return nil
}
