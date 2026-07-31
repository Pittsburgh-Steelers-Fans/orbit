package serviceiface

import (
	"context"
	"time"
)

// User is the minimal user shape exchanged between services.
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Project is the minimal project shape exchanged between services.
type Project struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Name    string `json:"name"`
}

// Task is the minimal task shape exchanged between services.
type Task struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	DueAt     time.Time `json:"due_at"`
}

// Comment is the minimal comment shape exchanged between services.
type Comment struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	AuthorID  string    `json:"author_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// Notification is the minimal notification shape exchanged between services.
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Kind      string    `json:"kind"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

// UserService describes user operations exposed to handlers and jobs through dependency injection.
type UserService interface {
	Create(ctx context.Context, user User) (User, error)
	Find(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
}

// ProjectService describes project operations exposed to handlers and jobs through dependency injection.
type ProjectService interface {
	Create(ctx context.Context, project Project) (Project, error)
	Find(ctx context.Context, id string) (Project, error)
	ListByOwner(ctx context.Context, ownerID string) ([]Project, error)
}

// TaskService describes task operations exposed to handlers and jobs through dependency injection.
type TaskService interface {
	Create(ctx context.Context, task Task) (Task, error)
	Find(ctx context.Context, id string) (Task, error)
	ListByProject(ctx context.Context, projectID string) ([]Task, error)
	Complete(ctx context.Context, id string) (Task, error)
}

// CommentService describes comment operations exposed to handlers and jobs through dependency injection.
type CommentService interface {
	Create(ctx context.Context, comment Comment) (Comment, error)
	ListByTask(ctx context.Context, taskID string) ([]Comment, error)
	Delete(ctx context.Context, id string) error
}

// NotificationService describes notification operations exposed to handlers and jobs through dependency injection.
type NotificationService interface {
	Notify(ctx context.Context, notification Notification) error
	ListByUser(ctx context.Context, userID string) ([]Notification, error)
	MarkRead(ctx context.Context, id string) error
}
