package users

import (
	"context"
	"errors"
	"sync"
)

// ErrNotFound is returned when a user cannot be found.
var ErrNotFound = errors.New("users: not found")

// Repository stores and retrieves users.
type Repository interface {
	List(ctx context.Context) ([]User, error)
	Get(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id string) error
}

// MemoryRepository is a concurrency-safe in-memory user repository.
type MemoryRepository struct {
	mu    sync.RWMutex
	users map[string]User
}

// NewMemoryRepository creates a repository seeded with optional users.
func NewMemoryRepository(seed ...User) *MemoryRepository {
	repo := &MemoryRepository{users: make(map[string]User, len(seed))}
	for _, user := range seed {
		repo.users[user.ID] = user
	}
	return repo
}

// List returns all users in the repository.
func (r *MemoryRepository) List(ctx context.Context) ([]User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	users := make([]User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

// Get returns one user by id.
func (r *MemoryRepository) Get(ctx context.Context, id string) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return user, nil
}

// Update stores a replacement user.
func (r *MemoryRepository) Update(ctx context.Context, user User) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[user.ID]; !ok {
		return User{}, ErrNotFound
	}
	r.users[user.ID] = user
	return user, nil
}

// Delete removes a user by id.
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return ErrNotFound
	}
	delete(r.users, id)
	return nil
}
