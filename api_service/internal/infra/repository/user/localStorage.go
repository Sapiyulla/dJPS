package user

import (
	"api-service/internal/application/ports"
	"api-service/internal/domain"
	"context"
	"sync"

	"github.com/google/uuid"
)

type UserLocalStorage struct {
	storage map[uuid.UUID]*domain.User

	mu sync.Mutex
}

func NewUserLocalStorage() *UserLocalStorage {
	return &UserLocalStorage{
		storage: make(map[uuid.UUID]*domain.User),
		mu:      sync.Mutex{},
	}
}

func (s *UserLocalStorage) Create(ctx context.Context, user *domain.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.storage[user.ID()]; exists {
		return ports.ErrUserAlreadyExists
	}
	for _, u := range s.storage {
		if u.Email() == user.Email() {
			return ports.ErrUserAlreadyExists
		}
	}

	s.storage[user.ID()] = user

	return nil
}

func (s *UserLocalStorage) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.storage[id]
	if !exists {
		return nil, ports.ErrUserNotFound
	}

	return user, nil
}

func (s *UserLocalStorage) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, user := range s.storage {
		if user.Email() == email {
			return user, nil
		}
	}

	return nil, ports.ErrUserNotFound
}

func (s *UserLocalStorage) Delete(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.storage[id]; !exists {
		return ports.ErrUserNotFound
	}

	delete(s.storage, id)

	return nil
}
