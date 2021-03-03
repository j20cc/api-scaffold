package mysql

import (
	"github.com/lukedever/api"
	"gorm.io/gorm"
)

var _ api.UserService = (*UserService)(nil)

// UserService represents a service for managing users
type UserService struct {
	db *gorm.DB
}

// NewUserService returns a new instance of UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// FindUserByID retrieves a user by ID
func (s *UserService) FindUserByID(id int) (*api.User, error) {
	return nil, nil
}

// FindUsers retrieves users by filter
func (s *UserService) FindUsers(filter api.UserFilter) ([]*api.User, int, error) {
	return nil, 0, nil
}

// CreateUser create a new user
func (s *UserService) CreateUser(user *api.User) error {
	return nil
}