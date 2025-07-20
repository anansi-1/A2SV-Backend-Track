package data

import (
	"context"
	"errors"

	"github/anansi-1/Task-Four-Task-Manager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type TaskServiceImpl struct {
	taskCollection *mongo.Collection
	ctx            context.Context
}

func NewTaskService(taskCollection *mongo.Collection, ctx context.Context) TaskService {
	return &TaskServiceImpl{
		taskCollection: taskCollection,
		ctx:            ctx,
	}
}

func (s *TaskServiceImpl) GetTasks() ([]*models.Task, error) {
	var tasks []*models.Task

	cursor, err := s.taskCollection.Find(s.ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(s.ctx)

	for cursor.Next(s.ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, errors.New("no tasks found")
	}
	return tasks, nil
}

func (s *TaskServiceImpl) GetTask(id string) (*models.Task, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid task ID")
	}

	var task models.Task
	err = s.taskCollection.FindOne(s.ctx, bson.M{"_id": objID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return &task, nil
}

func (s *TaskServiceImpl) AddTask(task *models.Task) error {
	task.ID = primitive.NewObjectID().Hex()
	_, err := s.taskCollection.InsertOne(s.ctx, task)
	return err
}

func (s *TaskServiceImpl) UpdateTask(id string, updatedTask *models.Task) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID")
	}

	update := bson.M{}
	if updatedTask.Title != "" {
		update["title"] = updatedTask.Title
	}
	if updatedTask.Description != "" {
		update["description"] = updatedTask.Description
	}
	if !updatedTask.DueDate.IsZero() {
		update["dueDate"] = updatedTask.DueDate
	}
	if updatedTask.Status != "" {
		update["status"] = updatedTask.Status
	}

	if len(update) == 0 {
		return errors.New("no fields to update")
	}

	result, err := s.taskCollection.UpdateOne(s.ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (s *TaskServiceImpl) RemoveTask(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID")
	}

	result, err := s.taskCollection.DeleteOne(s.ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
