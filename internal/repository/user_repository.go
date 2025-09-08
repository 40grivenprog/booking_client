package repository

import (
	"booking_client/internal/models"
	"sync"
)

// UserRepository handles in-memory user storage using sync.Map
type UserRepository struct {
	storage *sync.Map
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		storage: &sync.Map{},
	}
}

// GetUser retrieves a user by chat ID
func (r *UserRepository) GetUser(chatID int64) (*models.User, bool) {
	value, exists := r.storage.Load(chatID)
	if !exists {
		return nil, false
	}

	user, ok := value.(*models.User)
	if !ok {
		return nil, false
	}

	return user, true
}

// SetUser stores a user with the given chat ID
func (r *UserRepository) SetUser(chatID int64, user *models.User) {
	r.storage.Store(chatID, user)
}

// DeleteUser removes a user by chat ID
func (r *UserRepository) DeleteUser(chatID int64) {
	r.storage.Delete(chatID)
}

// GetAllUsers returns all stored users (for debugging/admin purposes)
func (r *UserRepository) GetAllUsers() map[int64]*models.User {
	users := make(map[int64]*models.User)

	r.storage.Range(func(key, value interface{}) bool {
		chatID, ok := key.(int64)
		if !ok {
			return true // Continue iteration
		}

		user, ok := value.(*models.User)
		if !ok {
			return true // Continue iteration
		}

		users[chatID] = user
		return true // Continue iteration
	})

	return users
}

// UserExists checks if a user exists for the given chat ID
func (r *UserRepository) UserExists(chatID int64) bool {
	_, exists := r.storage.Load(chatID)
	return exists
}

// Count returns the number of stored users
func (r *UserRepository) Count() int {
	count := 0
	r.storage.Range(func(key, value interface{}) bool {
		count++
		return true // Continue iteration
	})
	return count
}
