package user

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type InMemoryRepository struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*User
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		users: make(map[uuid.UUID]*User),
	}
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func (r *InMemoryRepository) CreateUser(user *User) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user

	return user, nil
}

func (r *InMemoryRepository) GetUserByID(id uuid.UUID) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (r *InMemoryRepository) GetUserByExternalID(externalID string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.ExternalID == externalID {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

func (r *InMemoryRepository) UpdateUser(user *User) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user

	return user, nil
}

func (r *InMemoryRepository) DeleteUser(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.users, id)

	return nil
}
