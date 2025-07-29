package Repositories_test

import (
    "context"
    "os"
    "testing"
    "time"
    "task-manager/Domain"
    "task-manager/Repositories"

    "github.com/stretchr/testify/assert"
    "github.com/tryvium-travels/memongo"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func setupUserTestDB(t *testing.T) (*mongo.Collection, func()) {
    t.Helper()

    mongoServer, err := memongo.StartWithOptions(&memongo.Options{
        MongoVersion: "4.0.5",
        StartupTimeout: 30 * time.Second,
    })
    if err != nil {
        t.Fatalf("Failed to start memongo: %v", err)
    }

    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoServer.URI()))
    if err != nil {
        t.Fatalf("Failed to connect to memongo: %v", err)
    }

    os.Setenv("MONGODB_URI", mongoServer.URI())

    collection := client.Database("task_db").Collection("user")

    cleanup := func() {
        _ = client.Disconnect(context.TODO())
        mongoServer.Stop()
    }

    return collection, cleanup
}

func createTestUser(email string) Domain.User {
    return Domain.User{
        Email:    email,
        Password: "securepassword123",
    }
}

func TestUserRepository_Create_AssignsRolesCorrectly(t *testing.T) {
    _, cleanup := setupUserTestDB(t)
    defer cleanup()

    repo := Repositories.NewUserRepository()

    user1 := createTestUser("admin@example.com")
    created1, err := repo.Create(user1)
    assert.NoError(t, err)
    assert.Equal(t, "admin", created1.Role)

    user2 := createTestUser("user@example.com")
    created2, err := repo.Create(user2)
    assert.NoError(t, err)
    assert.Equal(t, "user", created2.Role)
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
    _, cleanup := setupUserTestDB(t)
    defer cleanup()

    repo := Repositories.NewUserRepository()
    user := createTestUser("duplicate@example.com")

    _, err := repo.Create(user)
    assert.NoError(t, err)

    _, err = repo.Create(user)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "email already registered")
}

func TestUserRepository_FindByEmail_Success(t *testing.T) {
    _, cleanup := setupUserTestDB(t)
    defer cleanup()

    repo := Repositories.NewUserRepository()
    user := createTestUser("findme@example.com")

    created, _ := repo.Create(user)

    fetched, err := repo.FindByEmail(user.Email)
    assert.NoError(t, err)
    assert.Equal(t, created.ID, fetched.ID)
    assert.Equal(t, user.Email, fetched.Email)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
    _, cleanup := setupUserTestDB(t)
    defer cleanup()

    repo := Repositories.NewUserRepository()

    _, err := repo.FindByEmail("nonexistent@example.com")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "mongo")
}

func TestUserRepository_Promote_ChangesRole(t *testing.T) {
    _, cleanup := setupUserTestDB(t)
    defer cleanup()

    repo := Repositories.NewUserRepository()
    user := createTestUser("promote@example.com")

    created, _ := repo.Create(user)
    assert.Equal(t, "admin", created.Role) 
    _, _ = repo.Create(createTestUser("second@example.com")) 
    promoted, err := repo.Promote(created.ID)
    assert.NoError(t, err)
    assert.Equal(t, "admin", promoted.Role)
}

func TestUserRepository_Promote_InvalidID(t *testing.T) {
    _, cleanup := setupUserTestDB(t)
    defer cleanup()

    repo := Repositories.NewUserRepository()

    _, err := repo.Promote("not-a-valid-hex")
    assert.Error(t, err)
}
