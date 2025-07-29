package Repositories_test

import (
    "context"
    "testing"
    "time"
    "task-manager/Domain"
    "task-manager/Repositories"

    "github.com/stretchr/testify/assert"
    "github.com/tryvium-travels/memongo"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

//setup inmemory mongodb
func setupTestDB(t *testing.T) (*mongo.Collection, func()) {
    t.Helper()

    mongoServer, err := memongo.Start("4.0.5")
    if err != nil {
        t.Fatalf("Failed to start memongo: %v", err)
    }

    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoServer.URI()))
    if err != nil {
        t.Fatalf("Failed to connect memongo: %v", err)
    }

    collection := client.Database("test_task_db").Collection("tasks")

    cleanup := func() {
        _ = client.Disconnect(context.TODO())
        mongoServer.Stop()
    }

    return collection, cleanup
}

func createSampleTask() Domain.Task {
    return Domain.Task{
        Title:       "First Task",
        Description: "this is a task description",
        DueDate:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
        Status:      "pending",
    }
}

func TestTaskRepository_Create_Success(t *testing.T) {
    collection, cleanup := setupTestDB(t)
    defer cleanup()

    repo := Repositories.NewTaskRepositoryWithCollection(collection)
    task := createSampleTask()

    created, err := repo.Create(task)
    assert.NoError(t, err)
    assert.NotEmpty(t, created.ID)
    assert.Equal(t, task.Title, created.Title)
    assert.Equal(t, task.Description, created.Description)
    assert.Equal(t, task.Status, created.Status)
    assert.True(t, task.DueDate.Equal(created.DueDate))

    fetched, err := repo.GetByID(created.ID)
    assert.NoError(t, err)
    assert.Equal(t, created.ID, fetched.ID)
}

func TestTaskRepository_GetByID_Success(t *testing.T) {
    collection, cleanup := setupTestDB(t)
    defer cleanup()

    repo := Repositories.NewTaskRepositoryWithCollection(collection)

    task := createSampleTask()
    created, _ := repo.Create(task)

    fetched, err := repo.GetByID(created.ID)
    assert.NoError(t, err)
    assert.Equal(t, created.ID, fetched.ID)
    assert.Equal(t, created.Title, fetched.Title)
}

func TestTaskRepository_GetByID_NotFound(t *testing.T) {
    collection, cleanup := setupTestDB(t)
    defer cleanup()

    repo := Repositories.NewTaskRepositoryWithCollection(collection)

    fakeID := "507f1f77bcf86cd799439011" 
    _, err := repo.GetByID(fakeID)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "no documents")
}

func TestTaskRepository_GetByID_InvalidID(t *testing.T) {
    collection, cleanup := setupTestDB(t)
    defer cleanup()

    repo := Repositories.NewTaskRepositoryWithCollection(collection)

    _, err := repo.GetByID("invalid-id")
    assert.Error(t, err)
}

func TestTaskRepository_GetAll_Success(t *testing.T) {
    collection, cleanup := setupTestDB(t)
    defer cleanup()

    repo := Repositories.NewTaskRepositoryWithCollection(collection)

    task1 := createSampleTask()
    task2 := Domain.Task{
        Title:  "Another Task",
        Status: "completed",
        DueDate: time.Now(),
    }

    _, _ = repo.Create(task1)
    _, _ = repo.Create(task2)

    tasks, err := repo.GetAll()
    assert.NoError(t, err)
    assert.Len(t, tasks, 2)

    titles := []string{tasks[0].Title, tasks[1].Title}
    assert.Contains(t, titles, task1.Title)
    assert.Contains(t, titles, task2.Title)
}

func TestTaskRepository_Update_Success(t *testing.T) {
    collection, cleanup := setupTestDB(t)
    defer cleanup()

    repo := Repositories.NewTaskRepositoryWithCollection(collection)

    original := createSampleTask()
    created, _ := repo.Create(original)

    update := Domain.Task{
        Title:       "Updated Task",
        Description: "Updated Desc",
        Status:      "completed",
        DueDate:     time.Now(),
    }

    updatedTask, err := repo.Update(created.ID, update)
    assert.NoError(t, err)
    assert.Equal(t, created.ID, updatedTask.ID)
    assert.Equal(t, update.Title, updatedTask.Title)

    fetched, _ := repo.GetByID(created.ID)
    assert.Equal(t, update.Title, fetched.Title)
    assert.Equal(t, update.Status, fetched.Status)
}

func TestTaskRepository_Update_NotFound(t *testing.T) {
    collection, cleanup := setupTestDB(t)
    defer cleanup()

    repo := Repositories.NewTaskRepositoryWithCollection(collection)

    update := createSampleTask()
    fakeID := "507f1f77bcf86cd799439011"

    _, err := repo.Update(fakeID, update)
    assert.Error(t, err)
    assert.Equal(t, mongo.ErrNoDocuments, err)
}

func TestTaskRepository_Delete_Success(t *testing.T) {
    collection, cleanup := setupTestDB(t)
    defer cleanup()

    repo := Repositories.NewTaskRepositoryWithCollection(collection)

    task := createSampleTask()
    created, _ := repo.Create(task)

    err := repo.Delete(created.ID)
    assert.NoError(t, err)

    _, err = repo.GetByID(created.ID)
    assert.Error(t, err)
}

func TestTaskRepository_Delete_NotFound(t *testing.T) {
    collection, cleanup := setupTestDB(t)
    defer cleanup()

    repo := Repositories.NewTaskRepositoryWithCollection(collection)

    fakeID := "507f1f77bcf86cd799439011"

    err := repo.Delete(fakeID)
    assert.Error(t, err)
    assert.Equal(t, mongo.ErrNoDocuments, err)
}
