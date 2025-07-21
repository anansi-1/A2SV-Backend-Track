package data

import (
	"context"
	"errors"

	"github/anansi-1/Task-Four-Task-Manager/models"
	"go.mongodb.org/mongo-driver/bson"
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

	c, err := s.taskCollection.Find(s.ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer c.Close(s.ctx)

	for c.Next(s.ctx) {
		var task models.Task
		if err := c.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := c.Err(); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, errors.New("no tasks found")
	}
	return tasks, nil
}

func (s *TaskServiceImpl) GetTask(id string) (*models.Task, error) {
	var task models.Task
	err := s.taskCollection.FindOne(s.ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &task, nil
}

func (s *TaskServiceImpl) AddTask(task *models.Task) error {
	if task.ID == "" {
		return errors.New("task ID is required")
	}
	_, err := s.taskCollection.InsertOne(s.ctx, task)
	return err
}

func (s *TaskServiceImpl) UpdateTask(id string, updatedTask *models.Task) error {
	update := bson.M{}
	if updatedTask.Title != "" {
		update["title"] = updatedTask.Title
	}
	if updatedTask.Description != "" {
		update["description"] = updatedTask.Description
	}
	if !updatedTask.DueDate.IsZero() {
		update["due_date"] = updatedTask.DueDate
	}
	if updatedTask.Status != "" {
		update["status"] = updatedTask.Status
	}

	if len(update) == 0 {
		return errors.New("no fields to update")
	}

	result, err := s.taskCollection.UpdateOne(s.ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (s *TaskServiceImpl) RemoveTask(id string) error {
	result, err := s.taskCollection.DeleteOne(s.ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
