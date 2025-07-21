package service

import (
	"context"
	"log"
	"task-six-authentication/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskService struct {
	TaskCollection *mongo.Collection
}

func NewTaskService(client *mongo.Client) *TaskService {
	collection := client.Database("taskdb").Collection("tasks")
	return &TaskService{TaskCollection: collection}
}
func (s *TaskService) GetTasks() []models.Task {
	var tasks []models.Task
	curr, err := s.TaskCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	for curr.Next(context.TODO()) {
		var task models.Task

		err = curr.Decode(&task)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	if err := curr.Err(); err != nil {
		log.Fatal(err)
	}
	curr.Close(context.TODO())
	return tasks
}

func (s *TaskService) GetTask(taskID primitive.ObjectID) (models.Task, error) {
	var task models.Task
	filter := bson.D{{Key: "_id", Value: taskID}}
	err := s.TaskCollection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (s *TaskService) AddTask(newTask models.Task) error {
	_, err := s.TaskCollection.InsertOne(context.TODO(), newTask)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) RemoveTask(taskID primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: taskID}}
	_, err := s.TaskCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) UpdateTask(taskID primitive.ObjectID, updatedTask models.Task) error {
	filter := bson.D{{Key: "_id", Value: taskID}}
	update := bson.D{{Key: "$set", Value: updatedTask}}
	_, err := s.TaskCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
