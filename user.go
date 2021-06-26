package api

// User represents a system account
type User struct {
	Model

	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

// UserService represents a service for managing users.
type UserService interface {
	FindUserByKV(key string, val interface{}) (*User, error)
	FindUsers(filter UserFilter) ([]*User, int, error)
	CreateUser(user *User) error
}

// UserFilter represents a filter passed to FindUsers().
type UserFilter struct {
	// Filtering fields.
	ID    *int    `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`

	// Restrict to subset of results.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
