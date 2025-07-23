package repositories

import (
	"context"
	"log"
	domain "task-seven/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	database   mongo.Database
	collection string
}

func NewTaskRepository(db mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}

func (tr *taskRepository) GetTasks(c context.Context) []domain.Task {
	collection := tr.database.Collection(tr.collection)
	var tasks []domain.Task
	curr, err := collection.Find(c, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	for curr.Next(c) {
		var task domain.Task

		err = curr.Decode(&task)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	if err := curr.Err(); err != nil {
		log.Fatal(err)
	}
	curr.Close(c)
	return tasks
}
func (tr *taskRepository) GetTask(c context.Context, taskID primitive.ObjectID) (domain.Task, error) {
	collection := tr.database.Collection(tr.collection)
	var task domain.Task
	filter := bson.D{{Key: "_id", Value: taskID}}
	err := collection.FindOne(c, filter).Decode(&task)
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (tr *taskRepository) AddTask(c context.Context, task *domain.Task) error {
	collection := tr.database.Collection(tr.collection)
	_, err := collection.InsertOne(c, task)
	if err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) RemoveTask(c context.Context, taskID primitive.ObjectID) error {
	collection := tr.database.Collection(tr.collection)
	filter := bson.D{{Key: "_id", Value: taskID}}
	_, err := collection.DeleteOne(c, filter)

	if err != nil {
		return err
	}
	return nil

}

func (tr *taskRepository) UpdateTask(c context.Context, taskID primitive.ObjectID, updatedTask *domain.Task) error {
	collection := tr.database.Collection(tr.collection)
	filter := bson.D{{Key: "_id", Value: taskID}}
	update := bson.D{{Key: "$set", Value: updatedTask}}
	_, err := collection.UpdateOne(c, filter, update)
	if err != nil {
		return err
	}
	return nil

}
