package Domain

import "time"

type Task struct {
	ID          string
	Title       string
	Description string
	DueDate     time.Time
	Status      string
}

type ITaskRepository interface {
	GetAll() ([]Task, error)
	GetByID(id string) (Task, error)
	Create(task Task) (Task, error)
	Update(id string, task Task) (Task, error)
	Delete(id string) error
}

type User struct {
	ID       string
	Email    string
	Password string
	Role     string
}

type IUserRepository interface {
    FindByEmail(email string) (User, error)
    Create(user User) (User, error)
    Promote(id string) (User, error)
}


type IPasswordService interface {
	Hash(password string) (string, error)
	Compare(plain, hashed string) bool
}

type AuthClaims struct {
	Email string
	Role  string
}
type IJWTService interface {
	GenerateToken(user User) (string, error)
	ValidateToken(tokenString string) (*AuthClaims, error) 
}